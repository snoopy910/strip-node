package blockchains

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
)

// NewSolanaBlockchain creates a new Solana blockchain instance
func NewSolanaBlockchain(networkType NetworkType) (IBlockchain, error) {
	apiKey := os.Getenv("HELIUS_API_KEY")
	heliusURL := "https://api.helius.xyz/v0/transactions?api-key=" + apiKey
	chainId := "900"
	network := Network{
		networkType: networkType,
		nodeURL:     "https://api.solana.com",
		networkID:   "mainnet",
	}

	if networkType == Devnet {
		apiKey = "6ccb4a2e-a0e6-4af3-afd0-1e06e1439547"
		network.nodeURL = "https://api.devnet.solana.com"
		network.networkID = "devnet"
		heliusURL = "https://api-devnet.helius.xyz/v0/transactions?api-key=" + apiKey
		chainId = "901"
	}

	client := rpc.New(network.nodeURL)

	return &SolanaBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Solana,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "base64",
			decimals:        7,
			chainID:         &chainId,
			opTimeout:       time.Second * 30,
			tokenSymbol:     "SOL",
		},
		client:    client,
		heliusURL: heliusURL,
	}, nil
}

// This is a type assertion to ensure that the SolanaBlockchain implements the IBlockchain interface
var _ IBlockchain = &SolanaBlockchain{}

// SolanaBlockchain implements the IBlockchain interface for Solana
type SolanaBlockchain struct {
	BaseBlockchain
	client    *rpc.Client
	heliusURL string
}

func validateAndOrderSignatures(tx *solana.Transaction) error {
	// -> check signature count in tests
	if len(tx.Signatures) != int(tx.Message.Header.NumRequiredSignatures) {
		return fmt.Errorf("signature count mismatch: got %d, want %d",
			len(tx.Signatures), tx.Message.Header.NumRequiredSignatures)
	}

	return nil
}

func (b *SolanaBlockchain) BroadcastTransaction(serializedTxn string, signatureBase58 string, _ *string) (string, error) {
	fmt.Printf("Solana Transaction Params:\n"+
		"  serializedTxn: %s\n"+
		"  chainId: %s\n"+
		"  keyCurve: %s\n"+
		"  signatureBase58: %s\n",
		serializedTxn, b.network.networkID, b.keyCurve, signatureBase58)

	// Decode the message data
	decodedTransactionData, err := base58.Decode(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("failed to decode message data: %v", err)
	}

	// Deserialize the binary data into a Solana transaction
	// This reconstructs the transaction object with all its instructions
	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		return "", fmt.Errorf("failed to deserialize transaction data: %v", err)
	}

	// Debug logging for message details
	fmt.Printf("Message Details:\n"+
		"  Header: %+v\n"+
		"  AccountKeys: %v\n"+
		"  RecentBlockhash: %s\n"+
		"  Instructions Count: %d\n",
		_tx.Message.Header, _tx.Message.AccountKeys, _tx.Message.RecentBlockhash, len(_tx.Message.Instructions))

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
	hash, err := b.client.SendTransaction(context.Background(), _tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// Return the transaction hash as a string
	return hash.String(), nil
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

func (b *SolanaBlockchain) GetTransfers(txnHash string, address *string) ([]common.Transfer, error) {
	fmt.Printf("Getting Solana transfers for transaction - chainId: %s, txnHash: %s\n", b.network.networkID, txnHash)
	// Parse Helius API response
	var heliusResponse []HeliusResponse

	heliusResponse, err := requestHeliusTransactionDetails(b.heliusURL, txnHash)
	if err != nil {
		fmt.Printf("Error getting transaction details from Helius: %v\n", err)
		return nil, err
	}

	if len(heliusResponse) == 0 {
		ticker := time.NewTicker(1 * time.Second)
		timeout := time.After(10 * time.Second)
		defer ticker.Stop()
	request:
		for {
			select {
			case <-ticker.C:
				heliusResponse, err = requestHeliusTransactionDetails(b.heliusURL, txnHash)
				if err != nil {
					fmt.Printf("Error getting transaction details from Helius: %v\n", err)
					return nil, err
				}
				if len(heliusResponse) > 0 {
					fmt.Printf("Got transaction details from Helius\n")
					break request
				}
			case <-timeout:
				fmt.Printf("Timeout waiting for transaction details from Helius\n")
				break request
			}
		}
	}

	fmt.Printf("Parsed Helius response - got %d transaction(s)\n", len(heliusResponse))
	if len(heliusResponse) == 0 {
		fmt.Printf("WARNING: Helius returned empty response for transaction %s\n", txnHash)
		return []common.Transfer{}, nil
	}

	var transfers []common.Transfer

	// Process each transaction in the response
	for i, response := range heliusResponse {
		fmt.Printf("Processing transaction %d - Native transfers: %d, Token transfers: %d\n",
			i+1, len(response.NativeTransfers), len(response.TokenTransfers))

		// Handle native SOL transfers
		for j, nativeTransfer := range response.NativeTransfers {
			// Convert amount to big.Int and format with 9 decimals (SOL decimal places)
			num, _ := new(big.Int).SetString(fmt.Sprintf("%d", nativeTransfer.Amount), 10)
			formattedAmount, _ := util.FormatUnits(num, 9)

			fmt.Printf("Native transfer %d: %s SOL from %s to %s\n",
				j+1, formattedAmount, nativeTransfer.FromUserAccount, nativeTransfer.ToUserAccount)

			// Create transfer record for native SOL
			transfers = append(transfers, common.Transfer{
				From:         nativeTransfer.FromUserAccount,
				To:           nativeTransfer.ToUserAccount,
				Amount:       formattedAmount,
				Token:        b.TokenSymbol(),
				IsNative:     true,
				TokenAddress: util.ZERO_ADDRESS,
				ScaledAmount: num.String(),
			})
		}

		// Handle SPL token transfers
		for j, tokenTransfer := range response.TokenTransfers {
			fmt.Printf("Token transfer %d: Mint: %s, Standard: %s\n",
				j+1, tokenTransfer.Mint, tokenTransfer.TokenStandard)

			// Skip non-fungible token transfers (e.g., NFTs)
			if tokenTransfer.TokenStandard != "Fungible" {
				fmt.Printf("Skipping non-fungible token transfer (standard: %s)\n", tokenTransfer.TokenStandard)
				continue
			}

			// Initialize RPC client for token metadata
			c := rpc.New(b.network.nodeURL)

			// Get token mint account address
			accountAddress := solana.MustPublicKeyFromBase58(tokenTransfer.Mint)
			fmt.Printf("Getting token metadata for mint: %s\n", accountAddress)

			// Fetch token mint account data for decimals
			accountInfo, err := c.GetAccountInfo(context.Background(), accountAddress)

			if err != nil {
				fmt.Printf("Failed to get account info for token %s: %v\n", tokenTransfer.Mint, err)
				return nil, fmt.Errorf("failed to get account info: %v", err)
			}

			// Decode mint account data to get token decimals
			var mint token.Mint
			err = bin.NewBinDecoder(accountInfo.GetBinary()).Decode(&mint)
			if err != nil {
				fmt.Printf("Failed to decode mint data: %v\n", err)
				panic(err)
			}
			fmt.Printf("Token %s has %d decimals\n", tokenTransfer.Mint, mint.Decimals)

			// Format token amount using the correct number of decimals
			num, _ := new(big.Int).SetString(fmt.Sprintf("%d", tokenTransfer.TokenAmount), 10)
			formattedAmount, err := util.FormatUnits(num, int(mint.Decimals))

			if err != nil {
				fmt.Printf("Error formatting token amount: %v\n", err)
				return nil, err
			}

			fmt.Printf("Token transfer: %s tokens from %s to %s\n",
				formattedAmount, tokenTransfer.FromUserAccount, tokenTransfer.ToUserAccount)

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

	fmt.Printf("Extracted %d transfers from transaction %s\n", len(transfers), txnHash)
	if len(transfers) == 0 {
		fmt.Printf("WARNING: No transfers found in transaction %s - this might indicate an issue with the transaction type or API parsing\n", txnHash)
	}

	return transfers, nil
}

func requestHeliusTransactionDetails(heliusTransactionUrl string, txnHash string) ([]HeliusResponse, error) {

	// Prepare request body with transaction hash
	requestBody := HeliusRequest{
		Transactions: []string{txnHash},
	}

	// Marshal request to JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error marshaling request body: %v\n", err)
		return nil, err
	}
	fmt.Printf("Sending request to Helius API with transaction: %s\n", txnHash)

	// Create HTTP request to Helius API
	req, err := http.NewRequest("POST", heliusTransactionUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		return nil, err
	}

	// Set content type for JSON request
	req.Header.Set("Content-Type", "application/json")

	// Send request to Helius API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request to Helius API: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("Helius API response status: %s\n", resp.Status)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Log response body for debugging
	bodyStr := string(body)
	if len(bodyStr) > 500 {
		fmt.Printf("Helius API response (truncated): %s...\n", bodyStr[:500])
	} else {
		fmt.Printf("Helius API response: %s\n", bodyStr)
	}

	// Parse Helius API response
	var heliusResponse []HeliusResponse
	err = json.Unmarshal(body, &heliusResponse)
	if err != nil {
		fmt.Printf("Failed to parse JSON response: %v\n", err)
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return heliusResponse, nil
}

func (b *SolanaBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	signature, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		return false, err
	}

	// Regarding the deprecation of GetConfirmedTransaction in Solana-Core v2, this has been updated to use GetTransaction.
	// https://spl_governance.crates.io/docs/rpc/deprecated/getconfirmedtransaction
	txResp, err := b.client.GetTransaction(context.Background(), signature, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentConfirmed,
	})

	if err != nil {
		fmt.Printf("Solana RPC Error: %v (type: %T)\n", err, err)
		// Check for specific error types to provide better diagnostics
		if err.Error() == "not found" {
			fmt.Printf("Transaction %s was not found on the blockchain - it may have been rejected or never submitted\n", txHash)
		}
		return false, err
	}

	// Log transaction details
	fmt.Printf("Transaction found! BlockTime: %v, Slot: %d, Confirmations: %d\n",
		txResp.BlockTime,
		txResp.Slot,
		txResp.Meta.Err)

	return true, nil
}

func (b *SolanaBlockchain) BuildWithdrawTx(account string,
	solverOutput string,
	recipient string,
	tokenAddress *string,
) (string, string, error) {
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	amount, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	if *tokenAddress == util.ZERO_ADDRESS {
		accountFrom := solana.MustPublicKeyFromBase58(account)
		accountTo := solana.MustPublicKeyFromBase58(recipient)

		// convert amount to uint64
		_amount, _ := big.NewInt(0).SetString(amount, 10)
		amountUint64 := _amount.Uint64()

		recentHash, err := b.client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
		if err != nil {
			return "", "", err
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					amountUint64,
					accountFrom,
					accountTo,
				).Build(),
			},
			recentHash.Value.Blockhash,
			solana.TransactionPayer(accountFrom),
		)

		if err != nil {
			return "", "", err
		}

		_msg, err := tx.ToBase64()
		if err != nil {
			return "", "", err
		}

		_msgBytes, _ := base64.StdEncoding.DecodeString(_msg)
		_msgBase58 := base58.Encode(_msgBytes)

		msg, err := tx.Message.MarshalBinary()
		if err != nil {
			return "", "", err
		}

		return _msgBase58, base58.Encode(msg), nil
	}

	accountFrom := solana.MustPublicKeyFromBase58(account)
	accountTo := solana.MustPublicKeyFromBase58(recipient)
	tokenMint := solana.MustPublicKeyFromBase58(*tokenAddress)

	senderTokenAccount, _, err := solana.FindAssociatedTokenAddress(accountFrom, tokenMint)
	if err != nil {
		return "", "", fmt.Errorf("failed to get sender token account: %v", err)
	}

	recipientTokenAccount, _, err := solana.FindAssociatedTokenAddress(accountTo, tokenMint)
	if err != nil {
		return "", "", fmt.Errorf("failed to get recipient token account: %v", err)
	}

	// convert amount to uint64
	_amount, _ := big.NewInt(0).SetString(amount, 10)
	amountUint64 := _amount.Uint64()

	recentHash, err := b.client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return "", "", err
	}

	transferInstruction := token.NewTransferInstruction(
		amountUint64,
		senderTokenAccount,
		recipientTokenAccount,
		accountFrom,
		nil, // No multisig signers
	).Build()

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			transferInstruction,
		},
		recentHash.Value.Blockhash,
		solana.TransactionPayer(accountFrom),
	)

	if err != nil {
		return "", "", err
	}

	_msg, err := tx.ToBase64()
	if err != nil {
		return "", "", err
	}

	_msgBytes, _ := base64.StdEncoding.DecodeString(_msg)
	_msgBase58 := base58.Encode(_msgBytes)

	msg, err := tx.Message.MarshalBinary()
	if err != nil {
		return "", "", err
	}

	return _msgBase58, base58.Encode(msg), nil
}

func (b *SolanaBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *SolanaBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}

func (b *SolanaBlockchain) ExtractDestinationAddress(serializedTxn string) (string, string, error) {
	// Decode base58 transaction and extract destination
	decodedTxn, err := base58.Decode(serializedTxn)
	destAddress := ""
	if err != nil {
		return "", "", fmt.Errorf("error decoding Solana transaction", err)
	}
	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTxn))
	if err != nil || len(tx.Message.Instructions) == 0 {
		return "", "", fmt.Errorf("error deserializing Solana transaction", err)
	}
	// Get the first instruction's destination account index
	destAccountIndex := tx.Message.Instructions[0].Accounts[1]
	// Get the actual account address from the message accounts
	destAddress = tx.Message.AccountKeys[destAccountIndex].String()

	return destAddress, "", nil
}
