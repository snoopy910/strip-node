package signer

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

	fmt.Println("processed keygen message with status: ", ok)
}

func generateKeygen(identity string, identityCurve string, keyCurve string, signers []string) {
	fmt.Println(signers)
	fmt.Println(NodePublicKey)

	Index := SliceIndexOfString(signers, NodePublicKey)

	if Index == -1 {
		fmt.Println("signer is not in consortium for keygen generation")
		return
	}

	TotalSigners := len(signers)

	if TotalSigners > MaximumSigners {
		fmt.Println("too many signers for keygen generation")
		return
	}

	if TotalSigners == 0 {
		fmt.Println("not enough signers for keygen generation")
		return
	}

	keyShare, err := GetKeyShare(identity, identityCurve, keyCurve)

	fmt.Println("key share from postgres: ", keyShare, err)

	if err != nil {
		fmt.Println("error from postgres:", err)
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

	// EdDSA channels (for EdDSA-based curves)
	saveChanEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanBitcoinEcdsa := make(chan *ecdsaKeygen.LocalPartySaveData)
	saveChanAptosEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanSuiEddsa := make(chan *eddsaKeygen.LocalPartySaveData)

	// ECDSA channels (for ECDSA-based curves)
	saveChanSecp256k1 := make(chan *ecdsaKeygen.LocalPartySaveData)
	saveChanEcdsa := make(chan *ecdsaKeygen.LocalPartySaveData)
	saveChanStellarEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanAlgorandEddsa := make(chan *eddsaKeygen.LocalPartySaveData)
	saveChanRippleEddsa := make(chan *eddsaKeygen.LocalPartySaveData)

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
	} else if keyCurve == SECP256K1_CURVE {
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		preParams, err := ecdsaKeygen.GeneratePreParams(2 * time.Minute)
		if err != nil {
			panic(err)
		}
		localParty := ecdsaKeygen.NewLocalParty(params, outChanKeygen, saveChanSecp256k1, *preParams)
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

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		case save := <-saveChanStellarEddsa:
			fmt.Println("saving key")

			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			// Get the public key bytes
			pkBytes := pk.Serialize()

			if len(pkBytes) != 32 {
				fmt.Println("Invalid public key length")
				return
			}

			// Version byte for ED25519 public key in Stellar
			versionByte := strkey.VersionByteAccountID // 6 << 3, or 48

			// Use Stellar SDK's strkey package to encode
			publicKeyStr, err := strkey.Encode(versionByte, pkBytes)
			if err != nil {
				fmt.Println("error encoding Stellar address: ", err)
				return
			}
			fmt.Println("new TSS Address (Stellar) is: ", publicKeyStr)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		case save := <-saveChanBitcoinEcdsa:
			fmt.Println("saving key")

			x := toHexInt(save.ECDSAPub.X())
			y := toHexInt(save.ECDSAPub.Y())
			publicKeyStr := "04" + x + y
			publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
			bitcoinAddressStr, _, _ := bitcoin.PublicKeyToBitcoinAddresses(publicKeyBytes)

			fmt.Println("new TSS Address (BTC) is: ", bitcoinAddressStr)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		case save := <-saveChanSecp256k1:
			fmt.Println("saving key")

			x := toHexInt(save.ECDSAPub.X())
			y := toHexInt(save.ECDSAPub.Y())
			publicKeyStr := "04" + x + y
			dogecoinAddressStr, err := dogecoin.PublicKeyToAddress(publicKeyStr)
			if err != nil {
				fmt.Println("Error generating Dogecoin address:", err)
			}

			fmt.Println("new TSS Address (DOGE) is: ", dogecoinAddressStr)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		case save := <-saveChanAptosEddsa:
			fmt.Println("saving key")

			pk := edwards.PublicKey{
				Curve: save.EDDSAPub.Curve(),
				X:     save.EDDSAPub.X(),
				Y:     save.EDDSAPub.Y(),
			}

			publicKeyStr := hex.EncodeToString(pk.Serialize())

			fmt.Println("new TSS Address is: ", publicKeyStr)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		case save := <-saveChanRippleEddsa:
			fmt.Println("saving key")

			publicKeyStr := ripple.PublicKeyToAddress(save)

			fmt.Println("new TSS Address (Ripple) is: ", publicKeyStr)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		case save := <-saveChanSuiEddsa:
			fmt.Println("saving key")

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

			fmt.Println("new Sui TSS Address is: ", suiAddress)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", suiAddress)
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

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		case save := <-saveChanAlgorandEddsa:
			fmt.Println("saving key")

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

			fmt.Println("new TSS Address (Algorand) is: ", publicKeyStr)

			out, err := json.Marshal(save)
			if err != nil {
				fmt.Println(err)
			}

			_json := string(out)
			AddKeyShare(identity, identityCurve, keyCurve, _json)

			signersOut, err := json.Marshal(signers)
			if err != nil {
				fmt.Println(err)
			}

			AddSignersForKeyShare(identity, identityCurve, keyCurve, string(signersOut))

			completed = true
			delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

			if val, ok := keygenGeneratedChan[identity+"_"+identityCurve+"_"+keyCurve]; ok {
				val <- "generated keygen"
			}

			fmt.Println("completed saving of new keygen ", publicKeyStr)
		}
	}
}
