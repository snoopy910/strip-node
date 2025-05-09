package blockchains

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

var EVMRegistry = map[BlockchainID]bool{
	Ethereum:   true,
	Arbitrum:   true,
	StripChain: true,
	Berachain:  true,
	Sonic:      true,
}

// RegisterEVMBlockchain adds a new EVM blockchain ID to the registry.
// It returns an error if a blockchain with the same ID is already registered.
// func RegisterEVMBlockchain(chainName BlockchainID) error {
// 	evmRegistryMu.Lock()
// 	defer evmRegistryMu.Unlock()

// 	fmt.Printf("evmRegistry[chainName]: %+v\n", evmRegistry[chainName])
// 	if _, exists := evmRegistry[chainName]; exists {
// 		return fmt.Errorf("blockchain with ID %s already registered", chainName)
// 	}

// 	fmt.Printf("added evmRegistry[chainName]: %+v\n", evmRegistry[chainName])
// 	evmRegistry[chainName] = true
// 	return nil
// }

// IsEVMBlockchain checks if a blockchain is an EVM blockchain by its ID.
// It returns true if the blockchain is an EVM blockchain, otherwise false.
func IsEVMBlockchain(id BlockchainID) bool {
	// evmRegistryMu.RLock()
	// defer evmRegistryMu.RUnlock()

	// fmt.Printf("evmRegistry: %+v\n", evmRegistry)
	// fmt.Printf("evmRegistry[id]: %+v\n", evmRegistry[id])
	// fmt.Printf("id: %+v\n", id)

	// return evmRegistry[id]
	return EVMRegistry[id]
}

// // GetAllEVMBlockchainIDs returns a slice of all registered EVM blockchain IDs.
// func GetAllEVMBlockchainIDs() []BlockchainID {
// 	evmRegistryMu.RLock()
// 	defer evmRegistryMu.RUnlock()

// 	chains := make([]BlockchainID, 0, len(evmRegistry))
// 	for id := range evmRegistry {
// 		chains = append(chains, id)
// 	}
// 	return chains
// }

// NewEVMBlockchain initializes a new EVMBlockchain
func NewEVMBlockchain(
	chainName BlockchainID,
	network Network,
	signingEncoding string,
	decimals uint,
	opTimeout time.Duration,
	chainID *string,
	tokenSymbol string,
) (EVMBlockchain, error) {
	client, err := ethclient.Dial(network.nodeURL)
	if err != nil {
		logger.Sugar().Errorw("failed to create ethclient", "error", err)
		return EVMBlockchain{}, err
	}

	// RegisterEVMBlockchain(chainName)

	return EVMBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       chainName,
			network:         network,
			keyCurve:        common.CurveEcdsa,
			signingEncoding: signingEncoding,
			decimals:        decimals,
			opTimeout:       opTimeout,
			chainID:         chainID,
			tokenSymbol:     tokenSymbol,
		},
		client: client,
	}, nil
}

// This is a type assertion to ensure that the EVMBlockchain implements the IBlockchain interface
var _ IBlockchain = &EVMBlockchain{}

// EVMBlockchain implements the IBlockchain interface for Ethereum
type EVMBlockchain struct {
	BaseBlockchain
	client *ethclient.Client
}

func (b *EVMBlockchain) BroadcastTransaction(txn string, signatureHex string, publicKey *string) (string, error) {
	// Handle hex strings with or without 0x prefix
	var serializedTx []byte
	var err error

	// Ensure the transaction has 0x prefix for hexutil.Decode
	if !strings.HasPrefix(txn, "0x") && !strings.HasPrefix(txn, "0X") {
		// No prefix, use standard hex decoding
		serializedTx, err = hex.DecodeString(txn)
	} else {
		// Has prefix, use hexutil.Decode
		serializedTx, err = hexutil.Decode(txn)
	}

	if err != nil {
		return "", fmt.Errorf("failed to decode transaction hex: %w", err)
	}

	var tx types.Transaction
	rlp.DecodeBytes(serializedTx, &tx)

	// Handle signatures with or without 0x prefix
	var sigData []byte

	// Ensure the signature has 0x prefix for hexutil.Decode
	if !strings.HasPrefix(signatureHex, "0x") && !strings.HasPrefix(signatureHex, "0X") {
		// No prefix, use standard hex decoding
		sigData, err = hex.DecodeString(signatureHex)
	} else {
		// Has prefix, use hexutil.Decode
		sigData, err = hexutil.Decode(signatureHex)
	}

	if err != nil {
		return "", fmt.Errorf("failed to decode signature hex: %w", err)
	}

	n, _ := new(big.Int).SetString(*b.chainID, 10)
	_tx, err := tx.WithSignature(types.NewLondonSigner(n), []byte(sigData))

	if err != nil {
		return "", err
	}

	err = b.client.SendTransaction(context.Background(), _tx)
	if err != nil {
		return "", err
	}

	return _tx.Hash().Hex(), nil
}

func (b *EVMBlockchain) GetTransfers(txnHash string, ecdsaAddr *string) ([]common.Transfer, error) {
	logger.Sugar().Infow("Processing Ethereum transaction", "txHash", txnHash, "chainID", *b.chainID)

	// Get the full transaction to examine content and data
	tx, isPending, err := b.client.TransactionByHash(context.Background(), ethCommon.HexToHash(txnHash))
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
	receipt, err := b.client.TransactionReceipt(context.Background(), ethCommon.HexToHash(txnHash))
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
			From:     ethCommon.HexToAddress(*ecdsaAddr),
			To:       tx.To(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		}

		// Call the transaction in a simulated environment to get error message
		_, err := b.client.CallContract(context.Background(), msg, receipt.BlockNumber)
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
					balance, balErr := b.client.BalanceAt(context.Background(), ethCommon.HexToAddress(*ecdsaAddr), nil)
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
			decimal, symbol, err := getERC20Details(b.client, ethCommon.BytesToAddress(log.Address.Bytes()))

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
			From:         *ecdsaAddr,
			To:           tx.To().String(),
			Amount:       WeiToEther(wei).String(),
			Token:        b.TokenSymbol(),
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

func (b *EVMBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	_, isPending, err := b.client.TransactionByHash(context.Background(), ethCommon.HexToHash(txHash))
	if err != nil {
		return false, err
	}

	return !isPending, nil
}

func (b *EVMBlockchain) BuildWithdrawTx(bridgeAddress string,
	solverOutput string,
	userAddress string,
	tokenAddress *string,
) (string, string, error) {
	// Parse solver output to get amount
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	amountStr, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	// Parse the amount as a big integer
	amount, success := new(big.Int).SetString(amountStr, 10)
	if !success {
		return "", "", fmt.Errorf("failed to parse amount as big integer: %s", amountStr)
	}

	// Validate Ethereum addresses
	if !ethCommon.IsHexAddress(bridgeAddress) {
		return "", "", fmt.Errorf("invalid bridge address format: %s", bridgeAddress)
	}

	if !ethCommon.IsHexAddress(userAddress) {
		return "", "", fmt.Errorf("invalid user address format: %s", userAddress)
	}

	// Setup the transaction based on whether it's an ERC20 or native token transfer
	var tx *types.Transaction
	var err error
	var dataToSign string

	if tokenAddress != nil && *tokenAddress != "" && *tokenAddress != util.ZERO_ADDRESS {
		// ERC20 token withdrawal
		if !ethCommon.IsHexAddress(*tokenAddress) {
			return "", "", fmt.Errorf("invalid token address format: %s", *tokenAddress)
		}

		parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
		if err != nil {
			return "", "", fmt.Errorf("failed to parse ERC20 ABI: %v", err)
		}

		// Pack the transfer function call
		tokenAddressObj := ethCommon.HexToAddress(*tokenAddress)
		userAddressObj := ethCommon.HexToAddress(userAddress)

		// Create the transfer function call data
		transferFnData, err := parsedABI.Pack("transfer", userAddressObj, amount)
		if err != nil {
			return "", "", fmt.Errorf("failed to pack transfer function data: %v", err)
		}

		// Create the transaction
		bridgeAddressObj := ethCommon.HexToAddress(bridgeAddress)
		nonce, err := b.client.PendingNonceAt(context.Background(), bridgeAddressObj)
		if err != nil {
			return "", "", fmt.Errorf("failed to get nonce: %v", err)
		}

		gasPrice, err := b.client.SuggestGasPrice(context.Background())
		if err != nil {
			return "", "", fmt.Errorf("failed to get gas price: %v", err)
		}

		gasLimit := uint64(100000) // Safe default for ERC20 transfers

		tx = types.NewTransaction(
			nonce,
			tokenAddressObj,
			big.NewInt(0), // No ETH value for ERC20 transfers
			gasLimit,
			gasPrice,
			transferFnData,
		)
	} else {
		// Native token (ETH/BNB/etc.) withdrawal
		bridgeAddressObj := ethCommon.HexToAddress(bridgeAddress)
		userAddressObj := ethCommon.HexToAddress(userAddress)

		nonce, err := b.client.PendingNonceAt(context.Background(), bridgeAddressObj)
		if err != nil {
			return "", "", fmt.Errorf("failed to get nonce: %v", err)
		}

		gasPrice, err := b.client.SuggestGasPrice(context.Background())
		if err != nil {
			return "", "", fmt.Errorf("failed to get gas price: %v", err)
		}

		gasLimit := uint64(21000) // Standard gas limit for ETH transfers

		tx = types.NewTransaction(
			nonce,
			userAddressObj,
			amount,
			gasLimit,
			gasPrice,
			nil, // No data for native transfers
		)
	}

	// Get the chain ID
	chainID, err := b.client.ChainID(context.Background())
	if err != nil {
		return "", "", fmt.Errorf("failed to get chain ID: %v", err)
	}

	// Prepare the transaction for signing
	signer := types.NewEIP155Signer(chainID)
	// Remove the "0x" prefix to make it compatible with the validator's signature generation
	dataToSign = fmt.Sprintf("%x", signer.Hash(tx).Bytes())

	// Serialize the transaction
	txBytes, err := tx.MarshalBinary()
	if err != nil {
		return "", "", fmt.Errorf("failed to serialize transaction: %v", err)
	}

	return hexutil.Encode(txBytes), dataToSign, nil
}

func (b *EVMBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *EVMBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}

// Common ERC20 error signatures
const (
	// Error function selector: Error(string)
	ErrorSig = "08c379a0"

	// Custom error selectors
	InsufficientBalanceSig     = "e450d38c"
	InsufficientAllowanceSig   = "4bd67a2d"
	TransferFromZeroAddressSig = "ea553b34"
	TransferToZeroAddressSig   = "d92e233d"
)

// Error ABI definitions for common ERC20 errors
const errorABI = `[
	{
		"inputs": [
			{
				"name": "reason",
				"type": "string"
			}
		],
		"name": "Error",
		"type": "error"
	},
	{
		"inputs": [
			{
				"name": "sender",
				"type": "address"
			},
			{
				"name": "balance",
				"type": "uint256"
			},
			{
				"name": "needed",
				"type": "uint256"
			}
		],
		"name": "InsufficientBalance",
		"type": "error"
	},
	{
		"inputs": [
			{
				"name": "owner",
				"type": "address"
			},
			{
				"name": "spender",
				"type": "address"
			},
			{
				"name": "allowance",
				"type": "uint256"
			},
			{
				"name": "needed",
				"type": "uint256"
			}
		],
		"name": "InsufficientAllowance",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "TransferFromZeroAddress",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "TransferToZeroAddress",
		"type": "error"
	}
]`

// DecodeERC20RevertReason attempts to decode Ethereum revert reason data into a human-readable format
// Uses ABI decoding for standard errors when possible, with fallback to manual decoding
func DecodeERC20RevertReason(errorMsg string) string {
	// Check if this is a revert error
	if !strings.Contains(errorMsg, "execution reverted") {
		return errorMsg
	}

	// Try to use ABI to decode the error first
	decodedABI := decodeErrorWithABI(errorMsg)
	if decodedABI != "" {
		return decodedABI
	}

	// Fall back to manual decoding if ABI approach failed
	return manualDecodeError(errorMsg)
}

// decodeErrorWithABI attempts to decode error using ABI definitions
func decodeErrorWithABI(errorMsg string) string {
	// Parse the ABI
	_, err := abi.JSON(strings.NewReader(errorABI))
	if err != nil {
		return "" // Fall back to manual decoding
	}

	// Extract the error data
	var selector string
	var arguments []byte

	// Check for custom error format with selector and data
	if strings.Contains(errorMsg, "custom error") {
		parts := strings.Split(errorMsg, "custom error ")
		if len(parts) > 1 {
			dataParts := strings.Split(parts[1], ": ")
			if len(dataParts) > 1 && len(dataParts[0]) >= 10 {
				selector = strings.TrimPrefix(dataParts[0], "0x")

				// Try to decode the hex data
				if hexData, err := hex.DecodeString(strings.Replace(dataParts[1], "0x", "", 1)); err == nil {
					arguments = hexData
				}
			}
		}

		// If we have a selector and arguments, try to match with ABI
		if selector != "" && len(arguments) > 0 {
			switch selector {
			case InsufficientBalanceSig:
				// Try to decode InsufficientBalance error
				var sender ethCommon.Address
				var balance, needed *big.Int

				// Manually decode arguments based on their types
				if len(arguments) >= 96 { // 3 parameters of 32 bytes each
					sender = ethCommon.BytesToAddress(arguments[12:32])
					balance = new(big.Int).SetBytes(arguments[32:64])
					needed = new(big.Int).SetBytes(arguments[64:96])

					return fmt.Sprintf("Insufficient Balance Error: address %s has balance %s but needs %s",
						sender.Hex(), balance.String(), needed.String())
				}

			case InsufficientAllowanceSig:
				// Try to decode InsufficientAllowance error
				if len(arguments) >= 128 { // 4 parameters of 32 bytes each
					owner := ethCommon.BytesToAddress(arguments[12:32])
					spender := ethCommon.BytesToAddress(arguments[44:64])
					allowance := new(big.Int).SetBytes(arguments[64:96])
					needed := new(big.Int).SetBytes(arguments[96:128])

					return fmt.Sprintf("Insufficient Allowance Error: spender %s is allowed %s by owner %s but needs %s",
						spender.Hex(), allowance.String(), owner.Hex(), needed.String())
				}

			case TransferFromZeroAddressSig:
				return "Error: Transfer from the zero address"

			case TransferToZeroAddressSig:
				return "Error: Transfer to the zero address"

			case ErrorSig:
				// Standard Error(string) revert
				if len(arguments) >= 96 {
					// String offset (first 32 bytes)
					// String length (next 32 bytes)
					stringLen := new(big.Int).SetBytes(arguments[32:64]).Int64()

					// String content starts at offset 64 and goes for stringLen bytes
					if len(arguments) >= int(64+stringLen) {
						return fmt.Sprintf("Error: %s", string(arguments[64:64+stringLen]))
					}
				}
			}
		}
	}

	// Standard "execution reverted: X" format
	if strings.Contains(errorMsg, "execution reverted: ") {
		parts := strings.Split(errorMsg, "execution reverted: ")
		if len(parts) > 1 && !strings.Contains(parts[1], "custom error") {
			return fmt.Sprintf("Error: %s", parts[1])
		}
	}

	return ""
}

// manualDecodeError is a fallback for when ABI decoding fails
// Uses string parsing and manual byte extraction to decode known error patterns
func manualDecodeError(errorMsg string) string {
	// Extract the custom error data
	customErrorData := ""
	if strings.Contains(errorMsg, "custom error") {
		parts := strings.Split(errorMsg, "custom error ")
		if len(parts) > 1 {
			dataParts := strings.Split(parts[1], ": ")
			if len(dataParts) > 1 {
				customErrorData = dataParts[0] + ": " + dataParts[1] // Combine error code and data
			} else {
				customErrorData = parts[1]
			}
		}
	}

	if customErrorData != "" {
		// Try to extract some meaning even without formal decoding
		if strings.HasPrefix(customErrorData, "0xe450d38c") {
			return "Insufficient Balance Error (manual decode)"
		} else if strings.HasPrefix(customErrorData, "0x4bd67a2d") {
			return "Insufficient Allowance Error (manual decode)"
		} else if strings.HasPrefix(customErrorData, "0xea553b34") {
			return "Transfer from the zero address"
		} else if strings.HasPrefix(customErrorData, "0xd92e233d") {
			return "Transfer to the zero address"
		}

		return fmt.Sprintf("Unknown error data: %s", customErrorData)
	}

	// Return the original message if all else fails
	return errorMsg
}

func (b *EVMBlockchain) ExtractDestinationAddress(serializedTxn string) (string, string, error) {
	// For EVM chains, decode the transaction to get the 'to' address
	destAddress := ""
	tokenAddress := ""
	txBytes, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return "", "", fmt.Errorf("error decoding EVM transaction", err)
	}
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(txBytes, tx); err != nil {
		return "", "", fmt.Errorf("error deserializing EVM transaction", err)
	}
	if tx.To() == nil {
		return "", "", fmt.Errorf("EVM transaction has nil To address")
	}
	if len(tx.Data()) >= 4 && bytes.Equal(tx.Data()[:4], []byte{0xa9, 0x05, 0x9c, 0xbb}) {
		// ERC20 transfer detected, extract recipient from call data
		if len(tx.Data()) < 36 {
			return "", "", fmt.Errorf("ERC20 transfer data too short to extract destination")
		}
		destAddress = ethCommon.BytesToAddress(tx.Data()[4:36]).Hex()

		// For ERC20 transfers, verify the token exists in the bridge contract
		tokenAddress = tx.To().Hex()
	} else {
		destAddress = tx.To().Hex()
	}
	return destAddress, tokenAddress, nil
}
