package stellar

import (
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
)

// GetClient returns a properly configured Stellar client based on chain settings
func GetClient(chainType string, chainUrl string) *horizonclient.Client {
	if chainType == "mainnet" {
		return horizonclient.DefaultPublicNetClient
	} else if chainUrl != "" {
		// Create a new client with custom URL
		client := &horizonclient.Client{
			HorizonURL: chainUrl,
			HTTP:      http.DefaultClient,
		}
		// Set timeout using the SDK's constant
		client.SetTimeout(horizonclient.HorizonTimeout)
		return client
	}
	return horizonclient.DefaultTestNetClient
}
