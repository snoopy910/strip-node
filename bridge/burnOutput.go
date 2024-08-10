package bridge

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetBurnOutput(
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

	eventSig := []byte("TokenBurned(address,address,uint256)")
	eventSigHash := crypto.Keccak256Hash(eventSig)

	contractABI := `[
		{
			"anonymous": false,
			"inputs": [
				{"indexed": false, "name": "user", "type": "address"},
				{"indexed": false, "name": "token", "type": "address"},
				{"indexed": false, "name": "amount", "type": "uint256"}
			],
			"name": "TokenBurned",
			"type": "event"
		}
	]`

	contractAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return "", err
	}

	for _, vLog := range receipt.Logs {
		if vLog.Topics[0].Hex() == eventSigHash.Hex() {
			event, err := contractAbi.Unpack("TokenBurned", vLog.Data)
			if err != nil {
				return "", err
			}

			value := event[2].(*big.Int)

			return value.String(), nil
		}
	}

	return "", errors.New("TokenBurned event not found")
}
