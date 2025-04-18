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
	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"github.com/stellar/go/clients/horizonclient"
)

// NewSolanaBlockchain creates a new Stellar blockchain instance
func NewSolanaBlockchain(networkType NetworkType) (IBlockchain, error) {
	apiKey := os.Getenv("HELIUS_API_KEY")
	heliusURL := "https://api.helius.xyz/v0/transactions?api-key=" + apiKey
	network := Network{
		networkType: networkType,
		nodeURL:     horizonclient.DefaultPublicNetClient.HorizonURL,
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = horizonclient.DefaultTestNetClient.HorizonURL
		network.networkID = "testnet"
		heliusURL = "https://api-devnet.helius.xyz/v0/transactions?api-key=" + apiKey
	}

	client := rpc.New(network.nodeURL)

	return &SolanaBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Solana,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "base64",
			decimals:        7,
			opTimeout:       time.Second * 10,
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

func (b *SolanaBlockchain) BroadcastTransaction(txn string, signatureBase58 string, _ *string) (string, error) {
	fmt.Printf("Solana Transaction Params:\n"+
		"  serializedTxn: %s\n"+
		"  chainId: %s\n"+
		"  keyCurve: %s\n"+
		"  signatureBase58: %s\n",
		txn, b.network.networkID, b.keyCurve, signatureBase58)

	// Decode the message data
	messageData, err := base58.Decode(txn)
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
	hash, err := b.client.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("error during sending transaction:", err)
		return "", err
	}

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

func (b *SolanaBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	// Configure Helius API URL based on chain ID
	// Currently only supports devnet (chainId 901)

	// Prepare request body with transaction hash
	requestBody := HeliusRequest{
		Transactions: []string{txHash},
	}

	// Marshal request to JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Create HTTP request to Helius API
	req, err := http.NewRequest("POST", b.heliusURL, bytes.NewBuffer(requestBodyBytes))
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
				Token:        b.TokenSymbol(),
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

			// Get token mint account address
			accountAddress := solana.MustPublicKeyFromBase58(tokenTransfer.Mint)
			// Fetch token mint account data for decimals
			accountInfo, err := b.client.GetAccountInfo(context.Background(), accountAddress)

			if err != nil {
				return nil, fmt.Errorf("failed to get account info: %v", err)
			}

			spew.Dump(accountInfo)

			// Decode mint account data to get token decimals
			var mint token.Mint
			// Account{}.Data.GetBinary() returns the *decoded* binary data
			// regardless the original encoding (it can handle them all).
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

	return transfers, nil
}

func (b *SolanaBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	signature, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		return false, err
	}

	// Regarding the deprecation of GetConfirmedTransaction in Solana-Core v2, this has been updated to use GetTransaction.
	// https://spl_governance.crates.io/docs/rpc/deprecated/getconfirmedtransaction
	_, err = b.client.GetTransaction(context.Background(), signature, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentConfirmed,
	})

	if err != nil {
		return false, err
	}

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

	if tokenAddress == nil {
		accountFrom := solana.MustPublicKeyFromBase58(account)
		accountTo := solana.MustPublicKeyFromBase58(recipient)

		// convert amount to uint64
		_amount, _ := big.NewInt(0).SetString(amount, 10)
		amountUint64 := _amount.Uint64()

		recentHash, err := b.client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
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
