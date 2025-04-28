package lending_solver

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Contract ABI strings
const (
	LendingPoolABI = `[
		{"inputs": [{"name": "token", "type": "address"}, {"name": "amount", "type": "uint256"}], "name": "supply", "outputs": [], "stateMutability": "nonpayable", "type": "function"},
		{"inputs": [{"name": "amount", "type": "uint256"}], "name": "borrowStripUSD", "outputs": [], "stateMutability": "nonpayable", "type": "function"},
		{"inputs": [{"name": "amount", "type": "uint256"}], "name": "repayStripUSD", "outputs": [], "stateMutability": "nonpayable", "type": "function"},
		{"inputs": [{"name": "token", "type": "address"}, {"name": "amount", "type": "uint256"}], "name": "withdrawCollateral", "outputs": [], "stateMutability": "nonpayable", "type": "function"},
		{"type": "event", "name": "Supply", "inputs": [{"name": "token", "type": "address", "indexed": true, "internalType": "address"}, {"name": "user", "type": "address", "indexed": true, "internalType": "address"}, {"name": "amount", "type": "uint256", "indexed": false, "internalType": "uint256"}], "anonymous": false},
		{"type": "event", "name": "Borrow", "inputs": [{"name": "token", "type": "address", "indexed": true, "internalType": "address"}, {"name": "user", "type": "address", "indexed": true, "internalType": "address"}, {"name": "amount", "type": "uint256", "indexed": false, "internalType": "uint256"}, {"name": "rate", "type": "uint256", "indexed": false, "internalType": "uint256"}], "anonymous": false},
		{"type": "event", "name": "CollateralWithdrawn", "inputs": [{"name": "token", "type": "address", "indexed": true, "internalType": "address"}, {"name": "user", "type": "address", "indexed": true, "internalType": "address"}, {"name": "amount", "type": "uint256", "indexed": false, "internalType": "uint256"}], "anonymous": false},
		{"type": "event", "name": "StripUSDRepaid", "inputs": [{"name": "user", "type": "address", "indexed": true, "internalType": "address"}, {"name": "amountRepaid", "type": "uint256", "indexed": false, "internalType": "uint256"}, {"name": "remainingDebt", "type": "uint256", "indexed": false, "internalType": "uint256"}], "anonymous": false}
	]`
)

type TxParams struct {
	Calldata  []byte
	Nonce     uint64
	GasLimit  uint64
	GasTipCap *big.Int
	GasFeeCap *big.Int
	Deadline  time.Time
}

// TransactionStatus represents the status of a transaction
type TransactionStatus struct {
	TxHash  string
	Status  string // pending, success, failure
	Receipt *types.Receipt
	Error   error
}

// LendingSolver handles interactions with the LendingPool contract
type LendingSolver struct {
	client          *ethclient.Client
	chainId         *big.Int
	lendingPool     common.Address
	abi             abi.ABI
	statuses        map[string]*TransactionStatus
	txParams        map[string]TxParams
	cleanupInterval time.Duration
	cleanupStop     chan struct{}
	cleanupMutex    sync.Mutex
}

func Start(
	rpcURL string,
	httpPort string,
	lendingPoolAddress string,
	chainId int64,
) {
	keepAlive := make(chan string)
	// Initialize the solver
	solver, err := NewLendingSolver(rpcURL, chainId, lendingPoolAddress)
	if err != nil {
		panic(err)
	}

	// Start the solver's API server
	go startHTTPServer(solver, httpPort)

	// Keep the process running
	<-keepAlive
}

// NewLendingSolver creates a new instance of LendingSolver
func NewLendingSolver(rpcURL string, chainId int64, lendingPool string) (*LendingSolver, error) {
	fmt.Printf("Initializing Lending solver with pool %s\n", lendingPool)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(LendingPoolABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	solver := &LendingSolver{
		client:          client,
		chainId:         big.NewInt(chainId),
		lendingPool:     common.HexToAddress(lendingPool),
		abi:             contractAbi,
		statuses:        make(map[string]*TransactionStatus),
		txParams:        make(map[string]TxParams),
		cleanupInterval: 5 * time.Minute,
		cleanupStop:     make(chan struct{}),
	}

	// Start cleanup goroutine
	go solver.cleanupExpiredParams()

	return solver, nil
}

// Solve executes the lending operation with a signed transaction
func (s *LendingSolver) Solve(intent Intent, opIndex int, signature string) (string, error) {
	fmt.Printf("Solving operation %d for intent %s\n", opIndex, intent.Identity)
	if opIndex >= len(intent.Operations) {
		return "", fmt.Errorf("operation index out of range")
	}

	// Calculate operation hash
	operation := intent.Operations[opIndex]
	opBytes, err := json.Marshal(operation)
	if err != nil {
		return "", fmt.Errorf("failed to marshal operation: %v", err)
	}
	opHash := "0x" + hex.EncodeToString(crypto.Keccak256(opBytes))

	// Get stored parameters
	s.cleanupMutex.Lock()
	params, ok := s.txParams[opHash]
	s.cleanupMutex.Unlock()
	if !ok {
		return "", fmt.Errorf("transaction parameters not found or expired")
	}
	if time.Now().After(params.Deadline) {
		s.cleanupMutex.Lock()
		delete(s.txParams, opHash)
		s.cleanupMutex.Unlock()
		return "", fmt.Errorf("transaction parameters expired")
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   s.chainId,
		Nonce:     params.Nonce,
		To:        &s.lendingPool,
		Value:     big.NewInt(0),
		Gas:       params.GasLimit,
		GasTipCap: params.GasTipCap,
		GasFeeCap: params.GasFeeCap,
		Data:      params.Calldata,
	})

	// Decode signature
	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return "", fmt.Errorf("invalid signature: %v", err)
	}
	if len(sig) != 65 {
		return "", fmt.Errorf("invalid signature length: got %d, want 65", len(sig))
	}

	// Create signed transaction
	signer := types.NewLondonSigner(s.chainId)
	signedTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return "", fmt.Errorf("failed to add signature: %v", err)
	}

	// Send the transaction
	err = s.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// Store the transaction status
	s.cleanupMutex.Lock()
	s.statuses[signedTx.Hash().Hex()] = &TransactionStatus{
		TxHash: signedTx.Hash().Hex(),
		Status: "pending",
	}
	s.cleanupMutex.Unlock()

	return signedTx.Hash().Hex(), nil
}

// Status checks the status of a lending operation
func (s *LendingSolver) Status(intent Intent, opIndex int) (string, error) {
	if opIndex >= len(intent.Operations) {
		return "", fmt.Errorf("operation index out of bounds")
	}

	op := intent.Operations[opIndex]
	if op.Result == "" {
		return "pending", nil
	}

	// Get transaction status from memory
	s.cleanupMutex.Lock()
	status, exists := s.statuses[op.Result]
	s.cleanupMutex.Unlock()
	if !exists {
		return "pending", nil
	}
	if status.Error != nil {
		return "failure", status.Error
	}

	// If we haven't checked the receipt yet
	if status.Receipt == nil {
		receipt, err := s.client.TransactionReceipt(context.Background(), common.HexToHash(status.TxHash))
		if err != nil {
			if err == ethereum.NotFound {
				return "pending", nil
			}
			status.Error = err
			return "failure", err
		}

		status.Receipt = receipt
		if receipt.Status == 1 {
			status.Status = "success"
		} else {
			status.Status = "failure"
			status.Error = fmt.Errorf("transaction reverted")
		}
	}

	return status.Status, status.Error
}

// GetOutput retrieves the result of a lending operation
func (s *LendingSolver) cleanupExpiredParams() {
	fmt.Println("Starting cleanup routine for expired parameters")
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			s.cleanupMutex.Lock()
			for hash, params := range s.txParams {
				if now.After(params.Deadline) {
					delete(s.txParams, hash)
				}
			}
			s.cleanupMutex.Unlock()
		case <-s.cleanupStop:
			return
		}
	}
}

func (s *LendingSolver) StopCleanup() {
	fmt.Println("Stopping cleanup routine")
	s.cleanupMutex.Lock()
	close(s.cleanupStop)
	s.cleanupMutex.Unlock()
}

func (s *LendingSolver) GetOutput(intent Intent, opIndex int) (string, error) {
	if opIndex >= len(intent.Operations) {
		return "", fmt.Errorf("operation index out of bounds")
	}

	op := intent.Operations[opIndex]
	s.cleanupMutex.Lock()
	status, exists := s.statuses[op.Result]
	s.cleanupMutex.Unlock()
	if !exists {
		return "", fmt.Errorf("transaction not found")
	}
	if status.Error != nil {
		return "", status.Error
	}

	// Get operation metadata
	var metadata LendingMetadata
	if err := json.Unmarshal(op.SolverMetadata, &metadata); err != nil {
		return "", fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	// Create output based on operation type
	output := LendingOutput{
		TxHash: status.TxHash,
	}

	// Get transaction receipt
	receipt := status.Receipt
	if receipt == nil {
		return "", fmt.Errorf("transaction receipt not available")
	}

	// Add operation-specific data from event logs
	switch metadata.Action {
	case "supply":
		found := false
		for _, log := range receipt.Logs {
			if log.Address == s.lendingPool {
				event, err := s.abi.EventByID(log.Topics[0])
				if err != nil {
					continue
				}
				if event.Name == "Supply" {
					if len(log.Topics) < 3 {
						return "", fmt.Errorf("invalid Supply event: insufficient topics")
					}
					// Extract token address from first topic
					output.Token = common.BytesToAddress(log.Topics[1][:]).Hex()
					dataMap, err := s.abi.Unpack(event.Name, log.Data)
					if err != nil {
						return "", fmt.Errorf("failed to unpack Supply event data: %v", err)
					}
					amount, ok := dataMap[0].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse amount from event data")
					}
					output.Amount.Int = amount.String()
					found = true
					break
				}
			}
		}
		if !found {
			return "", fmt.Errorf("Supply event not found in receipt")
		}

	case "borrow":
		found := false
		for _, log := range receipt.Logs {
			if log.Address == s.lendingPool {
				event, err := s.abi.EventByID(log.Topics[0])
				if err != nil {
					continue
				}
				if event.Name == "Borrow" {
					if len(log.Topics) < 3 {
						return "", fmt.Errorf("invalid Borrow event: insufficient topics")
					}
					// Extract token address from first topic
					output.Token = common.BytesToAddress(log.Topics[1][:]).Hex()
					dataMap, err := s.abi.Unpack(event.Name, log.Data)
					if err != nil {
						return "", fmt.Errorf("failed to unpack Borrow event data: %v", err)
					}
					amount, ok := dataMap[0].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse amount from event data")
					}
					rate, ok := dataMap[1].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse rate from event data")
					}
					output.Amount.Int = amount.String()
					output.Rate = rate.String()
					found = true
					break
				}
			}
		}
		if !found {
			return "", fmt.Errorf("Borrow event not found in receipt")
		}

	case "repay":
		found := false
		for _, log := range receipt.Logs {
			if log.Address == s.lendingPool {
				event, err := s.abi.EventByID(log.Topics[0])
				if err != nil {
					continue
				}
				if event.Name == "StripUSDRepaid" {
					if len(log.Topics) < 2 {
						return "", fmt.Errorf("invalid StripUSDRepaid event: insufficient topics")
					}
					dataMap, err := s.abi.Unpack(event.Name, log.Data)
					if err != nil {
						return "", fmt.Errorf("failed to unpack StripUSDRepaid event data: %v", err)
					}
					amountRepaid, ok := dataMap[0].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse amountRepaid from event data")
					}
					remaining, ok := dataMap[1].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse remainingDebt from event data")
					}
					output.Amount.Int = amountRepaid.String()
					output.RemainingDebt = remaining.String()
					found = true
					break
				}
			}
		}
		if !found {
			return "", fmt.Errorf("StripUSDRepaid event not found in receipt")
		}

	case "withdraw":
		found := false
		for _, log := range receipt.Logs {
			if log.Address == s.lendingPool {
				event, err := s.abi.EventByID(log.Topics[0])
				if err != nil {
					continue
				}
				if event.Name == "CollateralWithdrawn" {
					if len(log.Topics) < 3 {
						return "", fmt.Errorf("invalid CollateralWithdrawn event: insufficient topics")
					}
					// Extract token address from first topic
					output.Token = common.BytesToAddress(log.Topics[1][:]).Hex()
					dataMap, err := s.abi.Unpack(event.Name, log.Data)
					if err != nil {
						return "", fmt.Errorf("failed to unpack CollateralWithdrawn event data: %v", err)
					}
					amount, ok := dataMap[0].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse amount from event data")
					}
					output.Amount.Int = amount.String()
					found = true
					break
				}
			}
		}
		if !found {
			return "", fmt.Errorf("CollateralWithdrawn event not found in receipt")
		}
	}
	outputBytes, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return string(outputBytes), nil
}

// Construct builds the transaction data for a lending operation
func (s *LendingSolver) Construct(intent Intent, opIndex int) (string, error) {
	if opIndex >= len(intent.Operations) {
		return "", fmt.Errorf("operation index out of bounds")
	}

	op := intent.Operations[opIndex]
	var metadata LendingMetadata
	if err := json.Unmarshal(op.SolverMetadata, &metadata); err != nil {
		return "", fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	// Get the function call data based on operation type
	var data string
	var err error
	switch metadata.Action {
	case "supply":
		data, err = s.constructSupply(metadata)
	case "borrow":
		data, err = s.constructBorrow(metadata)
	case "repay":
		data, err = s.constructRepay(metadata)
	case "withdraw":
		data, err = s.constructWithdraw(metadata)
	default:
		return "", fmt.Errorf("unknown action: %s", metadata.Action)
	}
	if err != nil {
		return "", err
	}

	// Get current nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), common.HexToAddress(intent.Identity))
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get gas price
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// Create unsigned transaction
	decodedData, err := hex.DecodeString(strings.TrimPrefix(data, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode data: %v", err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   s.chainId,
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)), // Set fee cap to 2x the tip
		Gas:       300000,                                    // Fixed gas limit, could be estimated
		To:        &s.lendingPool,
		Value:     big.NewInt(0),
		Data:      decodedData,
	})

	// Get the hash to be signed
	signer := types.NewLondonSigner(s.chainId)
	hash := signer.Hash(tx)

	// Calculate operation hash
	opBytes, err := json.Marshal(op)
	if err != nil {
		return "", fmt.Errorf("failed to marshal operation: %v", err)
	}
	opHash := "0x" + hex.EncodeToString(crypto.Keccak256(opBytes))

	// Store parameters with 5-minute deadline
	s.txParams[opHash] = TxParams{
		Calldata:  decodedData,
		Nonce:     nonce,
		GasLimit:  300000,
		GasTipCap: gasPrice,
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		Deadline:  time.Now().Add(5 * time.Minute),
	}

	return "0x" + hex.EncodeToString(hash.Bytes()), nil
}

// constructSupply builds transaction data for supplying assets
func (s *LendingSolver) constructSupply(metadata LendingMetadata) (string, error) {
	amount := new(big.Int)
	amount.SetString(metadata.Amount.Int, 10)

	data, err := s.abi.Pack("supply",
		common.HexToAddress(metadata.Token),
		amount,
	)
	if err != nil {
		return "", fmt.Errorf("failed to pack supply data: %v", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}

// constructBorrow builds transaction data for borrowing assets
func (s *LendingSolver) constructBorrow(metadata LendingMetadata) (string, error) {
	amount := new(big.Int)
	amount.SetString(metadata.Amount.Int, 10)

	data, err := s.abi.Pack("borrowStripUSD", amount)
	if err != nil {
		return "", fmt.Errorf("failed to pack borrow data: %v", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}

// constructRepay builds transaction data for repaying debt
func (s *LendingSolver) constructRepay(metadata LendingMetadata) (string, error) {
	amount := new(big.Int)
	amount.SetString(metadata.Amount.Int, 10)

	data, err := s.abi.Pack("repayStripUSD", amount)
	if err != nil {
		return "", fmt.Errorf("failed to pack repay data: %v", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}

// constructWithdraw builds transaction data for withdrawing assets
func (s *LendingSolver) constructWithdraw(metadata LendingMetadata) (string, error) {
	amount := new(big.Int)
	amount.SetString(metadata.Amount.Int, 10)

	data, err := s.abi.Pack("withdrawCollateral",
		common.HexToAddress(metadata.Token),
		amount,
	)
	if err != nil {
		return "", fmt.Errorf("failed to pack withdraw data: %v", err)
	}

	return "0x" + hex.EncodeToString(data), nil
}
