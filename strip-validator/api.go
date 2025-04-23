package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/dogecoin"
	identityVerification "github.com/StripChain/strip-node/identity"
	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/ripple"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/stellar/go/strkey"

	"github.com/StripChain/strip-node/bitcoin"
	"github.com/StripChain/strip-node/solver"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/mr-tron/base58"
)

type BurnMetadata struct {
	Token string `json:"token"`
}

// BurnSyntheticMetadata defines the metadata required for the BURN_SYNTHETIC operation
type BurnSyntheticMetadata struct {
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

// BridgeDepositMetadata defines the metadata required for bridgeDeposit operations
type BridgeDepositMetadata struct {
	BlockchainID blockchains.BlockchainID `json:"blockchainID"`
	Result       string                   `json:"result"` // The transaction hash/result
	Token        string                   `json:"token"`  // Token address
}

type WithdrawMetadata struct {
	Token  string `json:"token"`
	Unlock bool   `json:"unlock"`
}

var messageChan = make(map[string]chan (Message))
var keygenGeneratedChan = make(map[string]chan (string))

var (
	ECDSA_CURVE       = "ecdsa"
	EDDSA_CURVE       = "eddsa"
	APTOS_EDDSA_CURVE = "aptos_eddsa"
	BITCOIN_CURVE     = "bitcoin_ecdsa"
	DOGECOIN_CURVE    = "dogecoin_ecdsa"
	SUI_EDDSA_CURVE   = "sui_eddsa"     // Sui uses Ed25519 for native transactions
	STELLAR_CURVE     = "stellar_eddsa" // Stellar uses Ed25519 with StrKey encoding
	ALGORAND_CURVE    = "algorand_eddsa"
	// Note: Hedera uses ECDSA_CURVE since it's compatible with EVM
	RIPPLE_CURVE  = "ripple_eddsa" // Ripple supports Ed25519 https://xrpl.org/docs/concepts/accounts/cryptographic-keys#signing-algorithms
	CARDANO_CURVE = "cardano_eddsa"
)

func generateKeygenMessage(identity string, identityCurve common.Curve, keyCurve common.Curve, signers []string) {
	message := Message{
		Type:          MESSAGE_TYPE_GENERATE_START_KEYGEN,
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      keyCurve,
		Signers:       signers,
	}

	broadcast(message)
}

func generateSignatureMessage(identity string, blockchainID blockchains.BlockchainID, identityCurve common.Curve, keyCurve common.Curve, msg []byte) {
	message := Message{
		Type:          MESSAGE_TYPE_START_SIGN,
		Hash:          msg,
		Identity:      identity,
		IdentityCurve: identityCurve,
		BlockchainID:  blockchainID,
		KeyCurve:      keyCurve,
	}

	broadcast(message)
}

type CreateWallet struct {
	Identity      string       `json:"identity"`
	IdentityCurve common.Curve `json:"identityCurve"`
	KeyCurve      common.Curve `json:"keyCurve"`
	Signers       []string     `json:"signers"`
}

type SignMessage struct {
	Message       string `json:"message"`
	Identity      string `json:"identity"`
	IdentityCurve string `json:"identityCurve"`
	KeyCurve      string `json:"keyCurve"`
}

func startHTTPServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/keygen", func(w http.ResponseWriter, r *http.Request) {
		requestIP := r.RemoteAddr
		headers := r.Header
		method := r.Method
		contentLength := r.ContentLength

		logger.Sugar().Infow("Received keygen request",
			"method", method,
			"remoteAddr", requestIP,
			"contentLength", contentLength,
			"userAgent", headers.Get("User-Agent"),
			"contentType", headers.Get("Content-Type"))

		var createWallet CreateWallet

		// Read the request body for logging
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Sugar().Errorw("Failed to read request body", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Log the raw request body if it's not too large
		if contentLength > 0 && contentLength < 1024 {
			logger.Sugar().Debugw("Request body", "body", string(bodyBytes))
		}

		// Restore the body for further processing
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err = json.NewDecoder(r.Body).Decode(&createWallet)
		if err != nil {
			logger.Sugar().Errorw("Failed to decode keygen request body", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Sugar().Infow("Processing keygen request",
			"identity", createWallet.Identity,
			"identityCurve", createWallet.IdentityCurve,
			"keyCurve", createWallet.KeyCurve,
			"signersCount", len(createWallet.Signers))

		key := createWallet.Identity + "_" + string(createWallet.IdentityCurve) + "_" + string(createWallet.KeyCurve)

		// Create channel for keygen results with error handling
		keygenGeneratedChan[key] = make(chan string)
		errorChan := make(chan error, 1)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Sugar().Errorw("Panic in keygen operation", "error", fmt.Sprintf("%v", r))
					errorChan <- fmt.Errorf("internal server error: panic in keygen operation")
				}
			}()

			generateKeygenMessage(createWallet.Identity, createWallet.IdentityCurve, createWallet.KeyCurve, createWallet.Signers)
		}()

		logger.Sugar().Infow("Waiting for keygen operation to complete", "key", key)

		// Add timeout to prevent hanging requests
		select {
		case <-keygenGeneratedChan[key]:
			logger.Sugar().Infow("Keygen operation completed successfully",
				"identity", createWallet.Identity,
				"identityCurve", createWallet.IdentityCurve,
				"keyCurve", createWallet.KeyCurve)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Keygen operation completed successfully"))
		case err := <-errorChan:
			logger.Sugar().Errorw("Keygen operation failed",
				"identity", createWallet.Identity,
				"error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case <-time.After(5 * time.Minute): // 5 minute timeout should be sufficient for keygen
			logger.Sugar().Errorw("Keygen operation timed out",
				"identity", createWallet.Identity,
				"identityCurve", createWallet.IdentityCurve,
				"keyCurve", createWallet.KeyCurve)
			http.Error(w, "keygen operation timed out", http.StatusGatewayTimeout)
		}

		delete(keygenGeneratedChan, key)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		peers := h.Network().Peers()
		if len(peers) < 2 {
			http.Error(w, "not enough peers", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")

		identityCurveEnum, err := common.ParseCurve(identityCurve)
		if err != nil {
			http.Error(w, "invalid identity curve", http.StatusBadRequest)
			return
		}

		var rawKeyEddsa *eddsaKeygen.LocalPartySaveData
		var rawKeyEcdsa *ecdsaKeygen.LocalPartySaveData

		addressesResponse := libs.AddressesResponse{
			Addresses: make(map[blockchains.BlockchainID]map[blockchains.NetworkType]string),
		}
		for _, keyCurve := range []common.Curve{common.CurveEcdsa, common.CurveEddsa} {
			keyShare, err := GetKeyShare(identity, identityCurveEnum, keyCurve)

			if err != nil {
				http.Error(w, "error from postgres", http.StatusBadRequest)
				return
			}

			if keyShare == "" {
				http.Error(w, "key share not found.", http.StatusBadRequest)
				return
			}

			for _, blockchainID := range blockchains.GetRegisteredBlockchains() {
				opBlockchain, err := blockchains.GetBlockchain(blockchainID, blockchains.NetworkType(blockchains.Mainnet))
				if err != nil {
					opBlockchain, err = blockchains.GetBlockchain(blockchainID, blockchains.NetworkType(blockchains.Testnet))
					if err != nil {
						logger.Sugar().Errorw("error getting blockchain", "error", err)
						http.Error(w, "invalid blockchain ID", http.StatusBadRequest)
						return
					}
				}
				if opBlockchain.KeyCurve() != keyCurve {
					continue
				}
				addressesResponse.Addresses[blockchainID] = make(map[blockchains.NetworkType]string)
				switch blockchainID {
				case blockchains.Solana:
					json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

					pk := edwards.PublicKey{
						Curve: tss.Edwards(),
						X:     rawKeyEddsa.EDDSAPub.X(),
						Y:     rawKeyEddsa.EDDSAPub.Y(),
					}

					publicKeyStr := base58.Encode(pk.Serialize())

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = publicKeyStr
				case blockchains.Bitcoin:
					json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)

					xStr := fmt.Sprintf("%064x", rawKeyEcdsa.ECDSAPub.X())
					prefix := "02"
					if rawKeyEcdsa.ECDSAPub.Y().Bit(0) == 1 {
						prefix = "03"
					}
					publicKeyStr := prefix + xStr
					publicKeyBytes, err := hex.DecodeString(publicKeyStr)
					if err != nil {
						http.Error(w, fmt.Sprintf("error decoding public key, %v", err), http.StatusBadRequest)
						return
					}
					mainnetAddress, testnetAddress, regtestAddress := bitcoin.PublicKeyToBitcoinAddresses(publicKeyBytes)

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = mainnetAddress
					addressesResponse.Addresses[blockchainID][blockchains.Testnet] = testnetAddress
					addressesResponse.Addresses[blockchainID][blockchains.Regnet] = regtestAddress
				case blockchains.Dogecoin:
					json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)

					x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
					y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())

					publicKeyStr := "04" + x + y

					mainnetAddress, err := dogecoin.PublicKeyToAddress(publicKeyStr)
					if err != nil {
						http.Error(w, fmt.Sprintf("error generating Dogecoin mainnet address: %v", err), http.StatusInternalServerError)
						return
					}

					testnetAddress, err := dogecoin.PublicKeyToTestnetAddress(publicKeyStr)
					if err != nil {
						http.Error(w, fmt.Sprintf("error generating Dogecoin testnet address: %v", err), http.StatusInternalServerError)
						return
					}

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = mainnetAddress
					addressesResponse.Addresses[blockchainID][blockchains.Testnet] = testnetAddress
				case blockchains.Sui:
					json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

					pk := edwards.PublicKey{
						Curve: tss.Edwards(),
						X:     rawKeyEddsa.EDDSAPub.X(),
						Y:     rawKeyEddsa.EDDSAPub.Y(),
					}

					// Serialize the Ed25519 public key
					pkBytes := pk.Serialize()

					// Full public key in hex
					// publicKeyHex := hex.EncodeToString(pkBytes)

					// Hash the public key with Blake2b-256 to get Sui address
					// hasher := blake2b.Sum256(pkBytes)
					// suiAddress := "0x" + hex.EncodeToString(hasher[:])
					flag := byte(0x00)
					hasher, _ := blake2b.New256(nil)
					hasher.Write([]byte{flag})
					hasher.Write(pkBytes)

					arr := hasher.Sum(nil)
					suiAddress := "0x" + hex.EncodeToString(arr)

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = suiAddress
				case blockchains.Aptos:
					json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

					pk := edwards.PublicKey{
						Curve: tss.Edwards(),
						X:     rawKeyEddsa.EDDSAPub.X(),
						Y:     rawKeyEddsa.EDDSAPub.Y(),
					}

					publicKeyStr := hex.EncodeToString(pk.Serialize())

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = "0x" + publicKeyStr
				case blockchains.Stellar:
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

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = address
				case blockchains.Algorand:
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

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = address
				case blockchains.Ripple:
					json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = ripple.PublicKeyToAddress(rawKeyEddsa)
				case blockchains.Cardano:
					json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

					// Get the public key
					pk := edwards.PublicKey{
						Curve: rawKeyEddsa.EDDSAPub.Curve(),
						X:     rawKeyEddsa.EDDSAPub.X(),
						Y:     rawKeyEddsa.EDDSAPub.Y(),
					}

					publicKeyStr := hex.EncodeToString(pk.Serialize())

					addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = publicKeyStr
				default:
					if blockchains.IsEVMBlockchain(blockchainID) {
						json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)

						x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
						y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())

						publicKeyStr := "04" + x + y
						publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
						address := publicKeyToAddress(publicKeyBytes)

						addressesResponse.Addresses[blockchainID][blockchains.Mainnet] = address
					} else {

						logger.Sugar().Errorw("unsupported blockchain ID", "blockchainID", blockchainID)
						http.Error(w, "unsupported blockchain ID", http.StatusBadRequest)
						return
					}
				}
			}
		}

		pk := edwards.PublicKey{
			Curve: tss.Edwards(),
			X:     rawKeyEddsa.EDDSAPub.X(),
			Y:     rawKeyEddsa.EDDSAPub.Y(),
		}

		addressesResponse.EDDSAAddress = hex.EncodeToString(pk.Serialize())

		x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
		y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())
		publicKeyStr := "04" + x + y
		publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
		addressesResponse.ECDSAAddress = publicKeyToAddress(publicKeyBytes)
		err = json.NewEncoder(w).Encode(addressesResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		}

	})

	http.HandleFunc("/signature", func(w http.ResponseWriter, r *http.Request) {
		// Set content type header for all responses
		w.Header().Set("Content-Type", "application/json")

		// the owner of the wallet must have created an intent and signed it.
		// we generate signature for an intent operation

		var intent libs.Intent

		err := json.NewDecoder(r.Body).Decode(&intent)
		if err != nil {
			logger.Sugar().Errorw("Failed to decode intent", "error", err)
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusBadRequest)
			return
		}

		operationIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))
		operationIndexInt := uint(operationIndex)

		// Validate intent has operations and the requested index is valid
		if int(operationIndexInt) >= len(intent.Operations) {
			errMsg := fmt.Sprintf("Invalid operation index: %d (intent has %d operations)", operationIndexInt, len(intent.Operations))
			logger.Sugar().Errorw(errMsg)
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", errMsg), http.StatusBadRequest)
			return
		}

		if intent.Expiry.Before(time.Now()) {
			logger.Sugar().Errorw("Intent has expired", "expiryTime", intent.Expiry, "currentTime", time.Now().Unix())
			http.Error(w, "{\"error\":\"Intent has expired\"}", http.StatusBadRequest)
			return
		}

		logger.Sugar().Infow("Processing signature request",
			"intentID", intent.ID,
			"operationIndex", operationIndexInt,
			"operationType", intent.Operations[operationIndexInt].Type)

		msg := ""

		intentBlockchain, err := blockchains.GetBlockchain(intent.BlockchainID, blockchains.NetworkType(intent.NetworkType))
		if err != nil {
			logger.Sugar().Errorw("error getting blockchain", "error", err)
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
			return
		}

		opBlockchain, err := blockchains.GetBlockchain(intent.Operations[operationIndexInt].BlockchainID, intent.Operations[operationIndexInt].NetworkType)
		if err != nil {
			logger.Sugar().Errorw("error getting blockchain", "error", err)
			return
		}

		operation := intent.Operations[operationIndexInt]
		switch operation.Type {
		case libs.OperationTypeTransaction:
			msg = ""
			if operation.DataToSign != nil {
				msg = *operation.DataToSign
			}
		case libs.OperationTypeSendToBridge:
			// Verify only operation for bridging
			// Get bridgewallet by calling /getwallet from sequencer api
			// req, err := http.NewRequest("GET", SequencerHost+"/getWallet?identity="+intent.Identity+"&identityCurve="+intent.IdentityCurve, nil)
			req, err := http.NewRequest("GET", SequencerHost+"/getBridgeAddress", nil)
			// req, err := http.NewRequest("GET", fmt.Sprintf("%s/getWallet?identity=%s&blockchain=%s", SequencerHost, intent.Identity, intent.BlockchainID), nil)
			if err != nil {
				logger.Sugar().Errorw("error creating request", "error", err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logger.Sugar().Errorw("error sending request", "error", err)
				return
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Sugar().Errorw("error reading response body", "error", err)
				return
			}

			var bridgeWallet db.WalletSchema
			err = json.Unmarshal(body, &bridgeWallet)
			if err != nil {
				logger.Sugar().Errorw("error unmarshalling response body", "error", err)
				return
			}

			// if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
			// 	chain, err := common.GetChain(operation.ChainId)
			// 	if err != nil {
			// 		logger.Sugar().Errorw("error getting chain", "error", err)
			// 		return
			// 	}

			// 	// Extract destination address from serialized transaction
			// 	var destAddress string
			// 	if chain.ChainType == "bitcoin" || chain.ChainType == "dogecoin" {
			// 		// For Bitcoin, decode the serialized transaction to get output address
			// 		var tx wire.MsgTx
			// 		txBytes, err := hex.DecodeString(operation.SerializedTxn)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("error decoding bitcoin&dogecoin transaction", "error", err)
			// 			return
			// 		}
			// 		if err := tx.Deserialize(bytes.NewReader(txBytes)); err != nil {
			// 			logger.Sugar().Errorw("error deserializing bitcoin&dogecoin transaction", "error", err)
			// 			return
			// 		}
			// 		// Get the first output's address (assuming it's the bridge address)
			// 		if len(tx.TxOut) > 0 {
			// 			_, addrs, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[0].PkScript, nil)
			// 			if err != nil || len(addrs) == 0 {
			// 				logger.Sugar().Errorw("error extracting bitcoin&dogecoin address", "error", err)
			// 				return
			// 			}
			// 			destAddress = addrs[0].String()
			// 		}
			// 	} else {
			// 		// For EVM chains, decode the transaction to get the 'to' address
			// 		txBytes, err := hex.DecodeString(operation.SerializedTxn)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("error decoding EVM transaction", "error", err)
			// 			return
			// 		}
			// 		tx := new(types.Transaction)
			// 		if err := rlp.DecodeBytes(txBytes, tx); err != nil {
			// 			logger.Sugar().Errorw("error deserializing EVM transaction", "error", err)
			// 			return
			// 		}
			// 		// Check if this is an ERC20 transfer (function signature: transfer(address,uint256))
			// 		// ERC20 transfer function signature is: 0xa9059cbb
			// 		isERC20Transfer := false
			// 		contractAddress := tx.To().Hex()
			// 		if tx.Data() != nil && len(tx.Data()) >= 4 && hex.EncodeToString(tx.Data()[:4]) == "a9059cbb" {
			// 			// This is an ERC20 transfer - extract the recipient address from the data
			// 			// The recipient address is the first parameter of the transfer function (32-byte padded)
			// 			isERC20Transfer = true
			// 			if len(tx.Data()) >= 36 {
			// 				// Extract the recipient address from position 4:36 (32 bytes)
			// 				recipientBytes := tx.Data()[4:36]
			// 				// Convert to address (take last 20 bytes for proper Ethereum address length)
			// 				destAddress = "0x" + hex.EncodeToString(recipientBytes[12:])
			// 				logger.Sugar().Infow("detected ERC20 transfer", "contract", contractAddress, "recipient", destAddress)
			// 			} else {
			// 				logger.Sugar().Errorw("invalid ERC20 transfer data length", "data", hex.EncodeToString(tx.Data()))
			// 				return
			// 			}
			// 		} else {
			// 			// For regular ETH transfers, use the 'to' address directly
			// 			destAddress = contractAddress
			// 		}

			// 		// Verify destination address matches bridge wallet
			// 		var expectedAddress string
			// 		if chain.ChainType == "bitcoin" {
			// 			expectedAddress = bridgeWallet.BitcoinMainnetPublicKey
			// 		} else if chain.ChainType == "dogecoin" {
			// 			expectedAddress = bridgeWallet.DogecoinMainnetPublicKey
			// 		} else {
			// 			// For Ethereum chains, determine if we need to use the bridge contract address
			// 			if chain.ChainType == "evm" && !isERC20Transfer {
			// 				// For native ETH transfers, get the bridge address from the dedicated endpoint
			// 				bridgeReq, err := http.NewRequest("GET", SequencerHost+"/getBridgeAddress", nil)
			// 				if err != nil {
			// 					logger.Sugar().Errorw("error creating bridge address request", "error", err)
			// 					return
			// 				}

			// 				bridgeReq.Header.Set("Content-Type", "application/json")
			// 				bridgeClient := &http.Client{}
			// 				bridgeResp, err := bridgeClient.Do(bridgeReq)
			// 				if err != nil {
			// 					logger.Sugar().Errorw("error fetching bridge address", "error", err)
			// 					return
			// 				}

			// 				defer bridgeResp.Body.Close()

			// 				bridgeBody, err := io.ReadAll(bridgeResp.Body)
			// 				if err != nil {
			// 					logger.Sugar().Errorw("error reading bridge address response", "error", err)
			// 					return
			// 				}

			// 				var bridgeAddressWallet db.WalletSchema
			// 				err = json.Unmarshal(bridgeBody, &bridgeAddressWallet)
			// 				if err != nil {
			// 					logger.Sugar().Errorw("error unmarshalling bridge address response", "error", err)
			// 					return
			// 				}

			// 				expectedAddress = bridgeAddressWallet.ECDSAPublicKey
			// 				logger.Sugar().Infow("using bridge contract address for native ETH transfer", "address", expectedAddress)
			// 			} else {
			// 				expectedAddress = bridgeWallet.ECDSAPublicKey
			// 			}
			// 		}

			// 		// Verify the extracted destination matches the bridge wallet
			// 		if !strings.EqualFold(destAddress, expectedAddress) {
			// 			logger.Sugar().Errorw("Invalid bridge destination address", "expected", expectedAddress, "got", destAddress)
			// 			return
			// 		}
			// 	}
			// } else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" || operation.KeyCurve == "cardano_eddsa" || operation.KeyCurve == "sui_eddsa" {
			// 	chain, err := common.GetChain(operation.ChainId)
			// 	if err != nil {
			// 		logger.Sugar().Errorw("error getting chain", "error", err)
			// 		return
			// 	}

			// 	// Verify destination address matches bridge wallet based on chain type
			// 	var validDestination bool
			// 	var destAddress string

			// 	// Extract destination address from serialized transaction based on chain type
			// 	switch chain.ChainType {
			// 	case "solana":
			// 		// Decode base58 transaction and extract destination
			// 		decodedTxn, err := base58.Decode(operation.SerializedTxn)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("error decoding Solana transaction", "error", err)
			// 			return
			// 		}
			// 		tx, err := solanasdk.TransactionFromDecoder(bin.NewBinDecoder(decodedTxn))
			// 		if err != nil || len(tx.Message.Instructions) == 0 {
			// 			logger.Sugar().Errorw("error deserializing Solana transaction", "error", err)
			// 			return
			// 		}
			// 		// Get the first instruction's destination account index
			// 		destAccountIndex := tx.Message.Instructions[0].Accounts[1]
			// 		// Get the actual account address from the message accounts
			// 		destAddress = tx.Message.AccountKeys[destAccountIndex].String()
			// 	case "aptos":
			// 		// For Aptos, the destination is in the transaction payload
			// 		var aptosPayload struct {
			// 			Function string   `json:"function"`
			// 			Args     []string `json:"arguments"`
			// 		}
			// 		if err := json.Unmarshal([]byte(operation.SerializedTxn), &aptosPayload); err != nil {
			// 			logger.Sugar().Errorw("error parsing Aptos transaction", "error", err)
			// 			return
			// 		}
			// 		if len(aptosPayload.Args) > 0 {
			// 			destAddress = aptosPayload.Args[0] // First arg is typically the recipient
			// 		}
			// 	case "stellar":
			// 		// For Stellar, parse the XDR transaction envelope
			// 		var txEnv xdr.TransactionEnvelope
			// 		err := xdr.SafeUnmarshalBase64(operation.SerializedTxn, &txEnv)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("error parsing Stellar transaction", "error", err)
			// 			return
			// 		}

			// 		// Get the first operation's destination
			// 		if len(txEnv.Operations()) > 0 {
			// 			if paymentOp, ok := txEnv.Operations()[0].Body.GetPaymentOp(); ok {
			// 				destAddress = paymentOp.Destination.Address()
			// 			}
			// 		}
			// 	case "algorand":
			// 		txnBytes, err := base64.StdEncoding.DecodeString(operation.SerializedTxn)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("failed to decode serialized transaction", "error", err)
			// 			return
			// 		}
			// 		var txn algorandTypes.Transaction
			// 		err = msgpack.Decode(txnBytes, &txn)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("failed to deserialize transaction", "error", err)
			// 			return
			// 		}
			// 		if txn.Type == algorandTypes.PaymentTx {
			// 			destAddress = txn.PaymentTxnFields.Receiver.String()
			// 		} else if txn.Type == algorandTypes.AssetTransferTx {
			// 			destAddress = txn.AssetTransferTxnFields.AssetReceiver.String()
			// 		} else {
			// 			logger.Sugar().Errorw("Unknown transaction type", "type", txn.Type)
			// 			return
			// 		}
			// 	case "ripple":
			// 		// For Ripple, the destination is in the transaction payload
			// 		// Decode the serialized transaction
			// 		txBytes, err := hex.DecodeString(strings.TrimPrefix(operation.SerializedTxn, "0x"))
			// 		if err != nil {
			// 			logger.Sugar().Errorw("error decoding transaction", "error", err)
			// 			return
			// 		}

			// 		// Parse the transaction
			// 		var tx data.Payment
			// 		err = json.Unmarshal(txBytes, &tx)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("error unmarshalling transaction", "error", err)
			// 			return
			// 		}
			// 		destAddress = tx.Destination.String()
			// 	case "cardano":
			// 		var tx cardanolib.Tx
			// 		txBytes, err := hex.DecodeString(operation.SerializedTxn)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("error decoding Cardano transaction", "error", err)
			// 			return
			// 		}
			// 		if err := json.Unmarshal(txBytes, &tx); err != nil {
			// 			logger.Sugar().Errorw("error parsing Cardano transaction", "error", err)
			// 			return
			// 		}
			// 		destAddress = tx.Body.Outputs[0].Address.String()
			// 	case "sui":
			// 		var tx sui_types.TransactionData
			// 		txBytes, err := base64.StdEncoding.DecodeString(operation.SerializedTxn)
			// 		if err != nil {
			// 			logger.Sugar().Errorw("error decoding Sui transaction", "error", err)
			// 			return
			// 		}
			// 		if err := json.Unmarshal(txBytes, &tx); err != nil {
			// 			logger.Sugar().Errorw("error parsing Sui transaction", "error", err)
			// 			return
			// 		}
			// 		if len(tx.V1.Kind.ProgrammableTransaction.Inputs) < 1 {
			// 			logger.Sugar().Errorw("wrong format sui transaction", "error", err)
			// 			return
			// 		}
			// 		destAddress = string(*tx.V1.Kind.ProgrammableTransaction.Inputs[0].Pure)
			// 	}

			// 	// Verify the extracted destination matches the bridge wallet
			// 	if destAddress == "" {
			// 		logger.Sugar().Errorw("Failed to extract destination address from %s transaction", chain.ChainType)
			// 		validDestination = false
			// 	} else {
			// 		switch chain.ChainType {
			// 		case "solana":
			// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.EDDSAPublicKey)
			// 		case "aptos":
			// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.AptosEDDSAPublicKey)
			// 		case "stellar":
			// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.StellarPublicKey)
			// 		case "algorand":
			// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.AlgorandEDDSAPublicKey)
			// 		case "ripple":
			// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.RippleEDDSAPublicKey)
			// 		case "cardano":
			// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.CardanoPublicKey)
			// 		case "sui":
			// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.SuiPublicKey)
			// 		}
			// 	}

			// 	if !validDestination {
			// 		logger.Sugar().Errorw("Invalid bridge destination address for", "chain", chain.ChainType)
			// 		return
			// 	}
			// }

			// Set message
			msg = ""
			if operation.DataToSign != nil {
				msg = *operation.DataToSign
			}
		case libs.OperationTypeSolver:
			intentBytes, err := json.Marshal(intent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			res, err := solver.Construct(operation.Solver, &intentBytes, int(operationIndexInt))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			msg = res
		case libs.OperationTypeBridgeDeposit:
			// For bridgeDeposit operations, extract transaction details from metadata
			logger.Sugar().Infow("Processing bridgeDeposit operation",
				"intentID", intent.ID,
				"operationIndex", operationIndexInt,
				"solverMetadata", operation.SolverMetadata)

			// Create depositOperation variable at the outer scope
			var depositOperation libs.Operation
			var needPreviousOp bool
			var tokenAddress string

			// Check if metadata is empty
			if operation.SolverMetadata == "" {
				logger.Sugar().Warnw("Bridge deposit has empty metadata, falling back to previous operation")
				needPreviousOp = true
			} else {
				// Extract the transaction hash from the operation metadata
				var metadata BridgeDepositMetadata
				err := json.Unmarshal([]byte(operation.SolverMetadata), &metadata)
				if err != nil {
					logger.Sugar().Errorw("Error unmarshalling bridge deposit metadata",
						"error", err,
						"solverMetadata", operation.SolverMetadata)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bridge deposit metadata"})
					return
				}

				logger.Sugar().Infow("Parsed bridge deposit metadata",
					"blockchainID", metadata.BlockchainID,
					"result", metadata.Result,
					"token", metadata.Token)

				// If metadata has no Result field, we need to get transaction details from previous operation
				if metadata.Result == "" {
					logger.Sugar().Warnw("No transaction hash in bridge deposit metadata, falling back to previous operation",
						"token", metadata.Token)
					needPreviousOp = true

					// Save token address for later use if present
					tokenAddress = metadata.Token
				} else {
					// Use the transaction result from metadata
					depositOperation = libs.Operation{
						BlockchainID: metadata.BlockchainID,
						Result:       metadata.Result,
					}

					logger.Sugar().Infow("Extracted bridge deposit transaction details",
						"blockchainID", depositOperation.BlockchainID,
						"txHash", depositOperation.Result)
				}
			}

			// If we need transaction details from previous operation
			if needPreviousOp {
				// If metadata is empty, try to get transaction details from previous operation
				if operationIndexInt == 0 {
					logger.Sugar().Errorw("No previous operation to get transaction details from")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Missing transaction information"})
					return
				}

				// Get the previous operation which should be sendToBridge
				prevOp := intent.Operations[operationIndexInt-1]
				if prevOp.Type != libs.OperationTypeSendToBridge {
					logger.Sugar().Errorw("Previous operation is not sendToBridge",
						"prevOpType", prevOp.Type,
						"prevOpIndex", operationIndexInt-1)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Previous operation is not sendToBridge"})
					return
				}

				// Use transaction details from previous operation
				depositOperation = libs.Operation{
					BlockchainID: prevOp.BlockchainID,
					Result:       prevOp.Result,
				}

				logger.Sugar().Infow("Using transaction details from previous operation",
					"blockchainID", depositOperation.BlockchainID,
					"txHash", depositOperation.Result,
					"savedToken", tokenAddress)
			}

			logger.Sugar().Infow("Processing bridgeDeposit for chain",
				"blockchainID", depositOperation.BlockchainID,
				"txHash", depositOperation.Result)

			transfers, err := opBlockchain.GetTransfers(depositOperation.Result, &intent.Identity)
			if err != nil {
				logger.Sugar().Errorw("error getting transfers", "error", err)
				return
			}

			if len(transfers) == 0 {
				// If we have a token address from metadata, create a minimal transfer to proceed
				if tokenAddress != "" {
					logger.Sugar().Warnw("No transfers found but token address provided in metadata, creating minimal transfer",
						"tokenAddress", tokenAddress,
						"txHash", depositOperation.Result)

					// Create a minimal transfer with the token address from metadata
					transfers = append(transfers, common.Transfer{
						TokenAddress: tokenAddress,
						// We don't have amount information, but we need a non-empty array to proceed
						Amount:   "1", // Minimal placeholder amount
						Token:    "",  // Unknown token symbol
						IsNative: false,
						From:     intent.Identity,
						To:       "", // Unknown recipient
					})

					logger.Sugar().Infow("Created minimal transfer from metadata token",
						"tokenAddress", tokenAddress,
						"from", intent.Identity)
				} else {
					logger.Sugar().Errorw("No transfers found for bridge deposit",
						"result", depositOperation.Result,
						"identity", intent.Identity,
						"blockchainID", depositOperation.BlockchainID)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "No transfers found in transaction"})
					return
				}
			}

			// check if the token exists
			transfer := transfers[0]
			srcAddress := transfer.TokenAddress

			logger.Sugar().Infow("Validating token for bridge deposit",
				"tokenAddress", srcAddress,
				"tokenSymbol", transfer.Token,
				"amount", transfer.Amount,
				"isNative", transfer.IsNative)

			chainID := opBlockchain.ChainID()
			if chainID == nil {
				logger.Sugar().Errorw("Chain ID is nil", "blockchainID", depositOperation.BlockchainID)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Chain ID is nil"})
				return
			}

			exists, peggedToken, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, *chainID, srcAddress)
			if err != nil {
				logger.Sugar().Errorw("Error checking token existence",
					"error", err,
					"tokenAddress", srcAddress,
					"blockchainID", depositOperation.BlockchainID)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to validate token"})
				return
			}

			if !exists {
				logger.Sugar().Errorw("Token does not exist for bridge deposit",
					"tokenAddress", srcAddress,
					"blockchainID", depositOperation.BlockchainID)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Token does not exist"})
				return
			}

			logger.Sugar().Infow("Token exists for bridge deposit",
				"tokenAddress", srcAddress,
				"peggedToken", peggedToken)

			// Set message for signing - first try SolverDataToSign
			msg = operation.SolverDataToSign

			dataToSign := ""
			if operation.DataToSign != nil {
				dataToSign = *operation.DataToSign
			}
			// Log detailed info about the message being signed
			logger.Sugar().Infow("Processing bridge deposit signature",
				"solverDataLength", len(operation.SolverDataToSign),
				"dataToSignLength", len(dataToSign))

			// If no SolverDataToSign is provided, use DataToSign as fallback
			if len(msg) == 0 {
				logger.Sugar().Infow("Using DataToSign for bridge deposit operation", "length", len(dataToSign))
				msg = dataToSign
			}

			if len(msg) == 0 {
				logger.Sugar().Errorw("No message data available for signing bridge deposit")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "No message data available for signing"})
				return
			}

			logger.Sugar().Infow("Bridge deposit message prepared for signing",
				"msgLength", len(msg),
				"transferToken", transfer.Token,
				"transferAmount", transfer.Amount)
		case libs.OperationTypeSwap:
			// Validate previous operation
			bridgeDeposit := intent.Operations[operationIndexInt-1]

			if operationIndexInt == 0 || !(bridgeDeposit.Type == libs.OperationTypeBridgeDeposit) {
				logger.Sugar().Errorw("Invalid operation type for swap")
				return
			}

			// Set message
			msg = operation.SolverDataToSign
		case libs.OperationTypeBurn:
			// Validate nearby operations
			bridgeSwap := intent.Operations[operationIndex-1]

			if operationIndex+1 >= len(intent.Operations) || intent.Operations[operationIndex+1].Type != libs.OperationTypeWithdraw {
				fmt.Println("BURN operation must be followed by a WITHDRAW operation")
				return
			}

			if operationIndex == 0 || !(bridgeSwap.Type == libs.OperationTypeSwap) {
				logger.Sugar().Errorw("Invalid operation type for swap")
				return
			}

			logger.Sugar().Infow("Burning tokens", "bridgeSwap", bridgeSwap)

			// Set message
			msg = operation.SolverDataToSign
		case libs.OperationTypeBurnSynthetic:
			// This operation allows direct burning of ERC20 tokens from the wallet
			// without requiring a prior swap operation
			burnSyntheticMetadata := BurnSyntheticMetadata{}

			// Verify that this operation is followed by a withdraw operation
			if operationIndex+1 >= len(intent.Operations) || intent.Operations[operationIndex+1].Type != libs.OperationTypeWithdraw {
				logger.Sugar().Errorw("BURN_SYNTHETIC operation must be followed by a WITHDRAW operation")
				return
			}

			err := json.Unmarshal([]byte(operation.SolverMetadata), &burnSyntheticMetadata)
			if err != nil {
				logger.Sugar().Errorw("Error unmarshalling burn synthetic metadata:", "error", err)
				return
			}

			// Get bridgewallet by calling /getwallet from sequencer api
			// req, err := http.NewRequest("GET", SequencerHost+"/getWallet?identity="+intent.Identity+"&identityCurve="+intent.IdentityCurve, nil)
			// req, err := http.NewRequest("GET", SequencerHost+"/getWallet?identity="+intent.Identity+"&identityCurve="+intent.IdentityCurve, nil)
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/getWallet?identity=%s&blockchainID=%s", SequencerHost, intent.Identity, intent.BlockchainID), nil)
			if err != nil {
				logger.Sugar().Errorw("error creating request", "error", err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logger.Sugar().Errorw("error sending request", "error", err)
				return
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Sugar().Errorw("error reading response body", "error", err)
				return
			}

			var bridgeWallet db.WalletSchema
			err = json.Unmarshal(body, &bridgeWallet)
			if err != nil {
				logger.Sugar().Errorw("error unmarshalling response body", "error", err)
				return
			}

			// Verify the user has sufficient token balance
			balance, err := ERC20.GetBalance(RPC_URL, burnSyntheticMetadata.Token, bridgeWallet.ECDSAPublicKey)
			if err != nil {
				logger.Sugar().Errorw("Error getting token balance:", "error", err)
				return
			}

			balanceBig, ok := new(big.Int).SetString(balance, 10)
			if !ok {
				logger.Sugar().Errorw("Error parsing balance")
				return
			}

			amountBig, ok := new(big.Int).SetString(burnSyntheticMetadata.Amount, 10)
			if !ok {
				logger.Sugar().Errorw("Error parsing amount")
				return
			}

			if balanceBig.Cmp(amountBig) < 0 {
				logger.Sugar().Errorw("Insufficient token balance")
				return
			}
			msg = operation.SolverDataToSign
		case libs.OperationTypeWithdraw:
			// Verify nearby operations
			burn := intent.Operations[operationIndex-1]

			if operationIndex == 0 || !(burn.Type == libs.OperationTypeBurn || burn.Type == libs.OperationTypeBurnSynthetic) {
				logger.Sugar().Errorw("Invalid operation type for withdraw after burn")
				return
			}

			var withdrawMetadata WithdrawMetadata
			json.Unmarshal([]byte(operation.SolverMetadata), &withdrawMetadata)

			// Handle different burn operation types
			var tokenToWithdraw string
			var burnTokenAddress string
			if burn.Type == libs.OperationTypeBurn {
				var burnMetadata BurnMetadata
				json.Unmarshal([]byte(burn.SolverMetadata), &burnMetadata)
				tokenToWithdraw = withdrawMetadata.Token
				burnTokenAddress = burnMetadata.Token
			} else if burn.Type == libs.OperationTypeBurnSynthetic {
				var burnSyntheticMetadata BurnSyntheticMetadata
				json.Unmarshal([]byte(burn.SolverMetadata), &burnSyntheticMetadata)
				tokenToWithdraw = withdrawMetadata.Token
				burnTokenAddress = burnSyntheticMetadata.Token
			}

			chainID := opBlockchain.ChainID()
			if chainID == nil {
				logger.Sugar().Errorw("Chain ID is nil", "blockchainID", operation.BlockchainID)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Chain ID is nil"})
				return
			}
			// verify these fields
			exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, *chainID, tokenToWithdraw)

			if err != nil {
				logger.Sugar().Errorw("error checking token existence", "error", err)
				return
			}

			if !exists {
				logger.Sugar().Errorw("Token does not exist", "token", tokenToWithdraw, "blockchainID", operation.BlockchainID)
				return
			}

			if destAddress != burnTokenAddress {
				logger.Sugar().Errorw("Token mismatch", "destAddress", destAddress, "token", burnTokenAddress)
				return
			}

			// Set message
			msg = operation.SolverDataToSign
		}

		identity := intent.Identity

		identityCurve := intentBlockchain.KeyCurve()
		keyCurve := opBlockchain.KeyCurve()
		log.Println("msg", msg)

		// verify signature
		intentStr, err := identityVerification.SanitiseIntent(intent)
		if err != nil {
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
			return
		}

		verified, err := identityVerification.VerifySignature(
			intent.Identity,
			intent.BlockchainID,
			intentStr,
			intent.Signature,
		)

		if err != nil {
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
			return
		}

		if !verified {
			http.Error(w, "{\"error\":\"Signature verification failed\"}", http.StatusBadRequest)
			return
		}

		switch operation.BlockchainID {
		case blockchains.Solana:
			msgBytes, err := base58.Decode(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
				return
			}
			if operation.Type == libs.OperationTypeSwap ||
				operation.Type == libs.OperationTypeBurn ||
				operation.Type == libs.OperationTypeBurnSynthetic ||
				operation.Type == libs.OperationTypeWithdraw {
				go generateSignatureMessage(BridgeContractAddress, blockchains.Arbitrum, common.CurveEcdsa, common.CurveEddsa, msgBytes)
			} else {
				go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, msgBytes)
			}
		case blockchains.Bitcoin:
			go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, []byte(msg))
		case blockchains.Dogecoin:
			go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, []byte(msg))
		case blockchains.Sui:
			// For Sui, we need to format the message according to Sui's standards
			// The message should be prefixed with "Sui Message:" for personal messages
			// suiMsg := []byte("Sui Message:" + msg)
			// go generateSignatureMessage(identity, identityCurve, keyCurve, suiMsg)
			msgBytes, _ := lib.NewBase64Data(msg)
			go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, *msgBytes)
		case blockchains.Stellar:
			msgBytes, err := base64.StdEncoding.DecodeString(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
				return
			}
			go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, msgBytes)
		case blockchains.Algorand:
			// For Algorand, decode the base32 message first
			// msgBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(msg)
			msgBytes, err := base64.StdEncoding.DecodeString(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
				return
			}
			go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, msgBytes)
		case blockchains.Ripple, blockchains.Cardano, blockchains.Aptos:
			msgBytes, err := hex.DecodeString(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
				return
			}
			go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, msgBytes)
		default:
			if blockchains.IsEVMBlockchain(operation.BlockchainID) {
				if operation.Type == libs.OperationTypeBridgeDeposit ||
					operation.Type == libs.OperationTypeSwap ||
					operation.Type == libs.OperationTypeBurn ||
					operation.Type == libs.OperationTypeBurnSynthetic ||
					operation.Type == libs.OperationTypeWithdraw {
					go generateSignatureMessage(BridgeContractAddress, operation.BlockchainID, identityCurve, keyCurve, []byte(msg))
				} else {
					go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, []byte(msg))
				}
			} else {
				http.Error(w, "{\"error\":\"Invalid key curve for signature\"}", http.StatusBadRequest)
				return
			}
		}

		// Create a channel using the message as the key. The key format varies by chain:
		// - Solana: base58 encoded string (from client)
		// - Bitcoin/Aptos/Ripple/Cardano: hex encoded string
		// - Algorand: base64 encoded string (from client)
		// - Stellar: base32 encoded string (from client)
		// - Ethereum: raw bytes as string
		// This same format must be used in message.go when looking up the channel
		messageChan[msg] = make(chan Message)

		// Wait for the signature to be sent back through the channel
		sig := <-messageChan[msg]

		signatureResponse := SignatureResponse{}

		switch operation.BlockchainID {
		case blockchains.Bitcoin:
			signatureResponse.Signature = string(sig.Message)
			signatureResponse.Address = sig.Address
			logger.Sugar().Infof("signatureResponse: %v", signatureResponse)
		case blockchains.Dogecoin:
			// signatureResponse.Signature = hex.EncodeToString(sig.Message)
			signatureResponse.Signature = string(sig.Message)

			signatureResponse.Address = sig.Address
		case blockchains.Sui:
			// For Sui, we return the signature in base64 format
			signatureResponse.Signature = string(sig.Message) // Already base64 encoded in generateSignature
			signatureResponse.Address = sig.Address
			logger.Sugar().Infof("generated Sui signature for address: %s", sig.Message)
		case blockchains.Aptos, blockchains.Stellar, blockchains.Ripple, blockchains.Cardano:
			signatureResponse.Signature = hex.EncodeToString(sig.Message)
			logger.Sugar().Infof("generated signature: %s", hex.EncodeToString(sig.Message))
			signatureResponse.Address = sig.Address
		case blockchains.Algorand:
			// For Algorand, encode the signature in base64 (Algorand's standard)
			signatureResponse.Signature = base64.StdEncoding.EncodeToString(sig.Message)
			signatureResponse.Address = sig.Address
			type algodMsg struct {
				IsRealTransaction bool
				Msg               string
			}
			m := algodMsg{IsRealTransaction: sig.AlgorandFlags.IsRealTransaction, Msg: msg}
			jsonBytes, err := json.Marshal(m)
			if err != nil {
				http.Error(w, fmt.Sprintf("{\"error\":\"Error marshaling algodMsg to JSON: %v\"}", err), http.StatusInternalServerError)
				return
			}
			v, err := identityVerification.VerifySignature(sig.Address, blockchains.Algorand, string(jsonBytes), signatureResponse.Signature)
			if !v {
				logger.Sugar().Errorf("invalid signature %s, err %v", signatureResponse.Signature, err)
			}
		case blockchains.Solana:
			signatureResponse.Signature = base58.Encode(sig.Message)
			signatureResponse.Address = sig.Address
		default:
			if blockchains.IsEVMBlockchain(operation.BlockchainID) {
				signatureResponse.Signature = string(sig.Message)
				signatureResponse.Address = sig.Address
			} else {
				logger.Sugar().Errorw("unsupported blockchain ID", "blockchainID", operation.BlockchainID)
				http.Error(w, "unsupported blockchain ID", http.StatusBadRequest)
				return
			}
		}

		// Validate we have a signature before responding
		if signatureResponse.Signature == "" {
			logger.Sugar().Errorw("Empty signature generated", "blockchainID", operation.BlockchainID, "address", signatureResponse.Address)
			http.Error(w, "{\"error\":\"Failed to generate signature\"}", http.StatusInternalServerError)
			return
		}

		// Log successful signature generation
		logger.Sugar().Infow("Successfully generated signature",
			"blockchainID", operation.BlockchainID,
			"address", signatureResponse.Address,
			"sigLength", len(signatureResponse.Signature))

		// Encode response with proper error handling
		if err := json.NewEncoder(w).Encode(signatureResponse); err != nil {
			logger.Sugar().Errorw("Error encoding signature response", "error", err)
			http.Error(w, fmt.Sprintf("{\"error\":\"Error building the response: %v\"}", err), http.StatusInternalServerError)
			return
		}

		delete(messageChan, msg)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

type SignatureResponse struct {
	Signature string `json:"signature"`
	Address   string `json:"address"`
}
