package sequencer

import (
	"fmt"
	"math/big"

	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum/common"
)

// Helper functions for BURN_SYNTHETIC operation debugging

// validateBurnSyntheticToken checks if the token address is valid and logs detailed information
func validateBurnSyntheticToken(tokenAddress string) (bool, error) {
	logger.Sugar().Infow("Validating BURN_SYNTHETIC token address", "tokenAddress", tokenAddress)
	
	if !common.IsHexAddress(tokenAddress) {
		logger.Sugar().Errorw("BURN_SYNTHETIC invalid token address format", "tokenAddress", tokenAddress)
		return false, fmt.Errorf("invalid token address format: %s", tokenAddress)
	}
	
	// Additional validation could be added here (e.g., checking if it's in a whitelist)
	logger.Sugar().Infow("BURN_SYNTHETIC token address validation successful", "tokenAddress", tokenAddress)
	return true, nil
}

// logBurnSyntheticSignature logs detailed information about the signature
func logBurnSyntheticSignature(signature string, dataToSign string) {
	if len(signature) < 20 {
		logger.Sugar().Errorw("BURN_SYNTHETIC signature too short", "signatureLength", len(signature))
		return
	}
	
	logger.Sugar().Infow("BURN_SYNTHETIC signature details",
		"signatureLength", len(signature),
		"signaturePrefix", signature[:20]+"...",
		"dataToSignLength", len(dataToSign),
		"dataToSignPrefix", dataToSign[:min(20, len(dataToSign))]+"...")
}

// verifyBurnSyntheticSignature performs local verification of the signature
func verifyBurnSyntheticSignature(dataToSign string, signature string, publicKey string) (bool, error) {
	logger.Sugar().Infow("Verifying BURN_SYNTHETIC signature locally",
		"publicKey", publicKey,
		"signatureLength", len(signature),
		"dataToSignLength", len(dataToSign))
	
	// This is a placeholder for actual verification logic
	// In a real implementation, this would use the appropriate crypto library
	// to verify the signature against the data and public key
	
	// For now, just do basic format checking
	if len(signature) < 64 {
		logger.Sugar().Errorw("BURN_SYNTHETIC signature verification failed - signature too short",
			"signatureLength", len(signature))
		return false, fmt.Errorf("signature too short: %d bytes", len(signature))
	}
	
	logger.Sugar().Infow("BURN_SYNTHETIC signature basic format check passed")
	return true, nil
}

// logBurnSyntheticBalanceCheck logs detailed balance information
func logBurnSyntheticBalanceCheck(balanceBig *big.Int, amountBig *big.Int, tokenAddress string) {
	logger.Sugar().Infow("BURN_SYNTHETIC balance check",
		"balance", balanceBig.String(),
		"requestedAmount", amountBig.String(),
		"sufficient", balanceBig.Cmp(amountBig) >= 0,
		"tokenAddress", tokenAddress)
}

// logBurnSyntheticGasEstimation logs detailed gas estimation information
func logBurnSyntheticGasEstimation(estimatedGas uint64, gasPrice *big.Int, tokenAddress string) {
	gasCost := new(big.Int).Mul(big.NewInt(int64(estimatedGas)), gasPrice)
	
	logger.Sugar().Infow("BURN_SYNTHETIC gas estimation",
		"estimatedGas", estimatedGas,
		"gasPrice", gasPrice.String(),
		"totalGasCost", gasCost.String(),
		"tokenAddress", tokenAddress)
}
