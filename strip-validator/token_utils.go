package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/util/logger"
)

// VerifyTokenBalance checks if a wallet has sufficient token balance for a given token and amount
// Returns true if the wallet has enough balance, false otherwise
func VerifyTokenBalance(identity string, blockchainID blockchains.BlockchainID, token string, amount string) (bool, error) {
	// Get bridgewallet by calling /getwallet from sequencer api
	reqURL := fmt.Sprintf("%s/getWallet?identity=%s&blockchain=%s", SequencerHost, identity, blockchainID)
	logger.Sugar().Infow("Making request to get wallet",
		"fullURL", reqURL,
		"identity", identity,
		"blockchainID", blockchainID,
		"token", token)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		logger.Sugar().Errorw("error creating request", "error", err)
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Sugar().Errorw("error sending request", "error", err)
		return false, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Sugar().Errorw("error reading response body", "error", err)
		return false, err
	}

	// Log the response for debugging
	logger.Sugar().Debugw("Received response from sequencer",
		"status", resp.Status,
		"responseBody", string(body),
		"requestURL", reqURL)

	// Check if we got an error response
	if resp.StatusCode != http.StatusOK {
		logger.Sugar().Errorw("Unable to get wallet from sequencer",
			"status", resp.Status,
			"responseBody", string(body),
			"requestURL", reqURL)
		return false, fmt.Errorf("failed to get wallet: %s", resp.Status)
	}

	var bridgeWallet db.WalletSchema
	err = json.Unmarshal(body, &bridgeWallet)
	if err != nil {
		// Log the error with the actual response body
		logger.Sugar().Errorw("error unmarshalling response body",
			"error", err,
			"responseBody", string(body),
			"status", resp.Status)

		// Try to unmarshal as a map to see what we actually got
		var rawResponse map[string]interface{}
		if jsonErr := json.Unmarshal(body, &rawResponse); jsonErr == nil {
			logger.Sugar().Infow("Response parsed as generic JSON", "content", rawResponse)

			// Check if there's an error message in the response
			if errMsg, ok := rawResponse["error"]; ok {
				logger.Sugar().Errorw("Server returned error", "error", errMsg)
				return false, fmt.Errorf("server error: %v", errMsg)
			}
		}

		return false, err
	}

	// Check if we got a valid wallet
	if bridgeWallet.ECDSAPublicKey == "" {
		logger.Sugar().Errorw("Invalid wallet structure received", "wallet", bridgeWallet)
		return false, fmt.Errorf("invalid wallet structure: missing ECDSAPublicKey")
	}

	// Verify the user has sufficient token balance
	balance, err := ERC20.GetBalance(RPC_URL, token, bridgeWallet.ECDSAPublicKey)
	if err != nil {
		logger.Sugar().Errorw("Error getting token balance:", "error", err)
		return false, err
	}

	balanceBig, ok := new(big.Int).SetString(balance, 10)
	if !ok {
		logger.Sugar().Errorw("Error parsing balance")
		return false, fmt.Errorf("error parsing balance")
	}

	amountBig, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		logger.Sugar().Errorw("Error parsing amount")
		return false, fmt.Errorf("error parsing amount")
	}

	// Return true if the balance is sufficient (≥ amount)
	return balanceBig.Cmp(amountBig) >= 0, nil
}
