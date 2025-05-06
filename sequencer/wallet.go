package sequencer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	pb "github.com/StripChain/strip-node/libs/proto"
	"github.com/StripChain/strip-node/util/logger"
)

// createWallet creates a new wallet with the specified identity and blockchain ID.
// Note: It selects a list of signers, ensuring the number of signers does not exceed MaximumSigners.
// The function performs the following steps:
// 1. Selects a random subset of signers if the total number exceeds MaximumSigners.
// 2. Determines the required key curve based on the blockchain.
// 3. Sends a gRPC request to the first signer to generate keys.
// 4. Waits for the key generation to complete.
//
// Parameters:
// - identity: A string representing the identity for the wallet.
// - blockchainID: The blockchain identifier for which to create the wallet.
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

	protoCurve, err := libs.CommonCurveToProto(blockchain.KeyCurve())
	if err != nil {
		return fmt.Errorf("failed to convert curve to proto: %w", err)
	}

	client, err := validatorClientManager.GetClient(signers[0].URL)
	if err != nil {
		return fmt.Errorf("failed to get validator client: %w", err)
	}

	_, err = client.Keygen(context.Background(), &pb.KeygenRequest{
		Identity:      identity,
		IdentityCurve: protoCurve,
		Signers:       signersPublicKeyList,
	})
	if err != nil {
		return fmt.Errorf("failed to keygen: %w", err)
	}

	resp, err := client.GetAddresses(context.Background(), &pb.GetAddressesRequest{
		Identity:      identity,
		IdentityCurve: protoCurve,
	})
	if err != nil {
		return fmt.Errorf("failed to get addresses: %w", err)
	}

	wallet := db.WalletSchema{
		Identity:       identity,
		BlockchainID:   blockchainID,
		Signers:        signersPublicKeyList,
		EDDSAPublicKey: resp.EddsaAddress,
		ECDSAPublicKey: resp.EcdsaAddress,
	}

	for protoBlockchainID, addresses := range resp.Addresses {
		blockchainID, err := libs.ProtoToBlockchainsID(pb.BlockchainID(protoBlockchainID))
		if err != nil {
			return fmt.Errorf("failed to convert proto blockchain ID to blockchain ID: %w", err)
		}

		switch blockchainID {
		case blockchains.Sui:
			wallet.SuiPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
		case blockchains.Algorand:
			wallet.AlgorandEDDSAPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
		case blockchains.Bitcoin:
			wallet.BitcoinMainnetPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
			wallet.BitcoinTestnetPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_TESTNET)].Address
			wallet.BitcoinRegtestPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_REGNET)].Address
		case blockchains.Dogecoin:
			wallet.DogecoinMainnetPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
			wallet.DogecoinTestnetPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_TESTNET)].Address
		case blockchains.Stellar:
			wallet.StellarPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
		case blockchains.Ripple:
			wallet.RippleEDDSAPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
		case blockchains.Cardano:
			wallet.CardanoPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
		case blockchains.Aptos:
			wallet.AptosEDDSAPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
		case blockchains.Solana:
			wallet.SolanaPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
		default:
			if blockchains.IsEVMBlockchain(blockchainID) {
				logger.Sugar().Infof("%s address: %s", blockchainID, addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address)
				wallet.EthereumPublicKey = addresses.NetworkAddresses[int32(pb.NetworkType_MAINNET)].Address
			} else {
				logger.Sugar().Errorw("unsupported blockchain ID", "blockchainID", blockchainID)
				return fmt.Errorf("unsupported blockchain ID: %s", blockchainID)
			}
		}
	}

	// add created wallet to the store
	_, err = db.AddWallet(&wallet)
	if err != nil {
		return fmt.Errorf("failed to add wallet to store: %w", err)
	}

	return nil
}
