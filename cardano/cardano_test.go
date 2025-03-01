package cardano

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/StripChain/strip-node/util"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
)

var (
	ErrInvalidSignature = fmt.Errorf("error submitting transaction: Invalid signature")
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

func Test1(t *testing.T) {
	eddsaKey := "2rQXJBwJ3vCYreCHgcyhyejEXBdrJ5fct9kJsAo2ogr1"
	pubKey, err := base58.Decode(eddsaKey)
	if err != nil {
		t.Fatalf("failed to decode public key: %v", err)
	}
	fmt.Printf("pubKey: %+v\n", hex.EncodeToString(pubKey))

	require.Equal(t, "7ffafceddd9ed4b1bdce621cbd3e3d9df3f782b5670239c4cf15569958aa44bb", hex.EncodeToString(pubKey))
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
		{
			name:          "Valid ADA Transaction",
			serializedTxn: "84a300d90102818258200db6e71dccd30746ef028a694c716f118f2f7e1ae9c2c6fd51876e628c35ee11000182825839003eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335023eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335021a000f4240825839003eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335023eb4cb6ca6f6a07a316949ce3601d23b825da66904f69114ae0335021b0000000253f99480021a00030d40a0f5f6",
			chainId:       "1006",
			publicKey:     "7ffafceddd9ed4b1bdce621cbd3e3d9df3f782b5670239c4cf15569958aa44bb",
			signatureHex:  "5bffd637f8174c3d4b3197711db28d41a610ae967f5c52dda3b340cb88b176db146b61c28079164cc7d3a64f2f87a120333b02f37f2de7578b23900c2f67c009",
			ExpectError:   false,
		},
		// {
		// 	name:          "Invalid Signature",
		// 	serializedTxn: "83a40081825820c26a40e2a30305288cc8e85950bce2be4dc47eb58c768583d03cdd3177bd0e3b000181825839018d98bea0414243dc84070f96265577e7e6cf702d62e871016885034ecc64bf258b8e330cf0cdd9fdb03e10b4e4ac08f5da1fdec6222a34681a002dc6c0021a0002a42f031a032dcd55a100818258206d8a0b425bd2ec9692af39b1c0cf0e51caa07a603550e22f54091e872c7df29058407cf261f45babf645e6a5e61cd6df6ba7d1694a0dd1e7da9f405f16d1e7ae2d3c4c4d4ea2b8f3e7b0d2cd18c4b7be54c2c3e9d3a7d5c6b2d0f2d40e9f0e",
		// 	chainId:       "1005",
		// 	publicKey:     "6d8a0b425bd2ec9692af39b1c0cf0e51caa07a603550e22f54091e872c7df290",
		// 	signatureHex:  "invalid_signature",
		// 	ExpectError:   true,
		// 	ErrorMessage:  ErrInvalidSignature,
		// },
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
