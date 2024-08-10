package bridge

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func BridgeBurnDataToSign(
	rpcURL string,
	bridgeContractAddress string,
	account string,
	amount string,
	token string,
) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}

	instance, err := NewBridge(common.HexToAddress(bridgeContractAddress), client)
	if err != nil {
		return "", err
	}

	_amountIn, _ := new(big.Int).SetString(amount, 10)

	nonce, err := instance.Nonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		return "", err
	}

	messageHash, err := instance.GetBurnMessageHash(
		&bind.CallOpts{},
		common.HexToAddress(account),
		nonce,
		_amountIn,
		common.HexToAddress(token),
	)

	if err != nil {
		return "", err
	}

	hash, err := instance.GetEthSignedMessageHash(
		&bind.CallOpts{},
		messageHash,
	)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash[:]), nil
}
