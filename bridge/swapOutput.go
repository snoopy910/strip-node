package bridge

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetSwapOutput(
	rpcURL string,
	txnHash string,
) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}

	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txnHash))

	if err != nil {
		return "", err
	}

	contractABI := `[
		{
			"anonymous": false,
			"inputs": [
				{"indexed": false, "name": "user", "type": "address"},
				{"indexed": false, "name": "tokenIn", "type": "address"},
				{"indexed": false, "name": "tokenOut", "type": "address"},
				{"indexed": false, "name": "amountIn", "type": "uint256"},
				{"indexed": false, "name": "amountOut", "type": "uint256"}
			],
			"name": "TokenSwapped",
			"type": "event"
		}
	]`

	type TokenSwappedEvent struct {
		User      common.Address
		TokenIn   common.Address
		TokenOut  common.Address
		AmountIn  *big.Int
		AmountOut *big.Int
	}

	contractAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return "", err
	}

	for _, vLog := range receipt.Logs {
		event := new(TokenSwappedEvent)
		err := contractAbi.UnpackIntoInterface(event, "TokenSwapped", vLog.Data)
		if err != nil {
			continue
		}

		fmt.Printf("User: %s\n", event.User.Hex())
		fmt.Printf("TokenIn: %s\n", event.TokenIn.Hex())
		fmt.Printf("TokenOut: %s\n", event.TokenOut.Hex())
		fmt.Printf("AmountIn: %s\n", event.AmountIn.String())
		fmt.Printf("AmountOut: %s\n", event.AmountOut.String())

		return event.AmountOut.String(), nil
	}

	return "", errors.New("TokenSwapped event not found")
}
