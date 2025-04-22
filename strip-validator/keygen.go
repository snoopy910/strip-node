package main

import (
	"encoding/json"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util/logger"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
)

func updateKeygen(identity string, identityCurve common.Curve, keyCurve common.Curve, from int, bz []byte, isBroadcast bool, to int, signers []string) {
	TotalSigners := len(signers)
	parties, _ := getParties(TotalSigners)

	//wait for 1 minute for party to be ready
	for i := 0; i < 6; i++ {
		if !partyProcesses[identity+"_"+string(identityCurve)+"_"+string(keyCurve)].Exists {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	if !partyProcesses[identity+"_"+string(identityCurve)+"_"+string(keyCurve)].Exists {
		return
	}

	party := *partyProcesses[identity+"_"+string(identityCurve)+"_"+string(keyCurve)].Party

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

	ok, err1 := party.Update(pMsg)
	if err1 != nil {
		panic(err1)
	}

	logger.Sugar().Infof("processed keygen message with status: %v", ok)
}

func generateKeygen(identity string, identityCurve common.Curve, keyCurve common.Curve, signers []string) {
	logger.Sugar().Infof("signers: %v", signers)
	logger.Sugar().Infof("NodePublicKey: %s", NodePublicKey)

	Index := SliceIndexOfString(signers, NodePublicKey)

	if Index == -1 {
		logger.Sugar().Errorw("signer is not in consortium for keygen generation")
		return
	}

	TotalSigners := len(signers)

	if TotalSigners > MaximumSigners {
		logger.Sugar().Errorw("too many signers for keygen generation")
		return
	}

	if TotalSigners == 0 {
		logger.Sugar().Errorw("not enough signers for keygen generation")
		return
	}

	keyShare, err := GetKeyShare(identity, identityCurve, keyCurve)

	logger.Sugar().Infof("key share from postgres: %s, error: %v", keyShare, err)

	if err != nil {
		logger.Sugar().Errorw("error from postgres", "error", err)
		return
	}

	if keyShare == "" {
		logger.Sugar().Infof("key share not found. continuing to generate key share")
	}

	if keyShare != "" {
		logger.Sugar().Infof("key share found. stopping to generate key share")
		return
	}

	delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))
	parties, partiesIds := getParties(TotalSigners)

	ctx := tss.NewPeerContext(parties)

	outChanKeygen := make(chan tss.Message)

	// EdDSA channels (for EdDSA-based curves)
	saveChanEddsa := make(chan *eddsaKeygen.LocalPartySaveData)

	// ECDSA channels (for ECDSA-based curves)
	saveChanEcdsa := make(chan *ecdsaKeygen.LocalPartySaveData)

	switch keyCurve {
	case common.CurveEddsa:
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanEddsa)
		partyProcesses[identity+"_"+string(identityCurve)+"_"+string(keyCurve)] = PartyProcess{&localParty, true}
		go localParty.Start()
	case common.CurveEcdsa:
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		preParams, err := ecdsaKeygen.GeneratePreParams(2 * time.Minute)
		if err != nil {
			panic(err)
		}
		localParty := ecdsaKeygen.NewLocalParty(params, outChanKeygen, saveChanEcdsa, *preParams)
		partyProcesses[identity+"_"+string(identityCurve)+"_"+string(keyCurve)] = PartyProcess{&localParty, true}
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
				Signers:       signers,
			}

			go broadcast(message)

		case save := <-saveChanEddsa:
			logger.Sugar().Infof("saving key")

			out, err := json.Marshal(save)
			if err != nil {
				logger.Sugar().Errorw("error marshalling save", "error", err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				logger.Sugar().Errorw("error marshalling signers", "error", err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

			if val, ok := keygenGeneratedChan[identity+"_"+string(identityCurve)+"_"+string(keyCurve)]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen")

		case save := <-saveChanEcdsa:
			logger.Sugar().Infof("saving key")
			out, err := json.Marshal(save)
			if err != nil {
				logger.Sugar().Errorw("error marshalling save", "error", err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				logger.Sugar().Errorw("error marshalling signers", "error", err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

			if val, ok := keygenGeneratedChan[identity+"_"+string(identityCurve)+"_"+string(keyCurve)]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen")
		}
	}
}
