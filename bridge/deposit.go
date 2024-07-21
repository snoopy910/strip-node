package bridge

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func BridgeDepositDataToSign(rpcURL string, bridgeContractAddress string, amount string, account string, token string) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}

	instance, err := NewBridge(common.HexToAddress(bridgeContractAddress), client)
	if err != nil {
		return "", err
	}

	amountBigInt, _ := big.NewInt(0).SetString(amount, 10)
	nonce, err := instance.MintNonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		return "", err
	}

	messageHash, err := instance.GetMessageHash(
		&bind.CallOpts{},
		common.HexToAddress(account),
		nonce,
		amountBigInt,
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
