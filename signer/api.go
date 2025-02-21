package signer

import (
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	identityVerification "github.com/StripChain/strip-node/identity"
	"github.com/StripChain/strip-node/sequencer"
	"github.com/stellar/go/strkey"

	"github.com/StripChain/strip-node/solver"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
)

var messageChan = make(map[string]chan (Message))
var keygenGeneratedChan = make(map[string]chan (string))

var (
	ECDSA_CURVE       = "ecdsa"
	EDDSA_CURVE       = "eddsa"
	APTOS_EDDSA_CURVE = "aptos_eddsa"
	SECP256K1_CURVE   = "secp256k1"
	STELLAR_CURVE     = "stellar_eddsa" // Stellar uses Ed25519 with StrKey encoding
	ALGORAND_CURVE    = "algorand_eddsa"
	// Note: Hedera uses ECDSA_CURVE since it's compatible with EVM
)

func generateKeygenMessage(identity string, identityCurve string, keyCurve string, signers []string) {
	message := Message{
		Type:          MESSAGE_TYPE_GENERATE_START_KEYGEN,
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      keyCurve,
		Signers:       signers,
	}

	broadcast(message)
}

func generateSignatureMessage(identity string, identityCurve string, keyCurve string, msg []byte) {
	message := Message{
		Type:          MESSAGE_TYPE_START_SIGN,
		Hash:          msg,
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      keyCurve,
	}

	broadcast(message)
}

type CreateWallet struct {
	Identity      string   `json:"identity"`
	IdentityCurve string   `json:"identityCurve"`
	KeyCurve      string   `json:"keyCurve"`
	Signers       []string `json:"signers"`
}

type SignMessage struct {
	Message       string `json:"message"`
	Identity      string `json:"identity"`
	IdentityCurve string `json:"identityCurve"`
	KeyCurve      string `json:"keyCurve"`
}

func startHTTPServer(port string) {
	http.HandleFunc("/keygen", func(w http.ResponseWriter, r *http.Request) {
		var createWallet CreateWallet

		err := json.NewDecoder(r.Body).Decode(&createWallet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key := createWallet.Identity + "_" + createWallet.IdentityCurve + "_" + createWallet.KeyCurve

		keygenGeneratedChan[key] = make(chan string)

		go generateKeygenMessage(createWallet.Identity, createWallet.IdentityCurve, createWallet.KeyCurve, createWallet.Signers)

		<-keygenGeneratedChan[key]
		delete(keygenGeneratedChan, key)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")
		keyCurve := r.URL.Query().Get("keyCurve")

		keyShare, err := GetKeyShare(identity, identityCurve, keyCurve)

		if err != nil {
			http.Error(w, "error from postgres", http.StatusBadRequest)
			return
		}

		if keyShare == "" {
			http.Error(w, "key share not found.", http.StatusBadRequest)
			return
		}

		var rawKeyEddsa *eddsaKeygen.LocalPartySaveData
		var rawKeyEcdsa *ecdsaKeygen.LocalPartySaveData

		if keyCurve == EDDSA_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

			pk := edwards.PublicKey{
				Curve: tss.Edwards(),
				X:     rawKeyEddsa.EDDSAPub.X(),
				Y:     rawKeyEddsa.EDDSAPub.Y(),
			}

			publicKeyStr := base58.Encode(pk.Serialize())

			getAddressResponse := GetAddressResponse{
				Address: publicKeyStr,
			}
			err := json.NewEncoder(w).Encode(getAddressResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else if keyCurve == SECP256K1_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)

			x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
			y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())

			publicKeyStr := "04" + x + y
			publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
			mainnetAddress, testnetAddress, regtestAddress := publicKeyToBitcoinAddresses(publicKeyBytes)

			getBitcoinAddressesResponse := GetBitcoinAddressesResponse{
				MainnetAddress: mainnetAddress,
				TestnetAddress: testnetAddress,
				RegtestAddress: regtestAddress,
			}
			err := json.NewEncoder(w).Encode(getBitcoinAddressesResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else if keyCurve == APTOS_EDDSA_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

			pk := edwards.PublicKey{
				Curve: tss.Edwards(),
				X:     rawKeyEddsa.EDDSAPub.X(),
				Y:     rawKeyEddsa.EDDSAPub.Y(),
			}

			publicKeyStr := hex.EncodeToString(pk.Serialize())

			getAddressResponse := GetAddressResponse{
				Address: "0x" + publicKeyStr,
			}
			err := json.NewEncoder(w).Encode(getAddressResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else if keyCurve == STELLAR_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

			pk := edwards.PublicKey{
				Curve: tss.Edwards(),
				X:     rawKeyEddsa.EDDSAPub.X(),
				Y:     rawKeyEddsa.EDDSAPub.Y(),
			}

			// Get the public key bytes
			pkBytes := pk.Serialize()

			// Stellar StrKey format:
			if len(pkBytes) != 32 {
				http.Error(w, "Invalid public key length", http.StatusInternalServerError)
				return
			}

			// Version byte for ED25519 public key in Stellar
			versionByte := strkey.VersionByteAccountID // 6 << 3, or 48

			// Use Stellar SDK's strkey package to encode
			address, err := strkey.Encode(versionByte, pkBytes)
			if err != nil {
				http.Error(w, fmt.Sprintf("error encoding Stellar address: %v", err), http.StatusInternalServerError)
				return
			}

			getAddressResponse := GetAddressResponse{
				Address: address,
			}
			err = json.NewEncoder(w).Encode(getAddressResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else if keyCurve == ALGORAND_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

			pk := edwards.PublicKey{
				Curve: tss.Edwards(),
				X:     rawKeyEddsa.EDDSAPub.X(),
				Y:     rawKeyEddsa.EDDSAPub.Y(),
			}

			// Get the public key bytes
			pkBytes := pk.Serialize()

			// Convert to Algorand address format
			// Algorand addresses are the last 32 bytes of the SHA512_256 of the public key
			hasher := sha512.New512_256()
			hasher.Write(pkBytes)
			checksum := hasher.Sum(nil)[28:] // Last 4 bytes
			// Add the prefix 'a' for Algorand address
			// Concatenate public key and checksum
			addressBytes := append(pkBytes, checksum...)

			// Encode in base32 without padding
			address := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(addressBytes)

			getAddressResponse := GetAddressResponse{
				Address: address,
			}
			err := json.NewEncoder(w).Encode(getAddressResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else {
			json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)

			x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
			y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())

			publicKeyStr := "04" + x + y
			publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
			address := publicKeyToAddress(publicKeyBytes)

			getAddressResponse := GetAddressResponse{
				Address: address,
			}
			err := json.NewEncoder(w).Encode(getAddressResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}

			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		}
	})

	http.HandleFunc("/signature", func(w http.ResponseWriter, r *http.Request) {
		// the owner of the wallet must have created an intent and signed it.
		// we generate signature for an intent operation

		var intent sequencer.Intent

		err := json.NewDecoder(r.Body).Decode(&intent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		operationIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))
		operationIndexInt := uint(operationIndex)

		if intent.Expiry < uint64(time.Now().Unix()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		msg := ""

		if intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_TRANSACTION {
			msg = intent.Operations[operationIndexInt].DataToSign
		} else if intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_SOLVER {
			intentBytes, err := json.Marshal(intent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			res, err := solver.Construct(intent.Operations[operationIndexInt].Solver, &intentBytes, int(operationIndexInt))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			msg = res
		} else if intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_BRIDGE_DEPOSIT {
			// validate if the DataToSign is actually correct by decoding the previous operation
			msg = intent.Operations[operationIndexInt].SolverDataToSign
		} else if intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_SWAP {
			msg = intent.Operations[operationIndexInt].SolverDataToSign
		} else if intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_BURN {
			msg = intent.Operations[operationIndexInt].SolverDataToSign
		} else if intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_WITHDRAW {
			msg = intent.Operations[operationIndexInt].SolverDataToSign
		}

		identity := intent.Identity
		identityCurve := intent.IdentityCurve
		keyCurve := intent.Operations[operationIndexInt].KeyCurve

		// verify signature
		intentStr, err := identityVerification.SanitiseIntent(intent)
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			return
		}

		verified, err := identityVerification.VerifySignature(
			intent.Identity,
			intent.IdentityCurve,
			intentStr,
			intent.Signature,
		)

		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			return
		}

		if !verified {
			http.Error(w, "signature verification failed", http.StatusBadRequest)
			return
		}

		if keyCurve == EDDSA_CURVE {
			msgBytes, err := base58.Decode(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
				return
			}
			go generateSignatureMessage(identity, identityCurve, keyCurve, msgBytes)
		} else if keyCurve == ECDSA_CURVE {
			if intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_BRIDGE_DEPOSIT ||
				intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_SWAP ||
				intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_BURN ||
				intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_WITHDRAW {
				go generateSignatureMessage(BridgeContractAddress, "ecdsa", "ecdsa", []byte(msg))
			} else {
				go generateSignatureMessage(identity, identityCurve, keyCurve, []byte(msg))
			}
		} else if keyCurve == SECP256K1_CURVE {
			if intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_BRIDGE_DEPOSIT ||
				intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_SWAP ||
				intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_BURN ||
				intent.Operations[operationIndexInt].Type == sequencer.OPERATION_TYPE_WITHDRAW {
				// For secp256k1, we need to hash the message first
				msgHash := crypto.Keccak256([]byte(msg))
				go generateSignatureMessage(BridgeContractAddress, "secp256k1", "secp256k1", msgHash)
			} else {
				msgHash := crypto.Keccak256([]byte(msg))
				go generateSignatureMessage(identity, identityCurve, keyCurve, msgHash)
			}
		} else if keyCurve == APTOS_EDDSA_CURVE {
			go generateSignatureMessage(identity, identityCurve, keyCurve, []byte(msg))
		} else if keyCurve == STELLAR_CURVE {
			msgBytes, err := base64.StdEncoding.DecodeString(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("error decoding Stellar message: %v", err), http.StatusInternalServerError)
				return
			}
			go generateSignatureMessage(identity, identityCurve, keyCurve, msgBytes)
		} else if keyCurve == ALGORAND_CURVE {
			// For Algorand, decode the base32 message first
			// msgBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(msg)
			msgBytes, err := base64.StdEncoding.DecodeString(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("error decoding Algorand message: %v", err), http.StatusInternalServerError)
				return
			}
			go generateSignatureMessage(identity, identityCurve, keyCurve, msgBytes)
		} else {
			http.Error(w, "invalid key curve", http.StatusBadRequest)
			return
		}

		// Create a channel using the message as the key. The key format varies by chain:
		// - Solana: base58 encoded string (from client)
		// - Bitcoin/Aptos: hex encoded string
		// - Algorand: base64 encoded string (from client)
		// - Stellar: base32 encoded string (from client)
		// - Ethereum: raw bytes as string
		// This same format must be used in message.go when looking up the channel
		messageChan[msg] = make(chan Message)

		// Wait for the signature to be sent back through the channel
		sig := <-messageChan[msg]

		w.Header().Set("Content-Type", "application/json")

		signatureResponse := SignatureReponse{}

		if keyCurve == ECDSA_CURVE {
			signatureResponse.Signature = string(sig.Message)
			signatureResponse.Address = sig.Address
		} else if keyCurve == SECP256K1_CURVE {
			signatureResponse.Signature = hex.EncodeToString(sig.Message)
			signatureResponse.Address = sig.Address
		} else if keyCurve == APTOS_EDDSA_CURVE || keyCurve == STELLAR_CURVE {
			signatureResponse.Signature = hex.EncodeToString(sig.Message)
			fmt.Println("generated signature", hex.EncodeToString(sig.Message))
			signatureResponse.Address = sig.Address
		} else if keyCurve == ALGORAND_CURVE {
			// For Algorand, encode the signature in base64 (Algorand's standard)
			signatureResponse.Signature = base64.StdEncoding.EncodeToString(sig.Message)
			signatureResponse.Address = sig.Address
			v, err := identityVerification.VerifySignature(sig.Address, "algorand_eddsa", msg, signatureResponse.Signature)
			if !v {
				http.Error(w, fmt.Sprintf("error verifying algorand signature: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			signatureResponse.Signature = base58.Encode(sig.Message)
			signatureResponse.Address = sig.Address
		}

		err = json.NewEncoder(w).Encode(signatureResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		}

		delete(messageChan, msg)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

type SignatureReponse struct {
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

type GetAddressResponse struct {
	Address string `json:"address"`
}

type GetBitcoinAddressesResponse struct {
	MainnetAddress string `json:"mainnetAddress"`
	TestnetAddress string `json:"testnetAddress"`
	RegtestAddress string `json:"regtestAddress"`
}
