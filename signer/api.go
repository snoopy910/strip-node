package signer

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/algorand"
	"github.com/StripChain/strip-node/aptos"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/cardano"
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/dogecoin"
	identityVerification "github.com/StripChain/strip-node/identity"
	"github.com/StripChain/strip-node/ripple"
	"github.com/StripChain/strip-node/sequencer"
	"github.com/StripChain/strip-node/stellar"
	"github.com/StripChain/strip-node/sui"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	algorandTypes "github.com/algorand/go-algorand-sdk/types"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	cardanolib "github.com/echovl/cardano-go"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/rubblelabs/ripple/data"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"

	"github.com/StripChain/strip-node/bitcoin"
	"github.com/StripChain/strip-node/solver"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
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
		} else if keyCurve == BITCOIN_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)

			x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
			y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())

			publicKeyStr := "04" + x + y
			publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
			mainnetAddress, testnetAddress, regtestAddress := bitcoin.PublicKeyToBitcoinAddresses(publicKeyBytes)

			getBitcoinAddressesResponse := GetBitcoinAddressesResponse{
				MainnetAddress: mainnetAddress,
				TestnetAddress: testnetAddress,
				RegtestAddress: regtestAddress,
			}
			err := json.NewEncoder(w).Encode(getBitcoinAddressesResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else if keyCurve == DOGECOIN_CURVE {
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

			getDogecoinAddressesResponse := GetDogecoinAddressesResponse{
				MainnetAddress: mainnetAddress,
				TestnetAddress: testnetAddress,
			}
			err = json.NewEncoder(w).Encode(getDogecoinAddressesResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else if keyCurve == SUI_EDDSA_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

			pk := edwards.PublicKey{
				Curve: tss.Edwards(),
				X:     rawKeyEddsa.EDDSAPub.X(),
				Y:     rawKeyEddsa.EDDSAPub.Y(),
			}

			// Serialize the Ed25519 public key
			pkBytes := pk.Serialize()

			// Full public key in hex
			publicKeyHex := hex.EncodeToString(pkBytes)

			// Hash the public key with Blake2b-256 to get Sui address
			// hasher := blake2b.Sum256(pkBytes)
			// suiAddress := "0x" + hex.EncodeToString(hasher[:])
			flag := byte(0x00)
			hasher, _ := blake2b.New256(nil)
			hasher.Write([]byte{flag})
			hasher.Write(pkBytes)

			arr := hasher.Sum(nil)
			suiAddress := "0x" + hex.EncodeToString(arr)

			// Prepare response
			getSuiAddressResponse := GetSuiAddressResponse{
				Address:   suiAddress,
				PublicKey: publicKeyHex,
			}

			// Set content type header
			w.Header().Set("Content-Type", "application/json")

			// Encode and send response
			if err := json.NewEncoder(w).Encode(getSuiAddressResponse); err != nil {
				log.Printf("Error encoding Sui address response: %v", err)
				http.Error(w, fmt.Sprintf("Error building response: %v", err), http.StatusInternalServerError)
				return
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
		} else if keyCurve == RIPPLE_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

			getAddressResponse := GetAddressResponse{
				Address: ripple.PublicKeyToAddress(rawKeyEddsa),
			}
			err = json.NewEncoder(w).Encode(getAddressResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else if keyCurve == CARDANO_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

			// Get the public key
			pk := edwards.PublicKey{
				Curve: rawKeyEddsa.EDDSAPub.Curve(),
				X:     rawKeyEddsa.EDDSAPub.X(),
				Y:     rawKeyEddsa.EDDSAPub.Y(),
			}

			publicKeyStr := hex.EncodeToString(pk.Serialize())

			getAddressResponse := GetAddressResponse{
				Address: publicKeyStr,
			}
			err = json.NewEncoder(w).Encode(getAddressResponse)
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
			http.Error(w, "Intent has expired", http.StatusBadRequest)
			return
		}

		msg := ""

		operation := intent.Operations[operationIndexInt]
		if operation.Type == sequencer.OPERATION_TYPE_TRANSACTION {
			msg = operation.DataToSign
		} else if operation.Type == sequencer.OPERATION_TYPE_SEND_TO_BRIDGE {
			// Verify only operation for bridging
			// Get bridgewallet by calling /getwallet from sequencer api
			req, err := http.NewRequest("GET", "/getWallet?identity="+intent.Identity+"&identityCurve="+intent.IdentityCurve, nil)
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

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Sugar().Errorw("error reading response body", "error", err)
				return
			}

			var bridgeWallet sequencer.WalletSchema
			err = json.Unmarshal(body, &bridgeWallet)
			if err != nil {
				logger.Sugar().Errorw("error unmarshalling response body", "error", err)
				return
			}

			if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
				chain, err := common.GetChain(operation.ChainId)
				if err != nil {
					logger.Sugar().Errorw("error getting chain", "error", err)
					return
				}

				// Extract destination address from serialized transaction
				var destAddress string
				if chain.ChainType == "bitcoin" || chain.ChainType == "dogecoin" {
					// For Bitcoin, decode the serialized transaction to get output address
					var tx wire.MsgTx
					txBytes, err := hex.DecodeString(operation.SerializedTxn)
					if err != nil {
						logger.Sugar().Errorw("error decoding bitcoin&dogecoin transaction", "error", err)
						return
					}
					if err := tx.Deserialize(bytes.NewReader(txBytes)); err != nil {
						logger.Sugar().Errorw("error deserializing bitcoin&dogecoin transaction", "error", err)
						return
					}
					// Get the first output's address (assuming it's the bridge address)
					if len(tx.TxOut) > 0 {
						_, addrs, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[0].PkScript, nil)
						if err != nil || len(addrs) == 0 {
							logger.Sugar().Errorw("error extracting bitcoin&dogecoin address", "error", err)
							return
						}
						destAddress = addrs[0].String()
					}
				} else {
					// For EVM chains, decode the transaction to get the 'to' address
					txBytes, err := hex.DecodeString(operation.SerializedTxn)
					if err != nil {
						logger.Sugar().Errorw("error decoding EVM transaction", "error", err)
						return
					}
					tx := new(types.Transaction)
					if err := rlp.DecodeBytes(txBytes, tx); err != nil {
						logger.Sugar().Errorw("error deserializing EVM transaction", "error", err)
						return
					}
					destAddress = tx.To().Hex()
				}

				// Verify destination address matches bridge wallet
				var expectedAddress string
				if chain.ChainType == "bitcoin" {
					expectedAddress = bridgeWallet.BitcoinMainnetPublicKey
				} else if chain.ChainType == "dogecoin" {
					expectedAddress = bridgeWallet.DogecoinMainnetPublicKey
				} else {
					expectedAddress = bridgeWallet.ECDSAPublicKey
				}

				if !strings.EqualFold(destAddress, expectedAddress) {
					logger.Sugar().Errorw("Invalid bridge destination address", "expected", expectedAddress, "got", destAddress)
					return
				}
			} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" || operation.KeyCurve == "cardano_eddsa" || operation.KeyCurve == "sui_eddsa" {
				chain, err := common.GetChain(operation.ChainId)
				if err != nil {
					logger.Sugar().Errorw("error getting chain", "error", err)
					return
				}

				// Verify destination address matches bridge wallet based on chain type
				var validDestination bool
				var destAddress string

				// Extract destination address from serialized transaction based on chain type
				switch chain.ChainType {
				case "solana":
					// Decode base58 transaction and extract destination
					decodedTxn, err := base58.Decode(operation.SerializedTxn)
					if err != nil {
						logger.Sugar().Errorw("error decoding Solana transaction", "error", err)
						return
					}
					tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTxn))
					if err != nil || len(tx.Message.Instructions) == 0 {
						logger.Sugar().Errorw("error deserializing Solana transaction", "error", err)
						return
					}
					// Get the first instruction's destination account index
					destAccountIndex := tx.Message.Instructions[0].Accounts[1]
					// Get the actual account address from the message accounts
					destAddress = tx.Message.AccountKeys[destAccountIndex].String()
				case "aptos":
					// For Aptos, the destination is in the transaction payload
					var aptosPayload struct {
						Function string   `json:"function"`
						Args     []string `json:"arguments"`
					}
					if err := json.Unmarshal([]byte(operation.SerializedTxn), &aptosPayload); err != nil {
						logger.Sugar().Errorw("error parsing Aptos transaction", "error", err)
						return
					}
					if len(aptosPayload.Args) > 0 {
						destAddress = aptosPayload.Args[0] // First arg is typically the recipient
					}
				case "stellar":
					// For Stellar, parse the XDR transaction envelope
					var txEnv xdr.TransactionEnvelope
					err := xdr.SafeUnmarshalBase64(operation.SerializedTxn, &txEnv)
					if err != nil {
						logger.Sugar().Errorw("error parsing Stellar transaction", "error", err)
						return
					}

					// Get the first operation's destination
					if len(txEnv.Operations()) > 0 {
						if paymentOp, ok := txEnv.Operations()[0].Body.GetPaymentOp(); ok {
							destAddress = paymentOp.Destination.Address()
						}
					}
				case "algorand":
					txnBytes, err := base64.StdEncoding.DecodeString(operation.SerializedTxn)
					if err != nil {
						logger.Sugar().Errorw("failed to decode serialized transaction", "error", err)
						return
					}
					var txn algorandTypes.Transaction
					err = msgpack.Decode(txnBytes, &txn)
					if err != nil {
						logger.Sugar().Errorw("failed to deserialize transaction", "error", err)
						return
					}
					if txn.Type == algorandTypes.PaymentTx {
						destAddress = txn.PaymentTxnFields.Receiver.String()
					} else if txn.Type == algorandTypes.AssetTransferTx {
						destAddress = txn.AssetTransferTxnFields.AssetReceiver.String()
					} else {
						logger.Sugar().Errorw("Unknown transaction type", "type", txn.Type)
						return
					}
				case "ripple":
					// For Ripple, the destination is in the transaction payload
					// Decode the serialized transaction
					txBytes, err := hex.DecodeString(strings.TrimPrefix(operation.SerializedTxn, "0x"))
					if err != nil {
						logger.Sugar().Errorw("error decoding transaction", "error", err)
						return
					}

					// Parse the transaction
					var tx data.Payment
					err = json.Unmarshal(txBytes, &tx)
					if err != nil {
						logger.Sugar().Errorw("error unmarshalling transaction", "error", err)
						return
					}
					destAddress = tx.Destination.String()
				case "cardano":
					var tx cardanolib.Tx
					txBytes, err := hex.DecodeString(operation.SerializedTxn)
					if err != nil {
						logger.Sugar().Errorw("error decoding Cardano transaction", "error", err)
						return
					}
					if err := json.Unmarshal(txBytes, &tx); err != nil {
						logger.Sugar().Errorw("error parsing Cardano transaction", "error", err)
						return
					}
					destAddress = tx.Body.Outputs[0].Address.String()
				case "sui":
					var tx sui_types.TransactionData
					txBytes, err := base64.StdEncoding.DecodeString(operation.SerializedTxn)
					if err != nil {
						logger.Sugar().Errorw("error decoding Sui transaction", "error", err)
						return
					}
					if err := json.Unmarshal(txBytes, &tx); err != nil {
						logger.Sugar().Errorw("error parsing Sui transaction", "error", err)
						return
					}
					if len(tx.V1.Kind.ProgrammableTransaction.Inputs) < 1 {
						logger.Sugar().Errorw("wrong format sui transaction", "error", err)
						return
					}
					destAddress = string(*tx.V1.Kind.ProgrammableTransaction.Inputs[0].Pure)
				}

				// Verify the extracted destination matches the bridge wallet
				if destAddress == "" {
					logger.Sugar().Errorw("Failed to extract destination address from %s transaction", chain.ChainType)
					validDestination = false
				} else {
					switch chain.ChainType {
					case "solana":
						validDestination = strings.EqualFold(destAddress, bridgeWallet.EDDSAPublicKey)
					case "aptos":
						validDestination = strings.EqualFold(destAddress, bridgeWallet.AptosEDDSAPublicKey)
					case "stellar":
						validDestination = strings.EqualFold(destAddress, bridgeWallet.StellarPublicKey)
					case "algorand":
						validDestination = strings.EqualFold(destAddress, bridgeWallet.AlgorandEDDSAPublicKey)
					case "ripple":
						validDestination = strings.EqualFold(destAddress, bridgeWallet.RippleEDDSAPublicKey)
					case "cardano":
						validDestination = strings.EqualFold(destAddress, bridgeWallet.CardanoPublicKey)
					case "sui":
						validDestination = strings.EqualFold(destAddress, bridgeWallet.SuiPublicKey)
					}
				}

				if !validDestination {
					logger.Sugar().Errorw("Invalid bridge destination address for", "chain", chain.ChainType)
					return
				}
			}

			// Set message
			msg = operation.DataToSign
		} else if operation.Type == sequencer.OPERATION_TYPE_SOLVER {
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
		} else if operation.Type == sequencer.OPERATION_TYPE_BRIDGE_DEPOSIT {
			// Validate previous operation
			depositOperation := intent.Operations[operationIndexInt-1]

			if operationIndexInt == 0 || !(depositOperation.Type == sequencer.OPERATION_TYPE_SEND_TO_BRIDGE) {
				logger.Sugar().Errorw("Invalid operation type for bridge deposit")
				return
			}

			chain, err := common.GetChain(operation.ChainId)
			if err != nil {
				logger.Sugar().Errorw("error getting chain", "error", err)
				return
			}

			var transfers []common.Transfer

			if chain.ChainType == "ethereum" {
				transfers, err = sequencer.GetEthereumTransfers(depositOperation.ChainId, depositOperation.Result, intent.Identity)
				if err != nil {
					logger.Sugar().Errorw("error getting ethereum transfers", "error", err)
					return
				}
			}

			if chain.ChainType == "solana" {
				// TODO: Helius API Key
				transfers, err = sequencer.GetSolanaTransfers(depositOperation.ChainId, depositOperation.Result, HeliusApiKey)
				if err != nil {
					logger.Sugar().Errorw("error getting solana transfers", "error", err)
					return
				}
			}

			if chain.ChainType == "dogecoin" {
				transfers, err = dogecoin.GetDogeTransfers(depositOperation.ChainId, depositOperation.Result)
				if err != nil {
					logger.Sugar().Errorw("error getting dogecoin transfers", "error", err)
					return
				}
			}

			if chain.ChainType == "aptos" {
				transfers, err = aptos.GetAptosTransfers(depositOperation.ChainId, depositOperation.Result)
				if err != nil {
					logger.Sugar().Errorw("error getting aptos transfers", "error", err)
					return
				}
			}

			if chain.ChainType == "bitcoin" {
				transfers, _, err = bitcoin.GetBitcoinTransfers(depositOperation.ChainId, depositOperation.Result)
				if err != nil {
					logger.Sugar().Errorw("error getting bitcoin transfers", "error", err)
					return
				}
			} else if chain.ChainType == "sui" {
				transfers, err = sui.GetSuiTransfers(depositOperation.ChainId, depositOperation.Result)
				if err != nil {
					logger.Sugar().Errorw("error getting sui transfers", "error", err)
					return
				}
			}

			if chain.ChainType == "algorand" {
				transfers, err = algorand.GetAlgorandTransfers(depositOperation.GenesisHash, depositOperation.Result)
				if err != nil {
					logger.Sugar().Errorw("error getting algorand transfers", "error", err)
					return
				}
			}
			if chain.ChainType == "stellar" {
				transfers, err = stellar.GetStellarTransfers(depositOperation.ChainId, depositOperation.Result)
				if err != nil {
					logger.Sugar().Errorw("error getting stellar transfers", "error", err)
					return
				}
			}

			if chain.ChainType == "ripple" {
				transfers, err = ripple.GetRippleTransfers(depositOperation.ChainId, depositOperation.Result)
				if err != nil {
					logger.Sugar().Errorw("error getting ripple transfers", "error", err)
					return
				}
			}

			if chain.ChainType == "cardano" {
				transfers, err = cardano.GetCardanoTransfers(depositOperation.ChainId, depositOperation.Result)
				if err != nil {
					logger.Sugar().Errorw("error getting cardano transfers", "error", err)
					return
				}
			}

			if len(transfers) == 0 {
				logger.Sugar().Errorw("No transfers found", "result", depositOperation.Result, "identity", intent.Identity)
				return
			}

			// check if the token exists
			transfer := transfers[0]
			srcAddress := transfer.TokenAddress

			exists, _, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, depositOperation.ChainId, srcAddress)
			if err != nil {
				logger.Sugar().Errorw("error checking token existence", "error", err)
				return
			}

			if !exists {
				logger.Sugar().Errorw("Token does not exist", "srcAddress", srcAddress, "chainId", depositOperation.ChainId)
				return
			}

			// Set message
			msg = operation.SolverDataToSign
		} else if operation.Type == sequencer.OPERATION_TYPE_SWAP {
			// Validate previous operation
			bridgeDeposit := intent.Operations[operationIndexInt-1]

			if operationIndexInt == 0 || !(bridgeDeposit.Type == sequencer.OPERATION_TYPE_BRIDGE_DEPOSIT) {
				logger.Sugar().Errorw("Invalid operation type for swap")
				return
			}

			// Set message
			msg = operation.SolverDataToSign
		} else if operation.Type == sequencer.OPERATION_TYPE_BURN {
			// Validate nearby operations
			bridgeSwap := intent.Operations[operationIndex-1]

			if operationIndex+1 >= len(intent.Operations) || intent.Operations[operationIndex+1].Type != sequencer.OPERATION_TYPE_WITHDRAW {
				fmt.Println("BURN operation must be followed by a WITHDRAW operation")
				return
			}

			if operationIndex == 0 || !(bridgeSwap.Type == sequencer.OPERATION_TYPE_SWAP) {
				logger.Sugar().Errorw("Invalid operation type for swap")
				return
			}

			logger.Sugar().Infow("Burning tokens", "bridgeSwap", bridgeSwap)

			// Set message
			msg = operation.SolverDataToSign
		} else if operation.Type == sequencer.OPERATION_TYPE_BURN_SYNTHETIC {
			// This operation allows direct burning of ERC20 tokens from the wallet
			// without requiring a prior swap operation
			burnSyntheticMetadata := BurnSyntheticMetadata{}

			// Verify that this operation is followed by a withdraw operation
			if operationIndex+1 >= len(intent.Operations) || intent.Operations[operationIndex+1].Type != sequencer.OPERATION_TYPE_WITHDRAW {
				logger.Sugar().Errorw("BURN_SYNTHETIC operation must be followed by a WITHDRAW operation")
				return
			}

			err := json.Unmarshal([]byte(operation.SolverMetadata), &burnSyntheticMetadata)
			if err != nil {
				logger.Sugar().Errorw("Error unmarshalling burn synthetic metadata:", "error", err)
				return
			}

			// Get bridgewallet by calling /getwallet from sequencer api
			req, err := http.NewRequest("GET", "/getWallet?identity="+intent.Identity+"&identityCurve="+intent.IdentityCurve, nil)
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

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Sugar().Errorw("error reading response body", "error", err)
				return
			}

			var bridgeWallet sequencer.WalletSchema
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
		} else if operation.Type == sequencer.OPERATION_TYPE_WITHDRAW {
			// Verify nearby operations
			burn := intent.Operations[operationIndex-1]

			if operationIndex == 0 || !(burn.Type == sequencer.OPERATION_TYPE_BURN || burn.Type == sequencer.OPERATION_TYPE_BURN_SYNTHETIC) {
				logger.Sugar().Errorw("Invalid operation type for withdraw after burn")
				return
			}

			var withdrawMetadata WithdrawMetadata
			json.Unmarshal([]byte(operation.SolverMetadata), &withdrawMetadata)

			// Handle different burn operation types
			var tokenToWithdraw string
			var burnTokenAddress string
			if burn.Type == sequencer.OPERATION_TYPE_BURN {
				var burnMetadata BurnMetadata
				json.Unmarshal([]byte(burn.SolverMetadata), &burnMetadata)
				tokenToWithdraw = withdrawMetadata.Token
				burnTokenAddress = burnMetadata.Token
			} else if burn.Type == sequencer.OPERATION_TYPE_BURN_SYNTHETIC {
				var burnSyntheticMetadata BurnSyntheticMetadata
				json.Unmarshal([]byte(burn.SolverMetadata), &burnSyntheticMetadata)
				tokenToWithdraw = withdrawMetadata.Token
				burnTokenAddress = burnSyntheticMetadata.Token
			}

			// verify these fields
			exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, operation.ChainId, tokenToWithdraw)

			if err != nil {
				logger.Sugar().Errorw("error checking token existence", "error", err)
				return
			}

			if !exists {
				logger.Sugar().Errorw("Token does not exist", "token", tokenToWithdraw, "chainId", operation.ChainId)
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
		identityCurve := intent.IdentityCurve
		keyCurve := operation.KeyCurve

		log.Println("keyCurve", keyCurve)
		log.Println("msg", msg)

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
			if operation.Type == sequencer.OPERATION_TYPE_BRIDGE_DEPOSIT ||
				operation.Type == sequencer.OPERATION_TYPE_SWAP ||
				operation.Type == sequencer.OPERATION_TYPE_BURN ||
				operation.Type == sequencer.OPERATION_TYPE_WITHDRAW {
				go generateSignatureMessage(BridgeContractAddress, "ecdsa", "ecdsa", []byte(msg))
			} else {
				go generateSignatureMessage(identity, identityCurve, keyCurve, []byte(msg))
			}
		} else if keyCurve == BITCOIN_CURVE {
			go generateSignatureMessage(identity, identityCurve, keyCurve, []byte(msg))
		} else if keyCurve == DOGECOIN_CURVE {
			go generateSignatureMessage(identity, identityCurve, keyCurve, []byte(msg))
		} else if keyCurve == SUI_EDDSA_CURVE {
			// For Sui, we need to format the message according to Sui's standards
			// The message should be prefixed with "Sui Message:" for personal messages
			// suiMsg := []byte("Sui Message:" + msg)
			// go generateSignatureMessage(identity, identityCurve, keyCurve, suiMsg)
			msgBytes, _ := lib.NewBase64Data(msg)
			go generateSignatureMessage(identity, identityCurve, keyCurve, *msgBytes)
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
		} else if keyCurve == RIPPLE_CURVE || keyCurve == CARDANO_CURVE {
			msgBytes, err := hex.DecodeString(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("error decoding %s message: %v", keyCurve, err), http.StatusInternalServerError)
				return
			}
			go generateSignatureMessage(identity, identityCurve, keyCurve, msgBytes)
		} else {
			http.Error(w, "invalid key curve", http.StatusBadRequest)
			return
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

		w.Header().Set("Content-Type", "application/json")

		signatureResponse := SignatureResponse{}

		if keyCurve == ECDSA_CURVE {
			signatureResponse.Signature = string(sig.Message)
			signatureResponse.Address = sig.Address
		} else if keyCurve == BITCOIN_CURVE {
			signatureResponse.Signature = string(sig.Message)
			signatureResponse.Address = sig.Address
			logger.Sugar().Infof("signatureResponse: %v", signatureResponse)
		} else if keyCurve == DOGECOIN_CURVE {
			// signatureResponse.Signature = hex.EncodeToString(sig.Message)
			signatureResponse.Signature = string(sig.Message)

			signatureResponse.Address = sig.Address
		} else if keyCurve == SUI_EDDSA_CURVE {
			// For Sui, we return the signature in base64 format
			signatureResponse.Signature = string(sig.Message) // Already base64 encoded in generateSignature
			signatureResponse.Address = sig.Address
			logger.Sugar().Infof("generated Sui signature for address: %s", sig.Message)
		} else if keyCurve == APTOS_EDDSA_CURVE || keyCurve == STELLAR_CURVE || keyCurve == RIPPLE_CURVE || keyCurve == CARDANO_CURVE {
			signatureResponse.Signature = hex.EncodeToString(sig.Message)
			logger.Sugar().Infof("generated signature: %s", hex.EncodeToString(sig.Message))
			signatureResponse.Address = sig.Address
		} else if keyCurve == ALGORAND_CURVE {
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
				http.Error(w, fmt.Sprintf("Error marshaling algodMsg to JSON: %v", err), http.StatusInternalServerError)
				return
			}
			v, err := identityVerification.VerifySignature(sig.Address, "algorand_eddsa", string(jsonBytes), signatureResponse.Signature)
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

type SignatureResponse struct {
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

type GetDogecoinAddressesResponse struct {
	MainnetAddress string `json:"mainnetAddress"`
	TestnetAddress string `json:"testnetAddress"`
}

type GetSuiAddressResponse struct {
	Address   string `json:"address"`
	PublicKey string `json:"publicKey"` // Full Ed25519 public key in hex
}
