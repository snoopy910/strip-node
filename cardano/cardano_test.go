package cardano

import (
	"fmt"
	"testing"

	"github.com/StripChain/strip-node/util"
	"github.com/stretchr/testify/require"
)

var (
	ErrInvalidSignature = fmt.Errorf("InvalidWitnessesUTXOW")
	ErrTxExpired        = fmt.Errorf("error submitting transaction: Transaction expired")
)

func TestGetCardanoTransfers(t *testing.T) {
	tests := []struct {
		Name      string
		ChainID   string
		TxHash    string
		Transfers []struct {
			From         string
			To           string
			Amount       string
			Token        string
			IsNative     bool
			TokenAddress string
			ScaledAmount string
		}
		ExpectError bool
	}{
		{
			Name:    "Single ADA Transfer in Transaction",
			ChainID: "1005",
			TxHash:  "259a74a59974bea81b35ff040a642283b6cb25dcb8018dbe2f29297816636367",
			Transfers: []struct {
				From         string
				To           string
				Amount       string
				Token        string
				IsNative     bool
				TokenAddress string
				ScaledAmount string
			}{
				{
					From:         "addr1q9eymn0aemrxakvprrsp626jvyh2nwljlt8f9fdqmvjgxrxk6dh558s7vymgqvplwkx3x0ypj9fujmwyfzzw8z3ctxjsjln78r",
					To:           "addr1vyc2fkntdhuyvzusnfk6es6t8g2vcsykr89p38yyne65tcs4me4az",
					Amount:       "5.681244",
					Token:        "ADA",
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: "5681244",
				},
			},
		},
		{
			Name:    "Multiple Native Token Transfer",
			ChainID: "1005",
			TxHash:  "27adfdc628ca244b6df56370f1f74987c537bb9746b69d9b2bdd04b2fc2f9052",
			Transfers: []struct {
				From         string
				To           string
				Amount       string
				Token        string
				IsNative     bool
				TokenAddress string
				ScaledAmount string
			}{
				{
					From:         "addr1w8rjw3pawl0kelu4mj3c8x20fsczf5pl744s9mxz9v8n7eg0fcr8k",
					To:           "addr1qywp9uddymvveq4z2l4dxpvzzfy74muxm2mdmwcxgt46zjz0yyusk4030ytwh4a7xr3l62thn0qmu73ru0zxmg9d0yysjrs2fr",
					Amount:       "1.176630",
					Token:        "ADA",
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: "1176630",
				},
				{
					From:         "addr1w8rjw3pawl0kelu4mj3c8x20fsczf5pl744s9mxz9v8n7eg0fcr8k",
					To:           "addr1qywp9uddymvveq4z2l4dxpvzzfy74muxm2mdmwcxgt46zjz0yyusk4030ytwh4a7xr3l62thn0qmu73ru0zxmg9d0yysjrs2fr",
					Amount:       "1.000000",
					Token:        "667Gremlin425",
					IsNative:     false,
					TokenAddress: "ff6be1474bbfe726665a9f37314c8cac55b2aa88a5c3b90af102fb0c3636374772656d6c696e343235",
					ScaledAmount: "1",
				},
				{
					From:         "addr1w8rjw3pawl0kelu4mj3c8x20fsczf5pl744s9mxz9v8n7eg0fcr8k",
					To:           "addr1qywp9uddymvveq4z2l4dxpvzzfy74muxm2mdmwcxgt46zjz0yyusk4030ytwh4a7xr3l62thn0qmu73ru0zxmg9d0yysjrs2fr",
					Amount:       "112.682723",
					Token:        "ADA",
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: "112682723",
				},
				{
					From:         "addr1w8rjw3pawl0kelu4mj3c8x20fsczf5pl744s9mxz9v8n7eg0fcr8k",
					To:           "addr1qywp9uddymvveq4z2l4dxpvzzfy74muxm2mdmwcxgt46zjz0yyusk4030ytwh4a7xr3l62thn0qmu73ru0zxmg9d0yysjrs2fr",
					Amount:       "7.080861",
					Token:        "ADA",
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: "7080861",
				},
			},
		},
		{
			Name:        "Empty TxHash",
			ChainID:     "1005",
			TxHash:      "",
			ExpectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			transfers, err := GetCardanoTransfers(tt.ChainID, tt.TxHash)
			if tt.ExpectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, transfers)
			require.Equal(t, len(tt.Transfers), len(transfers), "number of transfers should match")

			for i, expectedTransfer := range tt.Transfers {
				require.Equal(t, expectedTransfer.From, transfers[i].From)
				require.Equal(t, expectedTransfer.To, transfers[i].To)
				require.Equal(t, expectedTransfer.Amount, transfers[i].Amount)
				require.Equal(t, expectedTransfer.Token, transfers[i].Token)
				require.Equal(t, expectedTransfer.IsNative, transfers[i].IsNative)
				require.Equal(t, expectedTransfer.TokenAddress, transfers[i].TokenAddress)
				require.Equal(t, expectedTransfer.ScaledAmount, transfers[i].ScaledAmount)
			}
		})
	}
}

func TestSendCardanoTransaction(t *testing.T) {
	tests := []struct {
		name          string
		serializedTxn string
		chainId       string
		publicKey     string
		signatureHex  string
		ExpectError   bool
		ErrorMessage  error
	}{
		// These values are valid transactions on the Cardano testnet but because of the UTXO they are only valid once
		// {
		// 	name:          "Valid ADA Transaction",
		// 	serializedTxn: "84a300d90102818258200db6e71dccd30746ef028a694c716f118f2f7e1ae9c2c6fd51876e628c35ee11000182825839003eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335023eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335021a000f4240825839003eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335023eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335021b0000000253f99480021a00030d40a0f5f6",
		// 	chainId:       "1006",
		// 	publicKey:     "7ffafceddd9ed4b1bdce621cbd3e3d9df3f782b5670239c4cf15569958aa44bb",
		// 	signatureHex:  "5bffd637f8174c3d4b3197711db28d41a610ae967f5c52dda3b340cb88b176db146b61c28079164cc7d3a64f2f87a120333b02f37f2de7578b23900c2f67c009",
		// 	ExpectError:   false,
		// },
		// {
		// 	name:          "Valid ADA Transaction",
		// 	serializedTxn: "84a300d901028282582027880d7a68d8543a3f78fab79e4a7bb19e650d8a60ec7f54700e84e35555faac0082582027880d7a68d8543a3f78fab79e4a7bb19e650d8a60ec7f54700e84e35555faac01018282583900ceb40e1b07552d02162706a9d2f75ac359330e894657fd67b83d71ceceb40e1b07552d02162706a9d2f75ac359330e894657fd67b83d71ce1a000f424082583900ceb40e1b07552d02162706a9d2f75ac359330e894657fd67b83d71ceceb40e1b07552d02162706a9d2f75ac359330e894657fd67b83d71ce1b0000000253f06cc0021a00030d40a0f5f6",
		// 	chainId:       "1006",
		// 	publicKey:     "dc0c836c405e1d3f421f9107f1f18da5c941223adbbf73d49d35de9c4d858bf8",
		// 	signatureHex:  "aa0326b03aa3af20e4ac01bf0e95420ae4f613ac70672741c49932354b9f0404353159ee2349756f5452791f61b1f1ed9bd04178637b322a2c80017239689e0e",
		// 	ExpectError:   false,
		// },

		{
			name:          "Invalid Signature",
			serializedTxn: "84a300d90102818258200db6e71dccd30746ef028a694c716f118f2f7e1ae9c2c6fd51876e628c35ee11000182825839003eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335023eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335021a000f4240825839003eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335023eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335021b0000000253f99480021a00030d40a0f5f6",
			chainId:       "1005",
			publicKey:     "6d8a0b425bd2ec9692af39b1c0cf0e51caa07a603550e22f54091e872c7df290",
			signatureHex:  "aa0326b03aa3af20e4ac01bf0e95420ae4f613ac70672741c49932354b9f0404353159ee2349756f5452791f61b1f1ed9bd04178637b322a2c80017239689e0e",
			ExpectError:   true,
			ErrorMessage:  ErrInvalidSignature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SendCardanoTransaction(tt.serializedTxn, tt.chainId, "", tt.publicKey, tt.signatureHex)
			if tt.ExpectError {
				require.Error(t, err)
				if tt.ErrorMessage != nil {
					require.Contains(t, err.Error(), tt.ErrorMessage.Error())
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckCardanoTransactionConfirmed(t *testing.T) {
	tests := []struct {
		Name         string
		ChainID      string
		TxHash       string
		Success      bool
		ExpectError  bool
		ErrorMessage error
	}{
		{
			Name:    "Successful transaction on mainnet",
			ChainID: "1005",
			TxHash:  "1142b4e3f62d2d252c621a9dbb4d9a293be0a3592d1102c1036083b3e4c839e2",
			Success: true,
		},
		{
			Name:    "Successful transaction on testnet",
			ChainID: "1006",
			TxHash:  "a98b462c928d51eb119b47b937cb3ca07ca5d9a54e6788e350fe76ee5e3809b8",
			Success: true,
		},
		{
			Name:        "Transaction not found on mainnet",
			ChainID:     "1005",
			TxHash:      "0000000000000000000000000000000000000000000000000000000000000000",
			ExpectError: false,
			Success:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			success, err := CheckCardanoTransactionConfirmed(tt.ChainID, tt.TxHash)
			if tt.ExpectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.ErrorMessage.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.Success, success)
		})
	}
}
