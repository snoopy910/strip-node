package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

// GetEthereumTransfers retrieves and parses transfer information from an Ethereum transaction
// Handles both native ETH transfers and ERC20 token transfers by analyzing transaction logs
// Uses the standard ERC20 Transfer event signature to detect token transfers
func GetEthereumTransfers(chainId string, txnHash string, ecdsaAddr string) ([]common.Transfer, error) {
	// Get chain configuration for RPC URL and native token symbol
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, err
	}

	// Initialize Ethereum client with chain RPC URL
	client, err := ethclient.Dial(chain.ChainUrl)
	if err != nil {
		logger.Sugar().Errorw("failed to dial ethclient", "error", err)
		return nil, err
	}

	logger.Sugar().Infow("Processing Ethereum transaction", "txHash", txnHash, "chainId", chainId)

	// Get the full transaction to examine content and data
	tx, isPending, err := client.TransactionByHash(context.Background(), ethCommon.HexToHash(txnHash))
	if err != nil {
		logger.Sugar().Errorw("failed to get transaction details", "txHash", txnHash, "error", err)
		return nil, err
	}

	// Log detailed transaction parameters to aid debugging
	logger.Sugar().Infow("Raw transaction details",
		"txHash", txnHash,
		"isPending", isPending,
		"to", tx.To().Hex(),
		"value", tx.Value().String(),
		"gas", tx.Gas(),
		"gasPrice", tx.GasPrice().String(),
		"dataLength", len(tx.Data()),
		"dataPrefixHex", func() string {
			if len(tx.Data()) == 0 {
				return ""
			}
			maxLen := 8
			if len(tx.Data()) < maxLen {
				maxLen = len(tx.Data())
			}
			return fmt.Sprintf("%x", tx.Data()[:maxLen])
		}())

	// Try to decode function selector if it exists
	if len(tx.Data()) >= 4 {
		selector := fmt.Sprintf("%x", tx.Data()[:4])
		logger.Sugar().Infow("Function selector", "selector", selector)

		// Common ERC20 function selectors for reference
		switch selector {
		case "a9059cbb":
			logger.Sugar().Infow("Detected ERC20 transfer() call")
			if len(tx.Data()) >= 68 {
				// Extract recipient address (second 32 bytes after function selector)
				recipient := ethCommon.BytesToAddress(tx.Data()[4+12 : 4+32]).Hex()
				// Extract amount (last 32 bytes)
				amount := new(big.Int).SetBytes(tx.Data()[4+32 : 4+64])
				logger.Sugar().Infow("ERC20 transfer details from call data",
					"recipient", recipient,
					"amount", amount.String())
			}
		case "23b872dd":
			logger.Sugar().Infow("Detected ERC20 transferFrom() call")
		}
	}

	// Get transaction receipt which contains logs of token transfers
	receipt, err := client.TransactionReceipt(context.Background(), ethCommon.HexToHash(txnHash))
	if err != nil {
		logger.Sugar().Errorw("failed to get transaction receipt", "txHash", txnHash, "error", err)
		return nil, err
	}

	logger.Sugar().Infow("Got transaction receipt", "txHash", txnHash, "logs", len(receipt.Logs), "status", receipt.Status)

	// Check for failed transaction and try to get revert reason
	if receipt.Status == 0 {
		logger.Sugar().Warnw("Transaction failed on-chain", "txHash", txnHash)

		// Try to get revert reason (this requires tracing, which may not be available on all nodes)
		msg := ethereum.CallMsg{
			From:     ethCommon.HexToAddress(ecdsaAddr),
			To:       tx.To(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		}

		// Call the transaction in a simulated environment to get error message
		_, err := client.CallContract(context.Background(), msg, receipt.BlockNumber)
		if err != nil {
			errorMsg := err.Error()
			logger.Sugar().Infow("Revert reason (raw)", "error", errorMsg)

			// Use our specialized decoder to get a more readable error message
			decodedMsg := DecodeERC20RevertReason(errorMsg)
			if decodedMsg != errorMsg {
				logger.Sugar().Infow("Revert reason (decoded)", "reason", decodedMsg)

				// Enhanced error logging based on the type of error detected
				if strings.Contains(decodedMsg, "Insufficient Balance") {
					// Try to extract account and amount information for better diagnostics
					logger.Sugar().Warnw("Transaction failed due to insufficient balance",
						"txHash", txnHash,
						"from", ecdsaAddr,
						"to", tx.To().Hex(),
						"value", tx.Value().String())

					// Check account balance for more context
					balance, balErr := client.BalanceAt(context.Background(), ethCommon.HexToAddress(ecdsaAddr), nil)
					if balErr == nil {
						logger.Sugar().Infow("Account balance",
							"address", ecdsaAddr,
							"balance", balance.String(),
							"balanceEth", WeiToEther(balance).String())
					}
				} else if strings.Contains(decodedMsg, "Insufficient Allowance") {
					// For token approvals, log additional context
					logger.Sugar().Warnw("Transaction failed due to insufficient token allowance",
						"txHash", txnHash)

					// If this is a token transfer, extract the token contract address
					if len(tx.Data()) >= 4 && (fmt.Sprintf("%x", tx.Data()[:4]) == "a9059cbb" || // transfer
						fmt.Sprintf("%x", tx.Data()[:4]) == "23b872dd") { // transferFrom

						logger.Sugar().Infow("Token approval required before transfer")
					}
				}
			}

			// Return a specialized error that includes the revert reason
			return nil, fmt.Errorf("transaction reverted: %s", decodedMsg)
		}

		// If we didn't get an error from CallContract but status is 0, report generic failure
		return nil, fmt.Errorf("transaction failed but revert reason could not be determined")
	}

	var transfers []common.Transfer

	// Process logs to find ERC20 token transfers
	for _, log := range receipt.Logs {
		// Check if log is an ERC20 Transfer event
		// Signature: Transfer(address,address,uint256)
		// Hash: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
		if log.Topics[0].Hex() == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
			logger.Sugar().Infow("Found ERC20 Transfer event", "txHash", txnHash, "logIndex", log.Index, "topics", len(log.Topics), "dataLength", len(log.Data))

			// Safety check for topics length
			if len(log.Topics) < 3 {
				logger.Sugar().Warnw("ERC20 Transfer event has insufficient topics", "txHash", txnHash, "topicsCount", len(log.Topics))
				continue
			}

			// Extract from and to addresses from log topics
			// Topics[1] = from address, Topics[2] = to address
			from := ethCommon.BytesToAddress(log.Topics[1].Bytes()).Hex()
			to := ethCommon.BytesToAddress(log.Topics[2].Bytes()).Hex()

			logger.Sugar().Infow("ERC20 transfer details", "from", from, "to", to, "contract", log.Address.Hex())

			// Get token details (decimals and symbol) from the token contract
			decimal, symbol, err := getERC20Details(client, ethCommon.BytesToAddress(log.Address.Bytes()))

			if err != nil {
				logger.Sugar().Warnw("Error getting ERC20 token details, using defaults", "contract", log.Address.Hex(), "error", err)
				// Use default values to avoid failing the entire transfer detection
				decimal = 18
				symbol = "UNKNOWN"
			}

			// Format token amount using correct number of decimals
			// log.Data contains the transfer amount in the token's smallest unit
			formattedAmount, err := util.FormatUnits(new(big.Int).SetBytes(log.Data), int(decimal))

			if err != nil {
				logger.Sugar().Warnw("Error formatting token amount", "error", err)
				// Use raw amount as string instead of failing
				formattedAmount = new(big.Int).SetBytes(log.Data).String()
			}

			logger.Sugar().Infow("ERC20 transfer amount", "raw", new(big.Int).SetBytes(log.Data).String(), "formatted", formattedAmount)

			// Create transfer record for ERC20 token
			transfers = append(transfers, common.Transfer{
				From:         from,
				To:           to,
				Amount:       formattedAmount,
				Token:        symbol,
				IsNative:     false,
				TokenAddress: log.Address.Hex(),
				ScaledAmount: new(big.Int).SetBytes(log.Data).String(),
			})
		}
	}

	// Get native ETH transfer amount in Wei
	wei := tx.Value()

	// If transaction value is non-zero, it's a native ETH transfer
	if wei.Cmp(big.NewInt(0)) != 0 {
		logger.Sugar().Infow("Detected native ETH transfer",
			"from", ecdsaAddr,
			"to", tx.To().Hex(),
			"value", wei.String())

		// Create transfer record for native ETH
		transfers = append(transfers, common.Transfer{
			From:         ecdsaAddr,
			To:           tx.To().String(),
			Amount:       WeiToEther(wei).String(),
			Token:        chain.TokenSymbol,
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: wei.String(),
		})
	}

	// Log summary of transfer detection
	logger.Sugar().Infow("Transfer detection results",
		"txHash", txnHash,
		"transfersFound", len(transfers),
		"transactionStatus", receipt.Status)

	return transfers, nil
}

func WeiToEther(wei *big.Int) *big.Float {
	weiFloat := new(big.Float).SetInt(wei)
	ether := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
	return ether
}

const (
	erc20ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`
)

func getERC20Details(client *ethclient.Client, tokenAddress ethCommon.Address) (uint8, string, error) {
	contractABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return 0, "", err
	}

	callData, err := contractABI.Pack("decimals")
	if err != nil {
		return 0, "", err
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, "", err
	}

	var decimals uint8
	err = contractABI.UnpackIntoInterface(&decimals, "decimals", result)
	if err != nil {
		return 0, "", err
	}

	callData, err = contractABI.Pack("symbol")
	if err != nil {
		return 0, "", err
	}

	msg = ethereum.CallMsg{
		To:   &tokenAddress,
		Data: callData,
	}

	result, err = client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, "", err
	}

	var symbol string
	err = contractABI.UnpackIntoInterface(&symbol, "symbol", result)
	if err != nil {
		return 0, "", err
	}

	return decimals, symbol, nil
}

func CheckEVMTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	client, err := ethclient.Dial(chain.ChainUrl)
	if err != nil {
		return false, fmt.Errorf("failed to dial EVM client: %v", err)
	}

	_, isPending, err := client.TransactionByHash(context.Background(), ethCommon.HexToHash(txnHash))
	if err != nil {
		return false, err
	}

	return !isPending, nil
}

func SendEVMTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureHex string) (string, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", err
	}

	client, err := ethclient.Dial(chain.ChainUrl)
	if err != nil {
		logger.Sugar().Errorw("failed to dial ethclient", "error", err)
		return "", err
	}

	serializedTx, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return "", err
	}

	var tx types.Transaction
	rlp.DecodeBytes(serializedTx, &tx)

	sigData, err := hex.DecodeString(signatureHex)

	if err != nil {
		return "", err
	}

	n, _ := new(big.Int).SetString(chainId, 10)
	_tx, err := tx.WithSignature(types.NewLondonSigner(n), []byte(sigData))

	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), _tx)
	if err != nil {
		return "", err
	}

	return _tx.Hash().Hex(), nil
}
