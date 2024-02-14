package signer

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	cmn "github.com/bnb-chain/tss-lib/common"
	"github.com/bnb-chain/tss-lib/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/ethereum/go-ethereum/crypto"
)

func updateSignature(networkId string, partyKey string, from int, bz []byte, isBroadcast bool, to int) {
	parties, _ := getParties(networks[networkId].TotalSigners, networks[networkId].StartKeyInt)

	//wait for 1 minute for party to be ready
	for i := 0; i < 6; i++ {
		if !partyProcesses[networkId][partyKey].Exists {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	if !partyProcesses[networkId][partyKey].Exists {
		return
	}

	party := *partyProcesses[networkId][partyKey].Party

	if to != -1 && to != party.PartyID().Index {
		return
	}

	if party.PartyID().Index == from {
		return
	}

	pMsg, err := tss.ParseWireMessage(bz, parties[from], isBroadcast)
	if err != nil {
		panic(err)
	}

	ok, err := party.Update(pMsg)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("processed signature generation message with status: ", ok)
}

func generateSignature(networkId string, hash []byte) {
	if networks[networkId].Key == nil {
		return
	}

	parties, partiesIds := getParties(networks[networkId].TotalSigners, networks[networkId].StartKeyInt)

	ctx := tss.NewPeerContext(parties)

	outChanKeygen := make(chan tss.Message)
	saveChan := make(chan cmn.SignatureData)

	params := tss.NewParameters(tss.S256(), ctx, partiesIds[networks[networkId].Index], len(parties), networks[networkId].Threshold)

	msg, _ := new(big.Int).SetString(string(hash), 16)

	localParty := signing.NewLocalParty(msg, params, *networks[networkId].Key, outChanKeygen, saveChan)
	partyProcesses[networkId][string(hash)] = PartyProcess{&localParty, true}

	go localParty.Start()

	completed := false
	for !completed {
		select {
		case msg := <-outChanKeygen:
			dest := msg.GetTo()
			bytes, _, _ := msg.WireBytes()
			to := 0
			if dest == nil {
				to = -1
			} else {
				to = dest[0].Index
			}

			message := Message{
				Type:        MESSAGE_TYPE_SIGN,
				From:        msg.GetFrom().Index,
				To:          to,
				Message:     bytes,
				IsBroadcast: msg.IsBroadcast(),
				Hash:        hash,
				NetworkId:   networkId,
			}

			go broadcast(message)

		case save := <-saveChan:
			completed = true

			final := hex.EncodeToString(save.Signature) + hex.EncodeToString(save.SignatureRecovery)
			data, err := hex.DecodeString(string(hash))
			if err != nil {
				panic(err)
			}

			sdata, err := hex.DecodeString(final)
			if err != nil {
				panic(err)
			}
			pubkey, err := crypto.Ecrecover(data, sdata)
			if err != nil {
				panic(err)
			}

			message := Message{
				Type:      MESSAGE_TYPE_SIGNATURE,
				Hash:      hash,
				Message:   []byte(final),
				Address:   publicKeyToAddress(pubkey),
				NetworkId: networkId,
			}

			delete(partyProcesses, string(hash))

			go broadcast(message)
		}
	}
}
