package ERC20

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetBalance(rpcURL string, token string, account string) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}

	instance, err := NewERC20(common.HexToAddress(token), client)
	if err != nil {
		return "", err
	}

	balance, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(account))

	if err != nil {
		return "", err
	}

	return balance.String(), nil
}
