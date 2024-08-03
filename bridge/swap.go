package bridge

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func BridgeSwapDataToSign(
	rpcURL string,
	bridgeContractAddress string,
	account string,
	tokenIn string,
	tokenOut string,
	amountIn string,
	deadline int64,
) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}

	instance, err := NewBridge(common.HexToAddress(bridgeContractAddress), client)
	if err != nil {
		return "", err
	}

	_amountIn, _ := new(big.Int).SetString(amountIn, 10)

	params := ISwapRouterExactInputSingleParams{
		AmountIn:          _amountIn,
		AmountOutMinimum:  big.NewInt(0),
		TokenIn:           common.HexToAddress(tokenIn),
		TokenOut:          common.HexToAddress(tokenOut),
		Fee:               big.NewInt(500),
		Recipient:         common.HexToAddress(account),
		Deadline:          big.NewInt(0).SetInt64(deadline),
		SqrtPriceLimitX96: big.NewInt(0),
	}

	nonce, err := instance.Nonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		return "", err
	}

	messageHash, err := instance.GetSwapMessageHash(
		&bind.CallOpts{},
		params,
		common.HexToAddress(account),
		nonce,
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
