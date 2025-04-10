package sequencer

import (
	"fmt"
	"math/big"
	"strings"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

// Common ERC20 error signatures
const (
	// Error function selector: Error(string)
	ErrorSig = "08c379a0"
	
	// Custom error selectors
	InsufficientBalanceSig = "e450d38c"
	InsufficientAllowanceSig = "4bd67a2d"
	TransferFromZeroAddressSig = "ea553b34"
	TransferToZeroAddressSig = "d92e233d"
)

// Error ABI definitions for common ERC20 errors
const errorABI = `[
	{
		"inputs": [
			{
				"name": "reason",
				"type": "string"
			}
		],
		"name": "Error",
		"type": "error"
	},
	{
		"inputs": [
			{
				"name": "sender",
				"type": "address"
			},
			{
				"name": "balance",
				"type": "uint256"
			},
			{
				"name": "needed",
				"type": "uint256"
			}
		],
		"name": "InsufficientBalance",
		"type": "error"
	},
	{
		"inputs": [
			{
				"name": "owner",
				"type": "address"
			},
			{
				"name": "spender",
				"type": "address"
			},
			{
				"name": "allowance",
				"type": "uint256"
			},
			{
				"name": "needed",
				"type": "uint256"
			}
		],
		"name": "InsufficientAllowance",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "TransferFromZeroAddress",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "TransferToZeroAddress",
		"type": "error"
	}
]`

// DecodeERC20RevertReason attempts to decode Ethereum revert reason data into a human-readable format
// Uses ABI decoding for standard errors when possible, with fallback to manual decoding
func DecodeERC20RevertReason(errorMsg string) string {
	// Check if this is a revert error
	if !strings.Contains(errorMsg, "execution reverted") {
		return errorMsg
	}
	
	// Try to use ABI to decode the error first
	decodedABI := decodeErrorWithABI(errorMsg)
	if decodedABI != "" {
		return decodedABI
	}
	
	// Fall back to manual decoding if ABI approach failed
	return manualDecodeError(errorMsg)
}

// decodeErrorWithABI attempts to decode error using ABI definitions
func decodeErrorWithABI(errorMsg string) string {
	// Parse the ABI
	_, err := abi.JSON(strings.NewReader(errorABI))
	if err != nil {
		return "" // Fall back to manual decoding
	}
	
	// Extract the error data
	var selector string
	var arguments []byte
	
	// Check for custom error format with selector and data
	if strings.Contains(errorMsg, "custom error") {
		parts := strings.Split(errorMsg, "custom error ")
		if len(parts) > 1 {
			dataParts := strings.Split(parts[1], ": ")
			if len(dataParts) > 1 && len(dataParts[0]) >= 10 {
				selector = strings.TrimPrefix(dataParts[0], "0x")
				
				// Try to decode the hex data
				if hexData, err := hex.DecodeString(strings.Replace(dataParts[1], "0x", "", 1)); err == nil {
					arguments = hexData
				}
			}
		}
		
		// If we have a selector and arguments, try to match with ABI
		if selector != "" && len(arguments) > 0 {
			switch selector {
			case InsufficientBalanceSig:
				// Try to decode InsufficientBalance error
				var sender ethCommon.Address
				var balance, needed *big.Int
				
				// Manually decode arguments based on their types
				if len(arguments) >= 96 { // 3 parameters of 32 bytes each
					sender = ethCommon.BytesToAddress(arguments[12:32])
					balance = new(big.Int).SetBytes(arguments[32:64])
					needed = new(big.Int).SetBytes(arguments[64:96])
					
					return fmt.Sprintf("Insufficient Balance Error: address %s has balance %s but needs %s",
						sender.Hex(), balance.String(), needed.String())
				}
				
			case InsufficientAllowanceSig:
				// Try to decode InsufficientAllowance error
				if len(arguments) >= 128 { // 4 parameters of 32 bytes each
					owner := ethCommon.BytesToAddress(arguments[12:32])
					spender := ethCommon.BytesToAddress(arguments[44:64])
					allowance := new(big.Int).SetBytes(arguments[64:96])
					needed := new(big.Int).SetBytes(arguments[96:128])
					
					return fmt.Sprintf("Insufficient Allowance Error: spender %s is allowed %s by owner %s but needs %s",
						spender.Hex(), allowance.String(), owner.Hex(), needed.String())
				}
				
			case TransferFromZeroAddressSig:
				return "Error: Transfer from the zero address"
				
			case TransferToZeroAddressSig:
				return "Error: Transfer to the zero address"
				
			case ErrorSig:
				// Standard Error(string) revert
				if len(arguments) >= 96 {
					// String offset (first 32 bytes)
					// String length (next 32 bytes)
					stringLen := new(big.Int).SetBytes(arguments[32:64]).Int64()
					
					// String content starts at offset 64 and goes for stringLen bytes
					if len(arguments) >= int(64+stringLen) {
						return fmt.Sprintf("Error: %s", string(arguments[64:64+stringLen]))
					}
				}
			}
		}
	}
	
	// Standard "execution reverted: X" format
	if strings.Contains(errorMsg, "execution reverted: ") {
		parts := strings.Split(errorMsg, "execution reverted: ")
		if len(parts) > 1 && !strings.Contains(parts[1], "custom error") {
			return fmt.Sprintf("Error: %s", parts[1])
		}
	}
	
	return ""
}

// manualDecodeError is a fallback for when ABI decoding fails
// Uses string parsing and manual byte extraction to decode known error patterns
func manualDecodeError(errorMsg string) string {
	// Extract the custom error data
	customErrorData := ""
	if strings.Contains(errorMsg, "custom error") {
		parts := strings.Split(errorMsg, "custom error ")
		if len(parts) > 1 {
			dataParts := strings.Split(parts[1], ": ")
			if len(dataParts) > 1 {
				customErrorData = dataParts[0] + ": " + dataParts[1] // Combine error code and data
			} else {
				customErrorData = parts[1]
			}
		}
	}
	
	if customErrorData != "" {
		// Try to extract some meaning even without formal decoding
		if strings.HasPrefix(customErrorData, "0xe450d38c") {
			return "Insufficient Balance Error (manual decode)"
		} else if strings.HasPrefix(customErrorData, "0x4bd67a2d") {
			return "Insufficient Allowance Error (manual decode)"
		} else if strings.HasPrefix(customErrorData, "0xea553b34") {
			return "Transfer from the zero address"
		} else if strings.HasPrefix(customErrorData, "0xd92e233d") {
			return "Transfer to the zero address"
		}
		
		return fmt.Sprintf("Unknown error data: %s", customErrorData)
	}
	
	// Return the original message if all else fails
	return errorMsg
}
