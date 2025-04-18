package blockchains

import (
	"errors"
	"net/http"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/stellar/go/clients/horizonclient"
)

// NewBitcoinBlockchain creates a new Stellar blockchain instance
func NewBitcoinBlockchain(networkType NetworkType) (IBlockchain, error) {
	network := Network{
		networkType: networkType,
		nodeURL:     horizonclient.DefaultPublicNetClient.HorizonURL,
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = horizonclient.DefaultTestNetClient.HorizonURL
		network.networkID = "testnet"
	}

	client := &horizonclient.Client{
		HorizonURL: network.nodeURL,
		HTTP:       http.DefaultClient,
	}
	// Set timeout using the SDK's constant
	client.SetHorizonTimeout(horizonclient.HorizonTimeout)

	return &BitcoinBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Bitcoin,
			network:         network,
			keyCurve:        common.CurveEcdsa,
			signingEncoding: "hex",
			decimals:        7,
			opTimeout:       time.Second * 10,
		},
		client: client,
	}, nil
}

// This is a type assertion to ensure that the BitcoinBlockchain implements the IBlockchain interface
var _ IBlockchain = &BitcoinBlockchain{}

// BitcoinBlockchain implements the IBlockchain interface for Stellar
type BitcoinBlockchain struct {
	BaseBlockchain
	client *horizonclient.Client
}

func (b *BitcoinBlockchain) BroadcastTransaction(txn string, signatureHex string, _ *string) (string, error) {
	return "", errors.New("BroadcastTransaction not implemented")
}

func (b *BitcoinBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	return nil, errors.New("GetTransfers not implemented")
}

func (b *BitcoinBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	return false, errors.New("IsTransactionBroadcastedAndConfirmed not implemented")
}

func (b *BitcoinBlockchain) BuildWithdrawTx(bridgeAddress string,
	solverOutput string,
	userAddress string,
	tokenAddress *string,
) (string, string, error) {
	return "", "", errors.New("BuildWithdrawTx not implemented")
}

func (b *BitcoinBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *BitcoinBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}
