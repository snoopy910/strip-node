package signer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Silent-Protocol/go-sio/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58/base58"
)

type MessageType uint

const (
	MESSAGE_TYPE_GENERATE_KEYGEN       MessageType = 0
	MESSAGE_TYPE_GENERATE_SIGNATURE    MessageType = 1
	MESSAGE_TYPE_GENERATE_START_KEYGEN MessageType = 2
	MESSAGE_TYPE_START_SIGN            MessageType = 3
	MESSAGE_TYPE_SIGN                  MessageType = 4
	MESSAGE_TYPE_SIGNATURE             MessageType = 5
)

type Message struct {
	Identity           string      `json:"identity"`
	IdentityCurve      string      `json:"identityCurve"`
	KeyCurve           string      `json:"keyCurve"`
	From               int         `json:"from"`
	To                 int         `json:"to"`
	Message            []byte      `json:"message"`
	Type               MessageType `json:"type"`
	IsToNewCommittee   bool        `json:"isToNewCommittee"`
	IsFromNewCommittee bool        `json:"isFromNewCommittee"`
	IsBroadcast        bool        `json:"isBroadcast"`
	Hash               []byte      `json:"hash"`
	Address            string      `json:"address"`
	Signature          []byte      `json:"signature"`
	Signers            []string    `json:"signers"`
}

type IsValid struct {
	Result bool `json:"result"`
}

func handleIncomingMessage(message []byte) {
	msg := Message{}
	json.Unmarshal(message, &msg)

	_validateMsg := msg
	// _signature := msg.Signature
	_validateMsg.Signature = nil
	messageBytes, err := json.Marshal(_validateMsg)

	if err != nil {
		panic(err)
	}

	hash := crypto.Keccak256Hash(messageBytes)

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), msg.Signature)
	if err != nil {
		log.Fatal(err)
	}
	pubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	compressedPubKey := crypto.CompressPubkey(pubKey)
	compressedPubKeyStr := hexutil.Encode(compressedPubKey)

	publicKeyStr := hex.EncodeToString(sigPublicKey)
	instance := getIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)

	signerExists, err := instance.Signers(&bind.CallOpts{}, common.PublicKeyStrToBytes32(compressedPubKeyStr))
	if err != nil {
		fmt.Println("message from not registered signer: ", "0x"+publicKeyStr)
		log.Fatal(err)
	}

	if !signerExists.Exists {
		return
	}

	if msg.Type == MESSAGE_TYPE_GENERATE_START_KEYGEN {
		go generateKeygen(msg.Identity, msg.IdentityCurve, msg.KeyCurve, msg.Signers)
	} else if msg.Type == MESSAGE_TYPE_GENERATE_KEYGEN {
		go updateKeygen(msg.Identity, msg.IdentityCurve, msg.KeyCurve, msg.From, msg.Message, msg.IsBroadcast, msg.To, msg.Signers)
	} else if msg.Type == MESSAGE_TYPE_START_SIGN {
		go generateSignature(msg.Identity, msg.IdentityCurve, msg.KeyCurve, msg.Hash)
	} else if msg.Type == MESSAGE_TYPE_SIGN {
		go updateSignature(msg.Identity, msg.IdentityCurve, msg.KeyCurve, msg.From, msg.Message, msg.IsBroadcast, msg.To)
	} else if msg.Type == MESSAGE_TYPE_SIGNATURE {
		if msg.KeyCurve == EDDSA_CURVE {
			if val, ok := messageChan[base58.Encode(msg.Hash)]; ok {
				val <- msg
			}
		} else {
			if val, ok := messageChan[string(msg.Hash)]; ok {
				val <- msg
			}
		}
	}
}

func broadcast(message Message) {

	messageBytes, err := json.Marshal(message)

	if err != nil {
		panic(err)
	}

	hash := crypto.Keccak256Hash(messageBytes)

	cleanedPrivateKey := strings.Replace(NodePrivateKey, "0x", "", 1)
	privateKey, err := crypto.HexToECDSA(cleanedPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	message.Signature = signature

	out, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	if err := topic.Publish(ctx, out); err != nil {
		fmt.Println("### Publish error")
		panic(err)
	}
}
