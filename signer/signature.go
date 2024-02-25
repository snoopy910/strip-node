package signer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/Silent-Protocol/go-sio/db"
	cmn "github.com/bnb-chain/tss-lib/v2/common"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	ecdsaSigning "github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	eddsaSigning "github.com/bnb-chain/tss-lib/v2/eddsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
)

func updateSignature(identity string, identityCurve string, keyCurve string, from int, bz []byte, isBroadcast bool, to int) {
	fmt.Println("xxxxxxx")

	signersString, err := db.GetSignersForKeyShare(identity, identityCurve, keyCurve)
	if err != nil && fmt.Sprint(err) != "redis: nil" {
		fmt.Println("error from redis:", err)
		return
	}

	if signersString == "" {
		fmt.Println("signers not found. stopping signing")
		return
	}

	if signersString != "" {
		fmt.Println("signers found. continuing to sign")
	}

	signers := []string{}
	json.Unmarshal([]byte(signersString), &signers)

	TotalSigners := len(signers)

	parties, _ := getParties(TotalSigners)

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
		fmt.Println("key share not found. stopping to sign")
		return
	}

	if keyShare != "" {
		fmt.Println("key share found. continuing to sign")
	}

	signersString, err := db.GetSignersForKeyShare(identity, identityCurve, keyCurve)
	if err != nil && fmt.Sprint(err) != "redis: nil" {
		fmt.Println("error from redis:", err)
		return
	}

	if signersString == "" {
		fmt.Println("signers not found. stopping to sign")
		return
	}

	if signersString != "" {
		fmt.Println("signers found. continuing to sign")
	}

	signers := []string{}
	json.Unmarshal([]byte(signersString), &signers)

	Index := SliceIndexOfString(signers, NodePublicKey)

	delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

	TotalSigners := len(signers)

	parties, partiesIds := getParties(TotalSigners)

	ctx := tss.NewPeerContext(parties)

	outChanKeygen := make(chan tss.Message)
	saveChan := make(chan *cmn.SignatureData)

	var rawKeyEddsa *eddsaKeygen.LocalPartySaveData
	var rawKeyEcdsa *ecdsaKeygen.LocalPartySaveData

	if keyCurve == EDDSA_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		msg := (&big.Int{}).SetBytes(hash)

		json.Unmarshal([]byte(keyShare), &rawKeyEddsa)
		localParty := eddsaSigning.NewLocalParty(msg, params, *rawKeyEddsa, outChanKeygen, saveChan)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}

		go localParty.Start()

	} else {
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		msg, _ := new(big.Int).SetString(string(hash), 16)
		json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)
		localParty := ecdsaSigning.NewLocalParty(msg, params, *rawKeyEcdsa, outChanKeygen, saveChan)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}

		go localParty.Start()
	}

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
			fmt.Println("Message broadcasted")

		case save := <-saveChan:
			completed = true

			if keyCurve == EDDSA_CURVE {
				pk := edwards.PublicKey{
					Curve: tss.Edwards(),
					X:     rawKeyEddsa.EDDSAPub.X(),
					Y:     rawKeyEddsa.EDDSAPub.Y(),
				}

				publicKeyStr := base58.Encode(pk.Serialize())

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
			} else {
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
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       []byte(final),
					Address:       publicKeyToAddress(pubkey),
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
				}

				delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

				go broadcast(message)
			}
		}
	}
}
