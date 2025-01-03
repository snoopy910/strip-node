package sequencer

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

// Bitcoin integration constants
const (
	BTC_TOKEN_SYMBOL = "BTC"
	SATOSHI_DECIMALS = 8
	BTC_ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"
)

// GetBitcoinTransfers fetches Bitcoin transaction details and parses them into transfers
func GetBitcoinTransfers(chainId string, txHash string) ([]Transfer, error) {
	// Get chain configuration
	chain, err := GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	// Create a new RPC client
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         chain.ChainUrl,
		User:         chain.RpcUsername, // TODO: need to replace username and password
		Pass:         chain.RpcPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client: %v", err)
	}
	defer client.Shutdown()

	// Convert transaction hash from string to *chainhash.Hash
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction hash: %v", err)
	}

	// Fetch transaction details
	rawTx, err := client.GetRawTransactionVerbose(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Bitcoin transaction: %v", err)
	}

	var transfers []Transfer

	// Process each input
	for _, input := range rawTx.Vin {
		// Fetch the previous transaction
		prevTxHash, err := chainhash.NewHashFromStr(input.Txid)
		if err != nil {
			return nil, fmt.Errorf("failed to parse previous transaction hash: %v", err)
		}

		prevTx, err := client.GetRawTransactionVerbose(prevTxHash)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch previous transaction: %v", err)
		}

		// Get the address from the previous transaction's output
		prevOut := prevTx.Vout[input.Vout]
		if len(prevOut.ScriptPubKey.Addresses) == 0 {
			continue
		}

		fromAddress := prevOut.ScriptPubKey.Addresses[0]

		// Process outputs
		for _, output := range rawTx.Vout {
			if len(output.ScriptPubKey.Addresses) == 0 {
				continue
			}

			amount := big.NewFloat(output.Value)
			amount.Mul(amount, big.NewFloat(1e8)) // Convert BTC to Satoshis
			floatValue, _ := amount.Float64()
			scaledAmount := fmt.Sprintf("%d", int64(floatValue))

			formattedAmount, err := getFormattedAmount(scaledAmount, SATOSHI_DECIMALS)
			if err != nil {
				return nil, fmt.Errorf("error formatting amount: %w", err)
			}

			transfers = append(transfers, Transfer{
				From:         fromAddress,
				To:           output.ScriptPubKey.Addresses[0],
				Amount:       formattedAmount,
				Token:        BTC_TOKEN_SYMBOL,
				IsNative:     true,
				TokenAddress: BTC_ZERO_ADDRESS,
				ScaledAmount: scaledAmount,
			})
		}
	}

	return transfers, nil
}

// Helper function to format amounts
func getFormattedAmount(amount string, decimal int) (string, error) {
	bigIntAmount := new(big.Int)

	_, success := bigIntAmount.SetString(amount, 10)
	if !success {
		return "", fmt.Errorf("error: Invalid number string")
	}

	formattedAmount, err := FormatUnits(bigIntAmount, decimal)
	if err != nil {
		return "", err
	}

	return formattedAmount, nil
}
