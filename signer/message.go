package signer

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"

	"github.com/StripChain/strip-node/common"
	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/util/logger"
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
	AlgorandFlags      *struct {
		IsRealTransaction bool `json:"isRealTransaction"`
	} `json:"algorandFlags,omitempty"`
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

	compressedPubKeyStr = compressedPubKeyStr[4:]
	compressedPubKeyStr = "0x" + compressedPubKeyStr

	instance := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)

	signerExists, err := instance.Signers(&bind.CallOpts{}, common.PublicKeyStrToBytes32(compressedPubKeyStr))
	if err != nil {
		logger.Sugar().Errorw("message from not registered signer", "error", err)
		log.Fatal(err)
	}

	if !signerExists {
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
		// When looking up the channel to send back the signature, we must encode the hash
		// in the same format that was used when creating the channel in api.go.
		// This ensures we find the correct channel for each chain's message format.

		sendMsg := make(chan bool)
		switch msg.KeyCurve {
		case EDDSA_CURVE:
			// Solana: Client sends base58 string -> decode -> process -> encode back to base58
			// Channel key must match the original base58 format from client
			if val, ok := messageChan[base58.Encode(msg.Hash)]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		case BITCOIN_CURVE:
			// Bitcoin: Client sends string -> hash -> process -> encode to hex
			// Channel key must match the hex encoded hash
			if val, ok := messageChan[string(msg.Hash)]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		case SECP256K1_CURVE:
			// Bitcoin: Client sends string -> hash -> process -> encode to hex
			// Channel key must match the hex encoded hash
			if val, ok := messageChan[hex.EncodeToString(msg.Hash)]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		case APTOS_EDDSA_CURVE, CARDANO_CURVE:
			// Aptos: Client sends string -> process -> encode to hex
			// Channel key must match the hex encoded format
			if val, ok := messageChan[hex.EncodeToString(msg.Hash)]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		case RIPPLE_CURVE:
			// Ripple: Client sends string -> process -> encode to hex
			// Channel key must match the hex encoded format
			if val, ok := messageChan[strings.ToUpper(hex.EncodeToString(msg.Hash))]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		case STELLAR_CURVE:
			// Stellar: Client sends base64 string -> decode -> process -> encode back to base64
			// Channel key must match the original base64 format from client, using StrKey encoding
			if val, ok := messageChan[base64.StdEncoding.EncodeToString(msg.Hash)]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		case SUI_EDDSA_CURVE:
			// Sui: Client sends base64 string -> decode -> process -> encode back to base64
			// Channel key must match the original base32 format from client
			if val, ok := messageChan[base64.StdEncoding.EncodeToString(msg.Hash)]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		case ALGORAND_CURVE:
			// Algorand: Client sends base64 string -> decode -> process -> encode back to base64
			// Channel key must match the original base32 format from client
			if val, ok := messageChan[base64.StdEncoding.EncodeToString(msg.Hash)]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		default:
			// Ethereum: Client sends string -> process raw bytes
			// Channel key must match the raw bytes as string
			if val, ok := messageChan[string(msg.Hash)]; ok {
				val <- msg
				go func() {
					sendMsg <- true
				}()
			}
		}

		val := <-sendMsg
		if val {
			logger.Sugar().Infof("message sent: %v", val)
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
		logger.Sugar().Errorw("failed to publish message", "error", err)
		panic(err)
	}
}
