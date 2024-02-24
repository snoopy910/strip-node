package signer

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/Silent-Protocol/go-sio/db"
	cmn "github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/eddsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/mr-tron/base58"
)

func updateSignature(identity string, identityCurve string, keyCurve string, from int, bz []byte, isBroadcast bool, to int) {
	parties, _ := getParties(TotalSigners, StartKey)

	//wait for 1 minute for party to be ready
	for i := 0; i < 6; i++ {
		if !partyProcesses[identity+"_"+identityCurve+"_"+keyCurve].Exists {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	if !partyProcesses[identity+"_"+identityCurve+"_"+keyCurve].Exists {
		return
	}

	party := *partyProcesses[identity+"_"+identityCurve+"_"+keyCurve].Party

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

func generateSignature(identity string, identityCurve string, keyCurve string, hash []byte) {
	keyShare, err := db.GetKeyShare(identity, identityCurve, keyCurve)

	if err != nil && fmt.Sprint(err) != "redis: nil" {
		fmt.Println("error from redis:", err)
		return
	}

	if keyShare == "" {
		fmt.Println("key share not found. stopping to generate key share")
		return
	}

	if keyShare != "" {
		fmt.Println("key share found. continuing to generate key share")
	}

	delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

	parties, partiesIds := getParties(TotalSigners, StartKey)

	ctx := tss.NewPeerContext(parties)

	outChanKeygen := make(chan tss.Message)
	saveChan := make(chan *cmn.SignatureData)

	params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), Threshold)

	msg := (&big.Int{}).SetBytes(hash)
	// msg, _ := new(big.Int).SetString(string(hash), 16)
	// fmt.Println(hex.EncodeToString(hash))

	var rawKey *keygen.LocalPartySaveData
	json.Unmarshal([]byte(keyShare), &rawKey)

	localParty := signing.NewLocalParty(msg, params, *rawKey, outChanKeygen, saveChan)
	partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}

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
				Type:          MESSAGE_TYPE_SIGN,
				From:          msg.GetFrom().Index,
				To:            to,
				Message:       bytes,
				IsBroadcast:   msg.IsBroadcast(),
				Hash:          hash,
				Identity:      identity,
				IdentityCurve: identityCurve,
				KeyCurve:      keyCurve,
			}

			go broadcast(message)

		case save := <-saveChan:
			completed = true

			// final := base58.Encode(save.Signature)

			pk := edwards.PublicKey{
				Curve: tss.Edwards(),
				X:     rawKey.EDDSAPub.X(),
				Y:     rawKey.EDDSAPub.Y(),
			}

			publicKeyStr := base58.Encode(pk.Serialize())

			newSig, err := edwards.ParseSignature(save.Signature)
			if err != nil {
				println("new sig error, ", err.Error())
			}

			ok := edwards.Verify(&pk, hash, newSig.R, newSig.S)
			fmt.Println(ok)

			verified := ed25519.Verify(ed25519.PublicKey(pk.Serialize()), hash, save.Signature)
			fmt.Println(verified)

			message := Message{
				Type:          MESSAGE_TYPE_SIGNATURE,
				Hash:          hash,
				Message:       save.Signature,
				Address:       publicKeyStr,
				Identity:      identity,
				IdentityCurve: identityCurve,
				KeyCurve:      keyCurve,
			}

			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			go broadcast(message)
		}
	}
}
