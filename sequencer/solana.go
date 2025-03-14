package sequencer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/mr-tron/base58"

	"github.com/gagliardetto/solana-go/rpc"
)

func TestBuildSolana() {
	// pubKeyByte, _ := base58.Decode("4dEgqPG9FtjCiwW1HUdReraBozwq6qUCcDFXD8BnUn9Z")

	// c := rpc.New("https://api.devnet.solana.com")
	// recentHash, err := c.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(recentHash.Value.Blockhash.String())

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

	// msg, err := tx.Message.MarshalBinary()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Input to sign: ", base58.Encode(msg))

	// fmt.Println("Input to sign: ", base58.Encode(msg))

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

	// // signature
	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		panic(err)
	}

	_tx.Signatures = append(_tx.Signatures, signature)

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

// GetSolanaTransfers retrieves and parses transfer information from a Solana transaction
// Uses Helius API for transaction data and native RPC for token metadata
// Handles both native SOL transfers and SPL token transfers
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
	body, err := ioutil.ReadAll(resp.Body)
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
			formattedAmount, _ := FormatUnits(num, 9)

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
			// Account{}.Data.GetBinary() returns the *decoded* binary data
			// regardless the original encoding (it can handle them all).
			err = bin.NewBinDecoder(accountInfo.GetBinary()).Decode(&mint)
			if err != nil {
				panic(err)
			}
			spew.Dump(mint)

			// Format token amount using the correct number of decimals
			num, _ := new(big.Int).SetString(fmt.Sprintf("%d", tokenTransfer.TokenAmount), 10)
			formattedAmount, err := FormatUnits(num, int(mint.Decimals))

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
