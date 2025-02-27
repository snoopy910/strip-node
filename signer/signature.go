package signer

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
	"strings"
	"time"

	"github.com/StripChain/strip-node/cardano"
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/dogecoin"
	"github.com/StripChain/strip-node/ripple"
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

func updateSignature(identity string, identityCurve string, keyCurve string, from int, bz []byte, isBroadcast bool, to int) {
	signersString, err := GetSignersForKeyShare(identity, identityCurve, keyCurve)
	if err != nil {
		fmt.Println("error from postgres:", err)
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
	keyShare, err := GetKeyShare(identity, identityCurve, keyCurve)

	if err != nil {
		fmt.Println("error from postgres:", err)
		return
	}

	if keyShare == "" {
		fmt.Println("key share not found. stopping to sign")
		return
	}

	if keyShare != "" {
		fmt.Println("key share found. continuing to sign")
	}

	signersString, err := GetSignersForKeyShare(identity, identityCurve, keyCurve)
	if err != nil {
		fmt.Println("error from postgres:", err)
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

	} else if keyCurve == SUI_EDDSA_CURVE {
		// Sui uses Ed25519 for transaction signing
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))

		msg := (&big.Int{}).SetBytes(hash)

		json.Unmarshal([]byte(keyShare), &rawKeyEddsa)
		localParty := eddsaSigning.NewLocalParty(msg, params, *rawKeyEddsa, outChanKeygen, saveChan)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}

		go localParty.Start()

	} else if keyCurve == SECP256K1_CURVE {
		params := tss.NewParameters(tss.S256(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		msg := new(big.Int).SetBytes(crypto.Keccak256(hash))
		json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)
		localParty := ecdsaSigning.NewLocalParty(msg, params, *rawKeyEcdsa, outChanKeygen, saveChan)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}

		go localParty.Start()
	} else if keyCurve == APTOS_EDDSA_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		msg, _ := new(big.Int).SetString(string(hash), 16)

		json.Unmarshal([]byte(keyShare), &rawKeyEddsa)
		localParty := eddsaSigning.NewLocalParty(msg, params, *rawKeyEddsa, outChanKeygen, saveChan)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}

		go localParty.Start()
	} else if keyCurve == STELLAR_CURVE || keyCurve == RIPPLE_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		msg := new(big.Int).SetBytes(hash)

		json.Unmarshal([]byte(keyShare), &rawKeyEddsa)
		localParty := eddsaSigning.NewLocalParty(msg, params, *rawKeyEddsa, outChanKeygen, saveChan)
		partyProcesses[identity+"_"+identityCurve+"_"+keyCurve] = PartyProcess{&localParty, true}

		go localParty.Start()
	} else if keyCurve == ALGORAND_CURVE {
		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[Index], len(parties), int(CalculateThreshold(TotalSigners)))
		msg := new(big.Int).SetBytes(hash)

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
			} else if keyCurve == SECP256K1_CURVE {
				x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
				y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())
				publicKeyStr := "04" + x + y
				publicKeyBytes, _ := hex.DecodeString(publicKeyStr)

				// Get chain information from metadata
				var address string
				var metadata map[string]interface{}
				if err := json.Unmarshal(hash, &metadata); err != nil {
					fmt.Println("Error parsing metadata:", err)
					// Default to Bitcoin address format if metadata is invalid
					address, _, _ = publicKeyToBitcoinAddresses(publicKeyBytes)
				} else {
					// Get chain information
					if chainId, ok := metadata["chainId"].(string); ok {
						chain, err := common.GetChain(chainId)
						if err != nil {
							fmt.Printf("Error getting chain info: %v\n", err)
							address, _, _ = publicKeyToBitcoinAddresses(publicKeyBytes)
						} else {
							switch chain.ChainType {
							case "dogecoin":
								// Use Dogecoin address format
								if strings.HasSuffix(chainId, "1") { // Testnet
									address, err = dogecoin.PublicKeyToTestnetAddress(publicKeyStr)
								} else { // Mainnet
									address, err = dogecoin.PublicKeyToAddress(publicKeyStr)
								}
								if err != nil {
									fmt.Printf("Error generating Dogecoin address: %v\n", err)
									address, _, _ = publicKeyToBitcoinAddresses(publicKeyBytes)
								}
							case "bitcoin":
								address, _, _ = publicKeyToBitcoinAddresses(publicKeyBytes)
							default:
								// Default to Bitcoin address format for unknown chains
								address, _, _ = publicKeyToBitcoinAddresses(publicKeyBytes)
							}
						}
					} else {
						// Default to Bitcoin if no chainId in metadata
						address, _, _ = publicKeyToBitcoinAddresses(publicKeyBytes)
					}
				}

				final := hex.EncodeToString(save.Signature) + hex.EncodeToString(save.SignatureRecovery)

				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       []byte(final),
					Address:       address,
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
				}

				delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

				go broadcast(message)
			} else if keyCurve == SUI_EDDSA_CURVE {
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
				}

				delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

				go broadcast(message)
			} else if keyCurve == APTOS_EDDSA_CURVE {
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
				}

				delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

				go broadcast(message)
			} else if keyCurve == STELLAR_CURVE {
				pk := edwards.PublicKey{
					Curve: tss.Edwards(),
					X:     rawKeyEddsa.EDDSAPub.X(),
					Y:     rawKeyEddsa.EDDSAPub.Y(),
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
				address, err := strkey.Encode(versionByte, pkBytes)
				if err != nil {
					fmt.Println("error encoding Stellar address: ", err)
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
				}

				delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

				go broadcast(message)
			} else if keyCurve == ALGORAND_CURVE {
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
				}

				delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)
				go broadcast(message)
			} else if keyCurve == RIPPLE_CURVE {
				message := Message{
					Type:          MESSAGE_TYPE_SIGNATURE,
					Hash:          hash,
					Message:       save.Signature,
					Address:       ripple.PublicKeyToAddress(rawKeyEddsa),
					Identity:      identity,
					IdentityCurve: identityCurve,
					KeyCurve:      keyCurve,
				}

				delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)
				go broadcast(message)
			} else if keyCurve == CARDANO_CURVE {
				address, err := cardano.PublicKeyToAddress(rawKeyEddsa, "1006")
				if err != nil {
					fmt.Println("error converting Cardano public key to address: ", err)
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

				fmt.Println("Address of the generated signature: ", publicKeyToAddress(pubkey))

				delete(partyProcesses, identity+"_"+identityCurve+"_"+keyCurve)

				go broadcast(message)
			}
		}
	}
}
