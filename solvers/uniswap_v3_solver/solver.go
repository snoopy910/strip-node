package uniswap_v3_solver

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
const npmABI = `[
	{"type":"function","name":"mint","inputs":[{"name":"params","type":"tuple","components":[{"name":"token0","type":"address"},{"name":"token1","type":"address"},{"name":"fee","type":"uint24"},{"name":"tickLower","type":"int24"},{"name":"tickUpper","type":"int24"},{"name":"amount0Desired","type":"uint256"},{"name":"amount1Desired","type":"uint256"},{"name":"amount0Min","type":"uint256"},{"name":"amount1Min","type":"uint256"},{"name":"recipient","type":"address"},{"name":"deadline","type":"uint256"}]}],"outputs":[{"name":"tokenId","type":"uint256"},{"name":"liquidity","type":"uint128"},{"name":"amount0","type":"uint256"},{"name":"amount1","type":"uint256"}],"stateMutability":"payable"},
	{"type":"function","name":"decreaseLiquidity","inputs":[{"name":"params","type":"tuple","components":[{"name":"tokenId","type":"uint256"},{"name":"liquidity","type":"uint128"},{"name":"amount0Min","type":"uint256"},{"name":"amount1Min","type":"uint256"},{"name":"deadline","type":"uint256"}]}],"outputs":[{"name":"amount0","type":"uint256"},{"name":"amount1","type":"uint256"}],"stateMutability":"nonpayable"},
	{"type":"event","name":"DecreaseLiquidity","inputs":[{"name":"tokenId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"liquidity","type":"uint128","indexed":false,"internalType":"uint128"},{"name":"amount0","type":"uint256","indexed":false,"internalType":"uint256"},{"name":"amount1","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},
	{"type":"event","name":"IncreaseLiquidity","inputs":[{"name":"tokenId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"liquidity","type":"uint128","indexed":false,"internalType":"uint128"},{"name":"amount0","type":"uint256","indexed":false,"internalType":"uint256"},{"name":"amount1","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},
	{"type":"function","name":"collect","inputs":[{"name":"tokenId","type":"uint256"},{"name":"recipient","type":"address"},{"name":"amount0Max","type":"uint256"},{"name":"amount1Max","type":"uint256"}],"outputs":[{"name":"amount0","type":"uint256"},{"name":"amount1","type":"uint256"}],"stateMutability":"nonpayable"}
]`

// TransactionStatus represents the status of a transaction
type TransactionStatus struct {
	TxHash  string
	Status  string // pending, success, failure
	Receipt *types.Receipt
	Error   error
}

type TxParams struct {
	Calldata  []byte
	Nonce     uint64
	GasLimit  uint64
	GasTipCap *big.Int
	GasFeeCap *big.Int
	Deadline  time.Time
}

// UniswapV3Solver handles interactions with the Uniswap V3 NonfungiblePositionManager contract
type UniswapV3Solver struct {
	client          *ethclient.Client
	chainId         *big.Int
	factory         common.Address
	npm             common.Address
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
	uniswapV3FactoryAddress string,
	npmAddress string,
	chainId int64,
) {
	keepAlive := make(chan string)
	// Initialize the solver
	solver, err := NewUniswapV3Solver(rpcURL, chainId, uniswapV3FactoryAddress, npmAddress)
	if err != nil {
		panic(err)
	}

	// Start the solver's API server
	go startHTTPServer(solver, httpPort)

	// Keep the process running
	<-keepAlive
}

// NewUniswapV3Solver creates a new instance of UniswapV3Solver
func NewUniswapV3Solver(rpcURL string, chainId int64, factoryAddress string, npmAddress string) (*UniswapV3Solver, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(npmABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	solver := &UniswapV3Solver{
		client:          client,
		chainId:         big.NewInt(chainId),
		factory:         common.HexToAddress(factoryAddress),
		npm:             common.HexToAddress(npmAddress),
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

func (s *UniswapV3Solver) cleanupExpiredParams() {
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

func (s *UniswapV3Solver) StopCleanup() {
	s.cleanupMutex.Lock()
	close(s.cleanupStop)
	s.cleanupMutex.Unlock()
}

// Solve executes the Uniswap V3 operation with a signed transaction
func (s *UniswapV3Solver) Solve(intent Intent, opIndex int, signature string) (string, error) {
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
		To:        &s.npm,
		Value:     big.NewInt(0),
		Gas:       params.GasLimit,
		GasTipCap: params.GasTipCap,
		GasFeeCap: params.GasFeeCap,
		Data:      params.Calldata,
	})

	signer := types.NewLondonSigner(s.chainId)
	txHash := signer.Hash(tx)
	fmt.Printf("Solve tx hash: %s\n", txHash.Hex()) // Debug

	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return "", fmt.Errorf("invalid signature: %v", err)
	}
	if len(sig) != 65 {
		return "", fmt.Errorf("invalid signature length: got %d, want 65", len(sig))
	}

	signedTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return "", fmt.Errorf("failed to add signature: %v", err)
	}

	sender, err := signer.Sender(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to recover sender: %v", err)
	}
	if sender != common.HexToAddress(intent.Identity) {
		return "", fmt.Errorf("signature mismatch: sender %s != identity %s", sender.Hex(), intent.Identity)
	}

	err = s.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	txHashStr := signedTx.Hash().Hex()
	s.statuses[txHashStr] = &TransactionStatus{
		TxHash: txHashStr,
		Status: "pending",
	}

	return txHashStr, nil
}

// Status checks the status of a Uniswap V3 operation
func (s *UniswapV3Solver) Status(intent Intent, opIndex int) (string, error) {
	if opIndex >= len(intent.Operations) {
		return "", fmt.Errorf("operation index out of range")
	}

	operation := intent.Operations[opIndex]
	// if operation.Result == "" {
	// 	return "pending", nil
	// }

	value, exists := s.statuses[operation.Result]
	if !exists {
		return "", fmt.Errorf("transaction not found")
	}

	status := value
	if status.Status == "pending" {
		receipt, err := s.client.TransactionReceipt(context.Background(), common.HexToHash(status.TxHash))
		if err == ethereum.NotFound {
			return "pending", nil
		}
		if err != nil {
			return "", fmt.Errorf("failed to get transaction receipt: %v", err)
		}

		if receipt.Status == 1 {
			status.Status = "success"
		} else {
			status.Status = "failure"
		}
		status.Receipt = receipt
		s.statuses[status.TxHash] = status
	}

	return status.Status, status.Error
}

// GetOutput retrieves the result of a Uniswap V3 operation
func (s *UniswapV3Solver) GetOutput(intent Intent, opIndex int) (string, error) {
	if opIndex >= len(intent.Operations) {
		return "", fmt.Errorf("operation index out of range")
	}

	operation := intent.Operations[opIndex]
	value, exists := s.statuses[operation.Result]
	if !exists {
		return "", fmt.Errorf("transaction not found")
	}

	status := value
	if status.Status != "success" {
		return "", fmt.Errorf("transaction not successful")
	}

	output := LPOutput{
		TxHash: status.TxHash,
	}

	receipt := status.Receipt
	if receipt == nil {
		return "", fmt.Errorf("transaction receipt not available")
	}

	var metadata LPMetadata
	if err := json.Unmarshal(operation.SolverMetadata, &metadata); err != nil {
		return "", fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	switch metadata.Action {
	case "mint":
		found := false
		for _, log := range receipt.Logs {
			if log.Address == s.npm {
				event, err := s.abi.EventByID(log.Topics[0])
				if err != nil {
					continue
				}
				if event.Name == "IncreaseLiquidity" {
					if len(log.Topics) < 2 {
						return "", fmt.Errorf("invalid IncreaseLiquidity event: insufficient topics")
					}
					tokenId := new(big.Int).SetBytes(log.Topics[1][:])
					dataMap, err := s.abi.Unpack(event.Name, log.Data)
					if err != nil {
						return "", fmt.Errorf("failed to unpack IncreaseLiquidity event data: %v", err)
					}
					liquidity, ok := dataMap[0].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse liquidity from event data")
					}
					amount0, ok := dataMap[1].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse amount0 from event data")
					}
					amount1, ok := dataMap[2].(*big.Int)
					if !ok {
						return "", fmt.Errorf("failed to parse amount1 from event data")
					}
					output.TokenId = uint(tokenId.Uint64())
					output.Liquidity = liquidity.String()
					output.AmountA = amount0.String()
					output.AmountB = amount1.String()
					found = true
					break
				}
			}
		}
		if !found {
			return "", fmt.Errorf("IncreaseLiquidity event not found in receipt")
		}

	case "exit":
		found := false
		for _, log := range receipt.Logs {
			if log.Address == s.npm {
				event, err := s.abi.EventByID(log.Topics[0])
				if err != nil {
					continue
				}
				if event.Name == "DecreaseLiquidity" {
					if len(log.Topics) < 2 {
						return "", fmt.Errorf("invalid DecreaseLiquidity event: insufficient topics")
					}
					tokenId := new(big.Int).SetBytes(log.Topics[1][:])
					dataMap, err := s.abi.Unpack(event.Name, log.Data)
					if err != nil {
						return "", fmt.Errorf("failed to unpack DecreaseLiquidity event data: %v", err)
					}
					liquidity, ok := dataMap[0].(*big.Int) // uint128 as *big.Int
					if !ok {
						return "", fmt.Errorf("failed to parse liquidity from event data")
					}
					amount0, ok := dataMap[1].(*big.Int) // uint256
					if !ok {
						return "", fmt.Errorf("failed to parse amount0 from event data")
					}
					amount1, ok := dataMap[2].(*big.Int) // uint256
					if !ok {
						return "", fmt.Errorf("failed to parse amount1 from event data")
					}
					output.TokenId = uint(tokenId.Uint64())
					output.Liquidity = liquidity.String()
					output.AmountA = amount0.String()
					output.AmountB = amount1.String()
					found = true
					break
				}
			}
		}
		if !found {
			return "", fmt.Errorf("DecreaseLiquidity event not found in receipt")
		}
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal output: %v", err)
	}

	return string(outputJSON), nil
}

// Construct builds the transaction data for a Uniswap V3 operation
func (s *UniswapV3Solver) Construct(intent Intent, opIndex int) (string, error) {
	if opIndex >= len(intent.Operations) {
		return "", fmt.Errorf("operation index out of range")
	}

	operation := intent.Operations[opIndex]
	var metadata LPMetadata
	if err := json.Unmarshal(operation.SolverMetadata, &metadata); err != nil {
		return "", fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	calldata, err := s.constructTxData(metadata, intent.Identity)
	if err != nil {
		return "", fmt.Errorf("failed to construct transaction data: %v", err)
	}

	nonce, err := s.client.PendingNonceAt(context.Background(), common.HexToAddress(intent.Identity))
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	msg := ethereum.CallMsg{
		From:      common.HexToAddress(intent.Identity),
		To:        &s.npm,
		Gas:       0,
		GasPrice:  gasPrice,
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		GasTipCap: gasPrice,
		Value:     big.NewInt(0),
		Data:      calldata,
	}
	gasLimit, err := s.client.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %v", err)
	}
	gasLimit = gasLimit + (gasLimit / 5)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   s.chainId,
		Nonce:     nonce,
		To:        &s.npm,
		Value:     big.NewInt(0),
		Gas:       gasLimit,
		GasTipCap: gasPrice,
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		Data:      calldata,
	})

	signer := types.NewLondonSigner(s.chainId)
	txHash := signer.Hash(tx)

	// Calculate operation hash
	opBytes, err := json.Marshal(operation)
	if err != nil {
		return "", fmt.Errorf("failed to marshal operation: %v", err)
	}
	opHash := "0x" + hex.EncodeToString(crypto.Keccak256(opBytes))

	// Store parameters with 5-minute deadline
	s.txParams[opHash] = TxParams{
		Calldata:  calldata,
		Nonce:     nonce,
		GasLimit:  gasLimit,
		GasTipCap: gasPrice,
		GasFeeCap: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		Deadline:  time.Now().Add(5 * time.Minute),
	}

	return "0x" + hex.EncodeToString(txHash[:]), nil
}

// constructTxData builds the transaction data based on the operation type
func (s *UniswapV3Solver) constructTxData(metadata LPMetadata, identity string) ([]byte, error) {
	switch metadata.Action {
	case "mint":
		return s.constructMint(metadata, identity)
	case "exit":
		return s.constructExit(metadata, identity)
	default:
		return nil, fmt.Errorf("unsupported action: %s", metadata.Action)
	}
}

// constructMint builds transaction data for minting a new position
func (s *UniswapV3Solver) constructMint(metadata LPMetadata, identity string) ([]byte, error) {
	// Create MintParams struct matching the ABI tuple structure
	type MintParams struct {
		Token0         common.Address
		Token1         common.Address
		Fee            *big.Int // Changed from uint32 to *big.Int for uint24
		TickLower      *big.Int // Changed from int32 to *big.Int for int24
		TickUpper      *big.Int // Changed from int32 to *big.Int for int24
		Amount0Desired *big.Int
		Amount1Desired *big.Int
		Amount0Min     *big.Int
		Amount1Min     *big.Int
		Recipient      common.Address
		Deadline       *big.Int
	}

	amount0Desired := new(big.Int)
	amount0Desired.SetString(metadata.AmountA, 10)
	amount1Desired := new(big.Int)
	amount1Desired.SetString(metadata.AmountB, 10)

	// Convert Fee, TickLower, and TickUpper to *big.Int
	fee := big.NewInt(int64(metadata.Fee))
	tickLower := big.NewInt(int64(metadata.TickLower))
	tickUpper := big.NewInt(int64(metadata.TickUpper))

	params := MintParams{
		Token0:         common.HexToAddress(metadata.TokenA),
		Token1:         common.HexToAddress(metadata.TokenB),
		Fee:            fee,
		TickLower:      tickLower,
		TickUpper:      tickUpper,
		Amount0Desired: amount0Desired,
		Amount1Desired: amount1Desired,
		Amount0Min:     big.NewInt(0),
		Amount1Min:     big.NewInt(0),
		Recipient:      common.HexToAddress(identity),
		Deadline:       big.NewInt(time.Now().Unix() + 1800), // 30 minutes from now
	}

	// Pack as a tuple argument named "params"
	// return s.abi.Pack("mint", struct{ Params MintParams }{Params: params})
	return s.abi.Pack("mint", params)
}

// constructExit builds transaction data for exiting a position
func (s *UniswapV3Solver) constructExit(metadata LPMetadata, identity string) ([]byte, error) {
	// Create DecreaseLiquidityParams struct matching the ABI structure
	type DecreaseLiquidityParams struct {
		TokenId    *big.Int
		Liquidity  *big.Int
		Amount0Min *big.Int
		Amount1Min *big.Int
		Deadline   *big.Int
	}

	tokenId := new(big.Int)
	tokenId.SetUint64(uint64(metadata.TokenId))
	liquidity, err := s.getPositionLiquidity(metadata.TokenId)
	if err != nil {
		// Fallback to a default if needed
		liquidity = big.NewInt(1000000000000000000)
	}

	params := DecreaseLiquidityParams{
		TokenId:    tokenId,
		Liquidity:  liquidity,
		Amount0Min: big.NewInt(0),
		Amount1Min: big.NewInt(0),
		Deadline:   big.NewInt(time.Now().Unix() + 1800), // 30 minutes from now
	}

	// Pack as a tuple argument named "params"
	return s.abi.Pack("decreaseLiquidity", params)
}

func (s *UniswapV3Solver) getPositionLiquidity(tokenId uint) (*big.Int, error) {

	positionABI := `[{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"positions","outputs":[{"internalType":"uint96","name":"nonce","type":"uint96"},{"internalType":"address","name":"operator","type":"address"},{"internalType":"address","name":"token0","type":"address"},{"internalType":"address","name":"token1","type":"address"},{"internalType":"uint24","name":"fee","type":"uint24"},{"internalType":"int24","name":"tickLower","type":"int24"},{"internalType":"int24","name":"tickUpper","type":"int24"},{"internalType":"uint128","name":"liquidity","type":"uint128"},{"internalType":"uint256","name":"feeGrowthInside0LastX128","type":"uint256"},{"internalType":"uint256","name":"feeGrowthInside1LastX128","type":"uint256"},{"internalType":"uint128","name":"tokensOwed0","type":"uint128"},{"internalType":"uint128","name":"tokensOwed1","type":"uint128"}],"stateMutability":"view","type":"function"}]`

	parsedABI, err := abi.JSON(strings.NewReader(positionABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse positions ABI: %v", err)
	}

	callData, err := parsedABI.Pack("positions", big.NewInt(int64(tokenId)))
	if err != nil {
		return nil, fmt.Errorf("failed to pack positions call: %v", err)
	}

	npmAddress := s.npm

	result, err := s.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &npmAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call positions method: %v", err)
	}

	var positionInfo struct {
		Nonce                    *big.Int
		Operator                 common.Address
		Token0                   common.Address
		Token1                   common.Address
		Fee                      *big.Int
		TickLower                *big.Int
		TickUpper                *big.Int
		Liquidity                *big.Int
		FeeGrowthInside0LastX128 *big.Int
		FeeGrowthInside1LastX128 *big.Int
		TokensOwed0              *big.Int
		TokensOwed1              *big.Int
	}

	err = parsedABI.UnpackIntoInterface(&positionInfo, "positions", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack position info: %v", err)
	}

	return positionInfo.Liquidity, nil
}
