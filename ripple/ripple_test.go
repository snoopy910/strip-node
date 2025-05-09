package ripple

import (
	"fmt"
	"testing"

	"github.com/StripChain/strip-node/util"
	"github.com/stretchr/testify/require"
)

var (
	ErrOldSequence      = fmt.Errorf("error submitting transaction: tefPAST_SEQ")
	ErrInvalidSignature = fmt.Errorf("error submitting transaction: invalidTransaction 0  fails local checks: Invalid signature")
	ErrTxExpired        = fmt.Errorf("error submitting transaction: tefPAST_EXPIRATION")
	ErrChainNotFound    = fmt.Errorf("error getting chain: chain not found")
	ErrInvalidTxHash    = fmt.Errorf("failed to parse transaction hash")
)

func TestGetRippleTransfers(t *testing.T) {
	tests := []struct {
		Name         string
		ChainID      string
		TxHash       string
		From         string
		To           string
		Amount       string
		Token        string
		IsNative     bool
		TokenAddress string
		ScaledAmount string
		ExpectError  bool
	}{
		{
			Name:         "Native XRP Transfer on Mainnet",
			ChainID:      "1003",
			TxHash:       "64BEF9A1AED96163CC80C6F415B8F687E9F30FEF7156F16629B6AA3EF8D566B2",
			From:         "rNhEokJ6ZoE7BwYPvLGZ88oMDadtAXV4Zm",
			To:           "rNxp4h8apvRis6mJf9Sh8C6iRxfrDWN7AV",
			Amount:       "1",
			Token:        "XRP",
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: "1000000",
		},
		{
			Name:         "USD Token Transfer on Mainnet",
			ChainID:      "1003",
			TxHash:       "935D7281CD03C5A1C95A25D69BED06F71E761EB463DF8F972FB0B706D6759F7D",
			From:         "rogue5HnPRSszD9CWGSUz8UGHMVwSSKF6",
			To:           "rogue5HnPRSszD9CWGSUz8UGHMVwSSKF6",
			Amount:       "0.000001107386733739162",
			Token:        "BITx",
			IsNative:     false,
			TokenAddress: "rBitcoiNXev8VoVxV7pwoQx1sSfonVP9i3",
			ScaledAmount: "1",
		},
		{
			Name:        "Empty TxHash",
			ChainID:     "1003",
			TxHash:      "",
			ExpectError: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			transfers, err := GetRippleTransfers(tt.ChainID, tt.TxHash)
			if tt.ExpectError {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			require.NotEmpty(t, transfers)

			require.Equal(t, tt.From, transfers[0].From)
			require.Equal(t, tt.To, transfers[0].To)
			require.Equal(t, tt.Amount, transfers[0].Amount)
			require.Equal(t, tt.Token, transfers[0].Token)
			require.Equal(t, tt.IsNative, transfers[0].IsNative)
			require.Equal(t, tt.TokenAddress, transfers[0].TokenAddress)
			require.Equal(t, tt.ScaledAmount, transfers[0].ScaledAmount)
		})
	}
}

func TestSendRippleTransaction(t *testing.T) {
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
			name:          "Valid XRP Transaction",
			serializedTxn: "7b225472616e73616374696f6e54797065223a225061796d656e74222c224163636f756e74223a22724a35726a4e685247504576774c77737539476a7339557472466f51636a67793857222c2244657374696e6174696f6e223a2272665648554e6a464c5a543778466b596a6e424a723547765038745a6254745a6633222c22416d6f756e74223a2231303030303030222c22466c616773223a302c2253657175656e6365223a353034343637322c22466565223a223132222c224c6173744c656467657253657175656e6365223a353034343639322c225369676e696e675075624b6579223a22454446324344423744393930373431344346334332463937333044373641333044354332344432464335334532444443443136433442464143414637343232313944227d",
			chainId:       "1004",
			publicKey:     "EDF2CDB7D9907414CF3C2F9730D76A30D5C24D2FC53E2DDCD16C4BFACAF742219D",
			signatureHex:  "ce30752b22bea759807e015394dd03c309c9042684d7b1666f6e43d959801e6fa96af524089744f41dcd8799b168ddecc96e8f523ae936f37aa774251ce27506",
			ExpectError:   false,
		},
		{
			name:          "Invalid Signature",
			serializedTxn: "7b225472616e73616374696f6e54797065223a225061796d656e74222c224163636f756e74223a2272396977466f506f6a395a6d6e5931426641714c54616a356969385632473542556b222c2244657374696e6174696f6e223a2272513359386277547564567476576a6d617447387a524c34393132337a3779684d54222c224e6574776f726b4944223a312c22416d6f756e74223a2231303030303030222c22466c616773223a302c2253657175656e6365223a353033303337352c22466565223a223132222c224c6173744c656467657253657175656e6365223a353033303435377d",
			chainId:       "1003",
			publicKey:     "acb6845b75c6f6c045a785c7f5165d76753d14207166fcaf5207565b60",
			signatureHex:  "240af8077c8ae5ca1c2398b45525f4e2c4d769d6726b1e67a76c65b3abe4e12860138eb5bf0150a36abd203003412aa6e059de10fa01b72e156b9215f1f87a0d",
			ExpectError:   true,
			ErrorMessage:  ErrInvalidSignature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SendRippleTransaction(tt.serializedTxn, tt.chainId, "", tt.publicKey, tt.signatureHex)
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

func TestCheckRippleTransactionConfirmed(t *testing.T) {
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
			ChainID: "1003",
			TxHash:  "B4F812F23479CE3744B6EAD7F99C3106E064E81E26BB2FC9905EECA5FAE41145",
			Success: true,
		},
		{
			Name:    "Successful transaction on testnet",
			ChainID: "1004",
			TxHash:  "6ADEA1249CB77B94F68624B87382818EC3804C648DACFE7788B189575CE7846C",
			Success: true,
		},
		{
			Name:        "Transaction not found on mainnet",
			ChainID:     "1003",
			TxHash:      "655B2D3EF03E530E0F64D3C38EE2C713AB7F7705AFFC2064C2B6754B1D756AE6",
			ExpectError: false,
			Success:     false,
		},
		{
			Name:         "Invalid txHash",
			ChainID:      "1003",
			TxHash:       "INVALID_HASH",
			ExpectError:  true,
			ErrorMessage: ErrInvalidTxHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			success, err := CheckRippleTransactionConfirmed(tt.ChainID, tt.TxHash)
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
