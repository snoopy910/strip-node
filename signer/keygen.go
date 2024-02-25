package signer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Silent-Protocol/go-sio/db"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/mr-tron/base58"
)

func updateKeygen(identity string, identityCurve string, keyCurve string, from int, bz []byte, isBroadcast bool, to int) {

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

	fmt.Println("processed keygen message with status: ", ok)
}

func generateKeygen(identity string, identityCurve string, keyCurve string) {
	keyShare, err := db.GetKeyShare(identity, identityCurve, keyCurve)

	fmt.Println("key share from redis: ", keyShare, err)

	if err != nil && fmt.Sprint(err) != "redis: nil" {
		fmt.Println("error from redis:", err)
		return
	}

	if keyShare == "" {
		fmt.Println("key share not found. continuing to generate key share")
	}

	if keyShare != "" {
		fmt.Println("key share found. stopping to generate key share")
		return
	}

	delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)
	parties, partiesIds := getParties(TotalSigners)

	ctx := tss.NewPeerContext(parties)

	outChanKeygen := make(chan tss.Message)

	saveChanEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanEcdsa := make(chan *ecdsaKeygen.LocalPartySaveData)

	if keyCurve == EDDSA_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), Threshold)
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanEddsa)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else {
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(Threshold))
		preParams, err := ecdsaKeygen.GeneratePreParams(2 * time.Minute)
		if err != nil {
			panic(err)
		}
		localParty := ecdsaKeygen.NewLocalParty(params, outChanKeygen, saveChanEcdsa, *preParams)
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
				Type:          MESSAGE_TYPE_GENERATE_KEYGEN,
				From:          msg.GetFrom().Index,
				To:            to,
				Message:       bytes,
				IsBroadcast:   msg.IsBroadcast(),
				Identity:      identity,
				IdentityCurve: identityCurve,
				KeyCurve:      keyCurve,
			}

			go broadcast(message)

		case save := <-saveChanEddsa:
			fmt.Println("saving key")

			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			publicKeyStr := base58.Encode(pk.Serialize())

			fmt.Println("new TSS Address is: ", publicKeyStr)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			json := string(out)
			db.AddKeyShare(identity, identityCurve, keyCurve, json)

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		case save := <-saveChanEcdsa:
			fmt.Println("saving key")

			x := toHexInt(save.ECDSAPub.X())
			y := toHexInt(save.ECDSAPub.Y())
			publicKeyStr := "04" + x + y
			publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
			newTssAddressStr := publicKeyToAddress(publicKeyBytes)

			fmt.Println("new TSS Address is: ", newTssAddressStr)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			json := string(out)
			db.AddKeyShare(identity, identityCurve, keyCurve, json)

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		}
	}
}
