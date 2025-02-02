package sequencer

import (
	"context"
	"math/big"
	"strings"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetEthereumTransfers(chainId string, txnHash string, ecdsaAddr string) ([]common.Transfer, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(chain.ChainUrl)
	if err != nil {
		return nil, err
	}

	receipt, err := client.TransactionReceipt(context.Background(), ethCommon.HexToHash(txnHash))
	if err != nil {
		return nil, err
	}

	var transfers []common.Transfer

	for _, log := range receipt.Logs {
		if log.Topics[0].Hex() == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
			from := ethCommon.BytesToAddress(log.Topics[1].Bytes()).Hex()
			to := ethCommon.BytesToAddress(log.Topics[2].Bytes()).Hex()

			decimal, symbol, err := getERC20Details(client, ethCommon.BytesToAddress(log.Address.Bytes()))

			if err != nil {
				return nil, err
			}

			formattedAmount, err := FormatUnits(new(big.Int).SetBytes(log.Data), int(decimal))

			if err != nil {
				return nil, err
			}

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

	tx, _, err := client.TransactionByHash(context.Background(), ethCommon.HexToHash(txnHash))
	if err != nil {
		return nil, err
	}

	wei := tx.Value()

	if wei.Cmp(big.NewInt(0)) != 0 {
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

	return transfers, nil
}

func FormatUnits(value *big.Int, decimals int) (string, error) {
	// Create the scaling factor as 10^decimals
	scalingFactor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// Convert the value to a big.Float
	valueFloat := new(big.Float).SetInt(value)

	// Divide the value by the scaling factor
	result := new(big.Float).Quo(valueFloat, scalingFactor)

	// Convert the result to a string with the appropriate precision
	return result.Text('f', decimals), nil
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
