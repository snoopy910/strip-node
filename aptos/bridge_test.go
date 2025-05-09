package aptos

import (
	"fmt"
	"strings"
	"testing"
)

var (
	ErrAccNotFound   = "failed to get accountInfo"
	ErrInvalidFormat = "invalid token address format"
)

func TestWithdrawAptosNativeGetSignature(t *testing.T) {
	// Test cases
	tests := []struct {
		Name      string
		RpcURL    string
		Account   string
		Amount    string
		Recipient string
		IsError   bool
		Error     string
	}{
		{
			Name:      "Native APT withdraw on mainnet",
			RpcURL:    "https://fullnode.mainnet.aptoslabs.com",
			Account:   "0xe6f6ecc5eb1d93963c9d9446e5f33d677facd01511cf8bab18290c67641ef083",
			Amount:    "1000",
			Recipient: "0xf6dc6b07e7c49f240d197d1880f2806b7dfa3d838762fb3b9391f707bc6afbf7",
			IsError:   false,
		},
		{
			Name:      "Account Not Found",
			RpcURL:    "https://fullnode.mainnet.aptoslabs.com",
			Account:   "0xe93dcd3dd5febf8d72bf8d33e1d85a6115300fedfe055c062834d264f103ce4c",
			Amount:    "1000",
			Recipient: "0xf6dc6b07e7c49f240d197d1880f2806b7dfa3d838762fb3b9391f707bc6afbf7",
			IsError:   true,
			Error:     ErrAccNotFound,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			serializedTxn, dataToSign, err := WithdrawAptosNativeGetSignature(tt.RpcURL, tt.Account, tt.Amount, tt.Recipient)
			if !tt.IsError && err != nil {
				t.Errorf("WithdrawAptosNativeGetSignature() error = %v", err)
			} else if tt.IsError && strings.Split(err.Error(), ":")[0] != tt.Error {
				t.Errorf("expected error = %v", err)
			}
			fmt.Println(serializedTxn)
			fmt.Println(dataToSign)
		})
	}
}

func TestWithdrawAptosTokenGetSignature(t *testing.T) {
	// Test cases
	tests := []struct {
		Name      string
		RpcURL    string
		Account   string
		Amount    string
		Recipient string
		TokenAddr string
		IsError   bool
		Error     string
	}{
		{
			Name:      "Aptos Token withdraw on mainnet",
			RpcURL:    "https://fullnode.mainnet.aptoslabs.com",
			Account:   "0xe6f6ecc5eb1d93963c9d9446e5f33d677facd01511cf8bab18290c67641ef083",
			Amount:    "1000",
			Recipient: "0xf6dc6b07e7c49f240d197d1880f2806b7dfa3d838762fb3b9391f707bc6afbf7",
			TokenAddr: "0x9d0595765a31f8d56e1d2aafc4d6c76f283c67a074ef8812d8c31bd8252ac2c3::asset::TOMA",
			IsError:   false,
		},
		{
			Name:      "Invalid Token Address",
			RpcURL:    "https://fullnode.mainnet.aptoslabs.com",
			Account:   "0xe6f6ecc5eb1d93963c9d9446e5f33d677facd01511cf8bab18290c67641ef083",
			Amount:    "1000",
			Recipient: "0xf6dc6b07e7c49f240d197d1880f2806b7dfa3d838762fb3b9391f707bc6afbf7",
			TokenAddr: "0x9d0595765a31f8d56e1d2aafc4d6c76f283c67a074ef8812d8c31bd8252ac2c3::asset",
			IsError:   true,
			Error:     ErrInvalidFormat,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			serializedTxn, dataToSign, err := WithdrawAptosTokenGetSignature(tt.RpcURL, tt.Account, tt.Amount, tt.Recipient, tt.TokenAddr)
			fmt.Println(err)
			if !tt.IsError && err != nil {
				t.Errorf("WithdrawAptosTokenGetSignature() error = %v", err)
			} else if tt.IsError && err.Error() != tt.Error {
				t.Errorf("expected error = %v", err)
			}
			fmt.Println(serializedTxn)
			fmt.Println(dataToSign)
		})
	}
}
