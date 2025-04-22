package sequencer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/util/logger"
)

// createWallet creates a new wallet with the specified identity and identity curve.
// Note: It selects a list of signers, ensuring the number of signers does not exceed MaximumSigners.
// The function performs the following steps:
// 1. Selects a random subset of signers if the total number exceeds MaximumSigners.
// 2. Constructs a CreateWalletRequest for both "eddsa" and "ecdsa" key curves.
// 3. Sends HTTP requests to the first signer's URL to generate keys.
// 4. Retrieves the generated addresses for both key curves.
// 5. Constructs a WalletSchema and adds it to the wallet store.
//
// Parameters:
// - identity: A string representing the identity for the wallet.
// - identityCurve: A string representing the curve type for the identity.
//
// Returns:
// - error: An error if any step in the wallet creation process fails.
func createWallet(identity string, blockchainID blockchains.BlockchainID) error {
	// select a list of nodes.
	// If length of selected nodes is more than maximum nodes then use maximum nodes length as signers.
	// If length of selected nodes is less than maximum nodes then use all nodes as signers.

	signers, err := SignersList()
	if err != nil {
		return fmt.Errorf("failed to get signers: %w", err)
	}

	if len(signers) > MaximumSigners {
		// select random number of max signers
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(signers), func(i, j int) { signers[i], signers[j] = signers[j], signers[i] })
		signers = signers[:MaximumSigners]
	}

	signersPublicKeyList := make([]string, len(signers))
	for i, signer := range signers {
		signersPublicKeyList[i] = signer.PublicKey
	}

	blockchain, err := blockchains.GetBlockchain(blockchainID, blockchains.NetworkType(blockchains.Mainnet))
	if err != nil {
		return fmt.Errorf("failed to get blockchain: %w", err)
	}

	// create the wallet whose keycurve is eddsa here
	createWalletRequest := libs.CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: blockchain.KeyCurve(),
		KeyCurve:      common.CurveEcdsa,
		Signers:       signersPublicKeyList,
	}

	marshalled, err := json.Marshal(createWalletRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal create wallet request: %w", err)
	}

	req, err := http.NewRequest("GET", signers[0].URL+"/keygen", bytes.NewReader(marshalled))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 3 * time.Minute}
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	// create the wallet whose keycurve is sui_eddsa here
	createWalletRequest = libs.CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: blockchain.KeyCurve(),
		KeyCurve:      common.CurveEddsa,
		Signers:       signersPublicKeyList,
	}

	marshalled, err = json.Marshal(createWalletRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal create wallet request: %w", err)
	}

	req, err = http.NewRequest("GET", signers[0].URL+"/keygen", bytes.NewReader(marshalled))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client = http.Client{Timeout: 3 * time.Minute}
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create sui wallet: %w", err)
	}

	// get the address of the wallet whose keycurve is eddsa here
	resp, err := http.Get(fmt.Sprintf("%s/address?identity=%s&identityCurve=%s", signers[0].URL, identity, blockchain.KeyCurve()))
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	var addressesResponse libs.AddressesResponse
	logger.Sugar().Infof("addresses response: %s", string(body))
	err = json.Unmarshal(body, &addressesResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal addresses response: %w", err)
	}

	wallet := db.WalletSchema{
		Identity:       identity,
		BlockchainID:   blockchainID,
		Signers:        signersPublicKeyList,
		EDDSAPublicKey: addressesResponse.EDDSAAddress,
		ECDSAPublicKey: addressesResponse.ECDSAAddress,
	}

	for blockchainID, addresses := range addressesResponse.Addresses {
		switch blockchainID {
		case blockchains.Sui:
			wallet.SuiPublicKey = addresses[blockchains.Mainnet]
		case blockchains.Algorand:
			wallet.AlgorandEDDSAPublicKey = addresses[blockchains.Mainnet]
		case blockchains.Bitcoin:
			wallet.BitcoinMainnetPublicKey = addresses[blockchains.Mainnet]
			wallet.BitcoinTestnetPublicKey = addresses[blockchains.Testnet]
			wallet.BitcoinRegtestPublicKey = addresses[blockchains.Regnet]
		case blockchains.Dogecoin:
			wallet.DogecoinMainnetPublicKey = addresses[blockchains.Mainnet]
			wallet.DogecoinTestnetPublicKey = addresses[blockchains.Testnet]
		case blockchains.Stellar:
			wallet.StellarPublicKey = addresses[blockchains.Mainnet]
		case blockchains.Ripple:
			wallet.RippleEDDSAPublicKey = addresses[blockchains.Mainnet]
		case blockchains.Cardano:
			wallet.CardanoPublicKey = addresses[blockchains.Mainnet]
		case blockchains.Aptos:
			wallet.AptosEDDSAPublicKey = addresses[blockchains.Mainnet]
		case blockchains.Solana:
			wallet.SolanaPublicKey = addresses[blockchains.Mainnet]
		default:
			if blockchains.IsEVMBlockchain(blockchainID) {
				logger.Sugar().Infof("%s address: %s", blockchainID, addresses[blockchains.Mainnet])
				wallet.EthereumPublicKey = addresses[blockchains.Mainnet]
			}
			logger.Sugar().Errorw("unsupported blockchain ID", "blockchainID", blockchainID)
			return fmt.Errorf("unsupported blockchain ID: %s", blockchainID)
		}
	}

	// add created wallet to the store
	_, err = db.AddWallet(&wallet)
	if err != nil {
		return fmt.Errorf("failed to add wallet to store: %w", err)
	}

	return nil
}
