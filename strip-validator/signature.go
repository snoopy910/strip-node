package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/StripChain/strip-node/bitcoin"
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs/blockchains"
	"github.com/StripChain/strip-node/ripple"
	"github.com/StripChain/strip-node/util/logger"
	cmn "github.com/bnb-chain/tss-lib/v2/common"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	ecdsaSigning "github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	eddsaSigning "github.com/bnb-chain/tss-lib/v2/eddsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/stellar/go/strkey"
	"golang.org/x/crypto/blake2b"
)

func updateSignature(identity string, identityCurve common.Curve, keyCurve common.Curve, from int, bz []byte, isBroadcast bool, to int) {
	signersString, err := GetSignersForKeyShare(identity, identityCurve, keyCurve)
	if err != nil {
		logger.Sugar().Errorw("error from postgres", "error", err)
		return
	}

	if signersString == "" {
		return
	}

	signers := []string{}
	json.Unmarshal([]byte(signersString), &signers)

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

	logger.Sugar().Infof("processed signature generation message with status: %v", ok)
}

func generateSignature(identity string, blockchainID blockchains.BlockchainID, identityCurve common.Curve, keyCurve common.Curve, hash []byte) {
	keyShare, err := GetKeyShare(identity, identityCurve, keyCurve)

	if err != nil {
		logger.Sugar().Errorw("error from postgres", "error", err)
		return
	}

	if keyShare == "" {
		logger.Sugar().Errorw("key share not found. stopping to sign")
		return
	}

	if keyShare != "" {
		logger.Sugar().Infof("key share found. continuing to sign")
	}

	signersString, err := GetSignersForKeyShare(identity, identityCurve, keyCurve)
	if err != nil {
		logger.Sugar().Errorw("error from postgres", "error", err)
		return
	}

	if signersString == "" {
		logger.Sugar().Errorw("signers not found. stopping to sign")
		return
	}

	if signersString != "" {
		logger.Sugar().Infof("signers found. continuing to sign")
	}

	signers := []string{}
	json.Unmarshal([]byte(signersString), &signers)

	Index := SliceIndexOfString(signers, NodePublicKey)

	delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

	TotalSigners := len(signers)

	parties, partiesIds := getParties(TotalSigners)

	ctx := tss.NewPeerContext(parties)

	outChanKeygen := make(chan tss.Message)
	saveChan := make(chan *cmn.SignatureData)

	var rawKeyEddsa *eddsaKeygen.LocalPartySaveData
	var rawKeyEcdsa *ecdsaKeygen.LocalPartySaveData

	switch keyCurve {
	case common.CurveEddsa:
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		msg := (&big.Int{}).SetBytes(hash)

		err = json.Unmarshal([]byte(keyShare), &rawKeyEddsa)
		if err != nil {
			logger.Sugar().Errorw("error unmarshalling key share", "error", err)
			return
		}
		localParty := eddsaSigning.NewLocalParty(msg, params, *rawKeyEddsa, outChanKeygen, saveChan)
		partyProcesses[identity+"_"+string(identityCurve)+"_"+string(keyCurve)] = PartyProcess{&localParty, true}

		go localParty.Start()
	case common.CurveEcdsa:
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		// msg := new(big.Int).SetBytes(crypto.Keccak256(hash))
		msg, _ := new(big.Int).SetString(string(hash), 16)
		err = json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)
		if err != nil {
			logger.Sugar().Errorw("error unmarshalling key share", "error", err)
			return
		}
		localParty := ecdsaSigning.NewLocalParty(msg, params, *rawKeyEcdsa, outChanKeygen, saveChan)
		partyProcesses[identity+"_"+string(identityCurve)+"_"+string(keyCurve)] = PartyProcess{&localParty, true}

		go localParty.Start()
	default:
		logger.Sugar().Errorw("invalid key curve", "keyCurve", keyCurve)
		return
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
				BlockchainID:  blockchainID,
				To:            to,
				Message:       bytes,
				IsBroadcast:   msg.IsBroadcast(),
				Hash:          hash,
				Identity:      identity,
				IdentityCurve: identityCurve,
				KeyCurve:      keyCurve,
			}

			go broadcast(message)
			logger.Sugar().Infof("Message broadcasted")

		case save := <-saveChan:
			completed = true

			switch blockchainID {
			case blockchains.Solana:
				pk := edwards.PublicKey{
					Curve: rawKeyEddsa.EDDSAPub.Curve(),
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
					BlockchainID:  blockchainID,
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

				go broadcast(message)
			case blockchains.Bitcoin:
				xStr := fmt.Sprintf("%064x", rawKeyEcdsa.ECDSAPub.X())
				prefix := "02"
				if rawKeyEcdsa.ECDSAPub.Y().Bit(0) == 1 {
					prefix = "03"
				}
				uncompressedPubKeyStr := prefix + xStr
				logger.Sugar().Infof("Uncompressed public key: %s", uncompressedPubKeyStr)
				compressedPubKeyStr, err := bitcoin.ConvertToCompressedPublicKey(uncompressedPubKeyStr)
				if err != nil {
					logger.Sugar().Errorw("error converting to compressed public key", "error", err)
					return
				}
				logger.Sugar().Infof("Compressed public key: %s", compressedPubKeyStr)

				final := hex.EncodeToString(save.Signature)

				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       []byte(final),
					Address:       compressedPubKeyStr, // we pass the public key in string format, hex string with length 130 starts with 04
					Identity:      identity,
					IdentityCurve: identityCurve,
					BlockchainID:  blockchainID,
					KeyCurve:      keyCurve,
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

				go broadcast(message)
			case blockchains.Dogecoin:
				x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
				y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())
				publicKeyStr := "04" + x + y
				compressedPubKeyStr, err := bitcoin.ConvertToCompressedPublicKey(publicKeyStr)
				if err != nil {
					logger.Sugar().Errorw("Error converting to compressed public key:", "error", err)
					return
				}

				final := hex.EncodeToString(save.Signature)
				logger.Sugar().Infof("Final message: %s", final)

				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       []byte(final),
					Address:       compressedPubKeyStr,
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
					BlockchainID:  blockchainID,
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

				go broadcast(message)
			case blockchains.Sui:
				// Get the Ed25519 public key
				pk := edwards.PublicKey{
					Curve: tss.Edwards(),
					X:     rawKeyEddsa.EDDSAPub.X(),
					Y:     rawKeyEddsa.EDDSAPub.Y(),
				}

				// Serialize the full Ed25519 public key
				pkBytes := pk.Serialize()

				// Convert to Sui address format (Blake2b-256 hash of public key)
				flag := byte(0x00)
				hasher, _ := blake2b.New256(nil)
				hasher.Write([]byte{flag})
				hasher.Write(pkBytes)

				arr := hasher.Sum(nil)
				suiAddress := "0x" + hex.EncodeToString(arr)

				// For Sui, we need to encode the signature in base64
				var signatureBytes [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
				signatureBuffer := bytes.NewBuffer([]byte{})
				scheme := sui_types.SignatureScheme{ED25519: &lib.EmptyEnum{}}
				signatureBuffer.WriteByte(scheme.Flag())
				signatureBuffer.Write(save.Signature)
				signatureBuffer.Write(pkBytes[:])
				copy(signatureBytes[:], signatureBuffer.Bytes())

				signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes[:])

				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       []byte(signatureBase64),
					Address:       suiAddress,
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
					BlockchainID:  blockchainID,
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

				go broadcast(message)
			case blockchains.Aptos:
				pk := edwards.PublicKey{
					Curve: tss.Edwards(),
					X:     rawKeyEddsa.EDDSAPub.X(),
					Y:     rawKeyEddsa.EDDSAPub.Y(),
				}

				publicKeyStr := hex.EncodeToString(pk.Serialize())

				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       save.Signature,
					Address:       publicKeyStr,
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
					BlockchainID:  blockchainID,
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

				go broadcast(message)
			case blockchains.Stellar:
				pk := edwards.PublicKey{
					Curve: tss.Edwards(),
					X:     rawKeyEddsa.EDDSAPub.X(),
					Y:     rawKeyEddsa.EDDSAPub.Y(),
				}

				// Get the public key bytes
				pkBytes := pk.Serialize()
				if len(pkBytes) != 32 {
					logger.Sugar().Errorw("invalid public key length")
					return
				}

				// Version byte for ED25519 public key in Stellar
				versionByte := strkey.VersionByteAccountID // 6 << 3, or 48

				// Use Stellar SDK's strkey package to encode
				address, err := strkey.Encode(versionByte, pkBytes)
				if err != nil {
					logger.Sugar().Errorw("error encoding Stellar address", "error", err)
					return
				}

				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       save.Signature,
					Address:       address,
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
					BlockchainID:  blockchainID,
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

				go broadcast(message)
			case blockchains.Algorand:
				pk := edwards.PublicKey{
					Curve: tss.Edwards(),
					X:     rawKeyEddsa.EDDSAPub.X(),
					Y:     rawKeyEddsa.EDDSAPub.Y(),
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
				address := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(addressBytes)

				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       save.Signature,
					Address:       address,
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
					BlockchainID:  blockchainID,
					AlgorandFlags: &struct {
						IsRealTransaction bool `json:"isRealTransaction"`
					}{
						IsRealTransaction: true,
					},
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))
				go broadcast(message)
			case blockchains.Ripple:
				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       save.Signature,
					Address:       ripple.PublicKeyToAddress(rawKeyEddsa),
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
					BlockchainID:  blockchainID,
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))
				go broadcast(message)
			case blockchains.Cardano:

				pk := edwards.PublicKey{
					Curve: rawKeyEddsa.EDDSAPub.Curve(),
					X:     rawKeyEddsa.EDDSAPub.X(),
					Y:     rawKeyEddsa.EDDSAPub.Y(),
				}

				publicKeyStr := hex.EncodeToString(pk.Serialize())

				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       save.Signature,
					Address:       publicKeyStr,
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
					BlockchainID:  blockchainID,
				}

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))
				go broadcast(message)
			default:
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
					BlockchainID:  blockchainID,
				}

				logger.Sugar().Infof("Address of the generated signature: %s", publicKeyToAddress(pubkey))

				delete(partyProcesses, identity+"_"+string(identityCurve)+"_"+string(keyCurve))

				go broadcast(message)
			}
		}
	}
}
