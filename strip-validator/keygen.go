package main

import (
	"crypto/sha512"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/StripChain/strip-node/bitcoin"
	"github.com/StripChain/strip-node/dogecoin"
	"github.com/StripChain/strip-node/ripple"
	"github.com/StripChain/strip-node/util/logger"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/mr-tron/base58"
	"github.com/stellar/go/strkey"
	"golang.org/x/crypto/blake2b"
)

func updateKeygen(identity string, identityCurve string, keyCurve string, from int, bz []byte, isBroadcast bool, to int, signers []string) {
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

	logger.Sugar().Infof("processed keygen message with status: %v", ok)
}

func generateKeygen(identity string, identityCurve string, keyCurve string, signers []string) {
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

	delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)
	parties, partiesIds := getParties(TotalSigners)

	ctx := tss.NewPeerContext(parties)

	outChanKeygen := make(chan tss.Message)

	// EdDSA channels (for EdDSA-based curves)
	saveChanEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanAptosEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanSuiEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanStellarEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanAlgorandEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanRippleEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanCardanoEddsa := make(chan *eddsaKeygen.LocalPartySaveData)

	// ECDSA channels (for ECDSA-based curves)
	saveChanBitcoinEcdsa := make(chan *ecdsaKeygen.LocalPartySaveData)
	saveChanDogecoinEcdsa := make(chan *ecdsaKeygen.LocalPartySaveData)
	saveChanEcdsa := make(chan *ecdsaKeygen.LocalPartySaveData)

	if keyCurve == EDDSA_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanEddsa)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else if keyCurve == SUI_EDDSA_CURVE {
		// Sui uses Ed25519 for native transactions
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanSuiEddsa)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else if keyCurve == STELLAR_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanStellarEddsa)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()

	} else if keyCurve == BITCOIN_CURVE {
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		preParams, err := ecdsaKeygen.GeneratePreParams(2 * time.Minute)
		if err != nil {
			panic(err)
		}
		localParty := ecdsaKeygen.NewLocalParty(params, outChanKeygen, saveChanBitcoinEcdsa, *preParams)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else if keyCurve == DOGECOIN_CURVE {
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		preParams, err := ecdsaKeygen.GeneratePreParams(2 * time.Minute)
		if err != nil {
			panic(err)
		}
		localParty := ecdsaKeygen.NewLocalParty(params, outChanKeygen, saveChanDogecoinEcdsa, *preParams)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else if keyCurve == APTOS_EDDSA_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanAptosEddsa)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else if keyCurve == ALGORAND_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanAlgorandEddsa)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else if keyCurve == RIPPLE_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanRippleEddsa)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else if keyCurve == CARDANO_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		localParty := eddsaKeygen.NewLocalParty(params, outChanKeygen, saveChanCardanoEddsa)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}
		go localParty.Start()
	} else {
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
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
				Signers:       signers,
			}

			go broadcast(message)

		case save := <-saveChanEddsa:
			logger.Sugar().Infof("saving key")

			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			publicKeyStr := base58.Encode(pk.Serialize())

			logger.Sugar().Infof("new TSS Address is: %s", publicKeyStr)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)
		case save := <-saveChanStellarEddsa:
			logger.Sugar().Infof("saving key")

			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			// Get the public key bytes
			pkBytes := pk.Serialize()

			if len(pkBytes) != 32 {
				logger.Sugar().Errorw("Invalid public key length")
				return
			}

			// Version byte for ED25519 public key in Stellar
			versionByte := strkey.VersionByteAccountID // 6 << 3, or 48

			// Use Stellar SDK's strkey package to encode
			publicKeyStr, err := strkey.Encode(versionByte, pkBytes)
			if err != nil {
				logger.Sugar().Errorw("error encoding Stellar address", "error", err)
				return
			}
			logger.Sugar().Infof("new TSS Address (Stellar) is: %s", publicKeyStr)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)
		case save := <-saveChanBitcoinEcdsa:
			logger.Sugar().Infof("saving key")

			xStr := fmt.Sprintf("%064x", save.ECDSAPub.X())
			prefix := "02"
			if save.ECDSAPub.Y().Bit(0) == 1 {
				prefix = "03"
			}
			publicKeyStr := prefix + xStr
			publicKeyBytes, err := hex.DecodeString(publicKeyStr)
			if err != nil {
				logger.Sugar().Errorw("error decoding public key", "error", err)
				return
			}
			bitcoinAddressStr, _, _ := bitcoin.PublicKeyToBitcoinAddresses(publicKeyBytes)

			logger.Sugar().Infof("new TSS Address (BTC) is: %s", bitcoinAddressStr)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}
			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)

		case save := <-saveChanDogecoinEcdsa:
			logger.Sugar().Infof("saving key")

			x := toHexInt(save.ECDSAPub.X())
			y := toHexInt(save.ECDSAPub.Y())
			publicKeyStr := "04" + x + y
			dogecoinAddressStr, err := dogecoin.PublicKeyToAddress(publicKeyStr)
			if err != nil {
				logger.Sugar().Errorw("Error generating Dogecoin address", "error", err)
			}

			logger.Sugar().Infof("new TSS Address (DOGE) is: %s", dogecoinAddressStr)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)
		case save := <-saveChanAptosEddsa:
			logger.Sugar().Infof("saving key")

			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			publicKeyStr := hex.EncodeToString(pk.Serialize())

			logger.Sugar().Infof("new TSS Address is: %s", publicKeyStr)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)
		case save := <-saveChanRippleEddsa:
			logger.Sugar().Infof("saving key")

			publicKeyStr := ripple.PublicKeyToAddress(save)

			logger.Sugar().Infof("new TSS Address (Ripple) is: %s", publicKeyStr)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)
		case save := <-saveChanCardanoEddsa:
			logger.Sugar().Infof("saving key")
			// Get the public key
			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			publicKeyStr := hex.EncodeToString(pk.Serialize())

			logger.Sugar().Infof("new TSS Address (Cardano) is: %s", publicKeyStr)
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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}
			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)
		case save := <-saveChanSuiEddsa:
			logger.Sugar().Infof("saving key")

			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			// Serialize the Ed25519 public key
			pkBytes := pk.Serialize()

			// Hash the public key with Blake2b-256 to get Sui address
			hasher := blake2b.Sum256(pkBytes)
			suiAddress := "0x" + hex.EncodeToString(hasher[:])

			logger.Sugar().Infof("new Sui TSS Address is: %s", suiAddress)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen %s", suiAddress)
		case save := <-saveChanEcdsa:
			logger.Sugar().Infof("saving key")

			x := toHexInt(save.ECDSAPub.X())
			y := toHexInt(save.ECDSAPub.Y())
			publicKeyStr := "04" + x + y
			publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
			newTssAddressStr := publicKeyToAddress(publicKeyBytes)

			logger.Sugar().Infof("new TSS Address is: %s", newTssAddressStr)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)
		case save := <-saveChanAlgorandEddsa:
			logger.Sugar().Infof("saving key")

			// For Algorand, we use the EdDSA key
			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			// Get the public key bytes
			pkBytes := pk.Serialize()

			// Calculate checksum (last 4 bytes of SHA512/256 hash)
			hasher := sha512.New512_256()
			hasher.Write(pkBytes)
			checksum := hasher.Sum(nil)[28:] // Last 4 bytes

			// Concatenate public key and checksum
			addressBytes := append(pkBytes, checksum...)

			// Encode in base32 without padding
			publicKeyStr := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(addressBytes)

			logger.Sugar().Infof("new TSS Address (Algorand) is: %s", publicKeyStr)

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
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			logger.Sugar().Infof("completed saving of new keygen %s", publicKeyStr)
		}
	}
}
