package solana

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/StripChain/strip-node/common"

	"github.com/StripChain/strip-node/util"
	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
)

func TestBuildSolana() {
	accountFrom := solana.MustPublicKeyFromBase58("DpZqkyDKkVv2S7Lhbd5dUVcVCPJz2Lypr4W5Cru2sHr7")
	accountTo := solana.MustPublicKeyFromBase58("5oNDL3swdJJF1g9DzJiZ4ynHXgszjAEpUkxVYejchzrY")
	amount := uint64(1)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				accountFrom,
				accountTo,
			).Build(),
		},
		solana.MustHashFromBase58("CBLp4VEPu9T9W2uzURoawLGqgAQ65LvmUwDYRHymgwbd"),
		solana.TransactionPayer(accountFrom),
	)

	if err != nil {
		panic(err)
	}

	sig, _ := base58.Decode("5jLFtNTCAnHA9uurWhyNNqzwHLwWCaSNrZBWG48AANMGkreX1DYGbkHL2VWNNt2Kz327QwzzsAacJj2YFdSsfkwN")
	signature := solana.SignatureFromBytes(sig)

	_msg, err := tx.ToBase64()
	if err != nil {
		panic(err)
	}

	_msgBytes, _ := base64.StdEncoding.DecodeString(_msg)
	_msgBase58 := base58.Encode(_msgBytes)

	decodedTransactionData, err := base58.Decode(_msgBase58)
	if err != nil {
		fmt.Println("Error decoding transaction data:", err)
		return
	}

	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		panic(err)
	}

	_tx.Signatures = []solana.Signature{signature} // Reset signatures array with fee payer signature

	err = _tx.VerifySignatures()
	if err != nil {
		fmt.Println("error during verification")
		fmt.Println(err)
	} else {
		fmt.Println("Signatures verified")
	}
}

type NativeTransfer struct {
	FromUserAccount string `json:"fromUserAccount"`
	ToUserAccount   string `json:"toUserAccount"`
	Amount          uint   `json:"amount"`
}

type TokenTransfer struct {
	FromUserAccount  string `json:"fromUserAccount"`
	ToUserAccount    string `json:"toUserAccount"`
	FromTokenAccount string `json:"fromTokenAccount"`
	ToTokenAccount   string `json:"toTokenAccount"`
	TokenAmount      uint   `json:"tokenAmount"`
	Mint             string `json:"mint"`
	TokenStandard    string `json:"tokenStandard"`
}

type HeliusResponse struct {
	NativeTransfers []NativeTransfer `json:"nativeTransfers"`
	TokenTransfers  []TokenTransfer  `json:"tokenTransfers"`
}

type HeliusRequest struct {
	Transactions []string `json:"transactions"`
}

func validateAndOrderSignatures(tx *solana.Transaction) error {
	// -> check signature count in tests
	if len(tx.Signatures) != int(tx.Message.Header.NumRequiredSignatures) {
		return fmt.Errorf("signature count mismatch: got %d, want %d",
			len(tx.Signatures), tx.Message.Header.NumRequiredSignatures)
	}

	return nil
}

// SendSolanaTransactionWithValidation submits a signed Solana transaction with thorough validation and detailed logging.
// It prints all input parameters, decodes and unmarshals the transaction, verifies the signature, and provides verbose error context upon failure.
// Intended for debugging and development; for a leaner production path use SendSolanaTransaction.
func SendSolanaTransactionWithValidation(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureBase58 string) (string, error) {
	fmt.Printf("Solana Transaction Params:\n"+
		"  serializedTxn: %s\n"+
		"  chainId: %s\n"+
		"  keyCurve: %s\n"+
		"  dataToSign: %s\n"+
		"  signatureBase58: %s\n",
		serializedTxn, chainId, keyCurve, dataToSign, signatureBase58)

	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", err
	}

	c := rpc.New(chain.ChainUrl)

	// Decode the message data
	messageData, err := base58.Decode(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("failed to decode message data: %v", err)
	}

	// Create a message decoder and decode the message
	decoder := bin.NewBinDecoder(messageData)
	message := new(solana.Message)
	err = message.UnmarshalWithDecoder(decoder)
	if err != nil {
		return "", fmt.Errorf("failed to decode message: %v", err)
	}

	// Debug logging for message details
	fmt.Printf("Message Details:\n"+
		"  Header: %+v\n"+
		"  AccountKeys: %v\n"+
		"  RecentBlockhash: %s\n"+
		"  Instructions Count: %d\n",
		message.Header, message.AccountKeys, message.RecentBlockhash, len(message.Instructions))

	// Create a new transaction with the decoded message
	tx := &solana.Transaction{
		Message: *message,
	}

	// Decode and add the signature
	sig, err := base58.Decode(signatureBase58)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// The first account (fee payer) must sign
	signature := solana.SignatureFromBytes(sig)
	tx.Signatures = []solana.Signature{signature}

	// Verify the transaction is well-formed
	if err := tx.VerifySignatures(); err != nil {
		// Get the message bytes that were signed
		msgBytes, mErr := message.MarshalBinary()
		if mErr != nil {
			return "", fmt.Errorf("failed to marshal message: %v (original error: %v)", mErr, err)
		}

		// Add more detailed error information
		feePayer := "no fee payer"
		if len(message.AccountKeys) > 0 {
			feePayer = message.AccountKeys[0].String()
		}

		fmt.Printf("Signature Verification Details:\n"+
			"  Fee Payer: %s\n"+
			"  Signature: %s\n"+
			"  Message (base58): %s\n",
			feePayer, signature, base58.Encode(msgBytes))
		return "", fmt.Errorf("signature verification failed: %v", err)
	}

	// Send the transaction
	hash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("error during sending transaction:", err)
		return "", err
	}

	return hash.String(), nil
}

func GetSolanaTransfers(chainId string, txnHash string, apiKey string) ([]common.Transfer, error) {
	// Configure Helius API URL based on chain ID
	// Currently only supports devnet (chainId 901)
	var url string
	if chainId == "901" {
		url = "https://api-devnet.helius.xyz/v0/transactions?api-key=" + apiKey
	} else {
		return nil, fmt.Errorf("unsupported chainId: %s", chainId)
	}

	// Get chain configuration for native token info and RPC URL
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, err
	}

	// Prepare request body with transaction hash
	requestBody := HeliusRequest{
		Transactions: []string{txnHash},
	}

	// Marshal request to JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Create HTTP request to Helius API
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	// Set content type for JSON request
	req.Header.Set("Content-Type", "application/json")

	// Send request to Helius API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse Helius API response
	var heliusResponse []HeliusResponse
	err = json.Unmarshal(body, &heliusResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	var transfers []common.Transfer

	// Process each transaction in the response
	for _, response := range heliusResponse {
		// Handle native SOL transfers
		for _, nativeTransfer := range response.NativeTransfers {
			// Convert amount to big.Int and format with 9 decimals (SOL decimal places)
			num, _ := new(big.Int).SetString(fmt.Sprintf("%d", nativeTransfer.Amount), 10)
			formattedAmount, _ := util.FormatUnits(num, 9)

			// Create transfer record for native SOL
			transfers = append(transfers, common.Transfer{
				From:         nativeTransfer.FromUserAccount,
				To:           nativeTransfer.ToUserAccount,
				Amount:       formattedAmount,
				Token:        chain.TokenSymbol,
				IsNative:     true,
				TokenAddress: util.ZERO_ADDRESS,
				ScaledAmount: num.String(),
			})
		}

		// Handle SPL token transfers
		for _, tokenTransfer := range response.TokenTransfers {
			// Skip non-fungible token transfers (e.g., NFTs)
			if tokenTransfer.TokenStandard != "Fungible" {
				continue
			}

			// Initialize RPC client for token metadata
			c := rpc.New(chain.ChainUrl)

			// Get token mint account address
			accountAddress := solana.MustPublicKeyFromBase58(tokenTransfer.Mint)
			// Fetch token mint account data for decimals
			accountInfo, err := c.GetAccountInfo(context.Background(), accountAddress)

			if err != nil {
				return nil, fmt.Errorf("failed to get account info: %v", err)
			}

			spew.Dump(accountInfo)

			// Decode mint account data to get token decimals
			var mint token.Mint
			err = bin.NewBinDecoder(accountInfo.GetBinary()).Decode(&mint)
			if err != nil {
				panic(err)
			}
			spew.Dump(mint)

			// Format token amount using the correct number of decimals
			num, _ := new(big.Int).SetString(fmt.Sprintf("%d", tokenTransfer.TokenAmount), 10)
			formattedAmount, err := util.FormatUnits(num, int(mint.Decimals))

			if err != nil {
				return nil, err
			}

			// Create transfer record for SPL token
			transfers = append(transfers, common.Transfer{
				From:         tokenTransfer.FromUserAccount,
				To:           tokenTransfer.ToUserAccount,
				Amount:       formattedAmount,
				Token:        tokenTransfer.Mint,
				IsNative:     false,
				TokenAddress: tokenTransfer.Mint,
				ScaledAmount: num.String(),
			})
		}
	}

	fmt.Println(transfers)
	return transfers, nil
}

func CheckSolanaTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	c := rpc.New(chain.ChainUrl)

	signature, err := solana.SignatureFromBase58(txnHash)
	if err != nil {
		return false, err
	}

	// Regarding the deprecation of GetConfirmedTransaction in Solana-Core v2, this has been updated to use GetTransaction.
	// https://spl_governance.crates.io/docs/rpc/deprecated/getconfirmedtransaction
	_, err = c.GetTransaction(context.Background(), signature, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentConfirmed,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// SendSolanaTransaction submits a signed Solana transaction to the network
func SendSolanaTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureBase58 string) (string, error) {
	// Get chain configuration for RPC endpoint
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", err
	}

	// Initialize Solana RPC client
	c := rpc.New(chain.ChainUrl)

	// Decode the base58-encoded transaction data
	// Solana transactions are serialized using a custom binary format and base58-encoded
	decodedTransactionData, err := base58.Decode(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction data: %v", err)
	}

	// Deserialize the binary data into a Solana transaction
	// This reconstructs the transaction object with all its instructions
	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		return "", fmt.Errorf("failed to deserialize transaction data: %v", err)
	}

	// Decode the base58-encoded signature and convert it to Solana's signature format
	// Solana uses 64-byte Ed25519 signatures
	sig, _ := base58.Decode(signatureBase58)
	signature := solana.SignatureFromBytes(sig)

	// Add the signature to the transaction
	// Solana transactions can have multiple signatures for multi-sig transactions
	_tx.Signatures = append(_tx.Signatures, signature)

	// Verify that all required signatures are present and valid
	// This checks signatures against the transaction data and account permissions
	err = _tx.VerifySignatures()
	if err != nil {
		return "", fmt.Errorf("failed to verify signatures: %v", err)
	}

	// Submit the transaction to the Solana network
	// The returned hash can be used to track the transaction status
	hash, err := c.SendTransaction(context.Background(), _tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// Return the transaction hash as a string
	return hash.String(), nil
}
