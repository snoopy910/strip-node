package aptos

import (
	"fmt"
	"testing"

	"github.com/StripChain/strip-node/util"
	"github.com/stretchr/testify/require"
)

var (
	ErrOldNonce         = fmt.Errorf("error submitting transaction: vm_error: Invalid transaction: Type: Validation Code: SEQUENCE_NUMBER_TOO_OLD")
	ErrInvalidSignature = fmt.Errorf("error submitting transaction: non-optional enum value is nil")
	ErrTxExpired        = fmt.Errorf("error submitting transaction: vm_error: Invalid transaction: Type: Validation Code: TRANSACTION_EXPIRED")
	ErrChainNotFound    = fmt.Errorf("error getting chain: chain not found")
	ErrInvalidTxHash    = fmt.Errorf("error getting transaction by hash: transaction_not_found: Transaction not found by Transaction hash(0x11af2d5770fe893b272fae56a854fafb074efe39bd5ab9c3f7a79a3995ff8b53)")
)

func TestGetAptosTransfers(t *testing.T) {
	// Test cases
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
	}{
		{
			Name:         "Native APT Transfer on Mainnet",
			ChainID:      "11",
			TxHash:       "0x0edbfc90e9a5c2d7b2932335088297f34ea42987113b75c752a26cc8978c2869",
			From:         "0x834d639b10d20dcb894728aa4b9b572b2ea2d97073b10eacb111f338b20ea5d7",
			To:           "0xf6dc6b07e7c49f240d197d1880f2806b7dfa3d838762fb3b9391f707bc6afbf7",
			Amount:       "799.99900000",
			Token:        "APT",
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: "79999900000",
		},
		{
			Name:         "Aptos Token Transfer on Mainnet",
			ChainID:      "11",
			TxHash:       "0x2a90d9b7f03525120bf2404c332d6531799dd68004ec208b15d1eb6435665176",
			From:         "0xe6f6ecc5eb1d93963c9d9446e5f33d677facd01511cf8bab18290c67641ef083",
			To:           "0x2621fd2b392dc53115dc3edea701ff6e4fe29c56a19388b9f9c455ecb65b1e4c",
			Amount:       "13536.630000",
			Token:        "TOMA",
			IsNative:     false,
			TokenAddress: "0x9d0595765a31f8d56e1d2aafc4d6c76f283c67a074ef8812d8c31bd8252ac2c3",
			ScaledAmount: "13536630000",
		},
		{
			Name:         "Fungible Asset Transfer on Mainnet",
			ChainID:      "11",
			TxHash:       "0xa10016b53c6b37d65e20210e04bfb9f79e09dc7defd74e063bb5f7c436e581f4",
			From:         "0x7c6018b4aa445a63f91954fe36ede071a6289ae5b35cfd117755c2a952713128",
			To:           "0x541e28fb12aa661a30358f2bebcd44460187ec918cb9cee075c2db86ee6aed93",
			Amount:       "1.00000000",
			Token:        "TVS",
			IsNative:     false,
			TokenAddress: "0x43782fca70e1416fc0c75954942dadd4af8d305a608b6153397ad5801b71e72d",
			ScaledAmount: "100000000",
		},
		{
			Name:         "Native APT Transfer",
			ChainID:      "167",
			TxHash:       "0x7f6fb178b1d6f2f57ac54bb39b3f778c1ec582ee7afa6004f82e1b6602f7c585",
			From:         "0xe93dcd3dd5febf8d72bf8d33e1d85a6115300fedfe055c062834d264f103ce4c",
			To:           "0xa2312d7ba5328b25f9fa7b9eed57911509946f90443283522c4e43ec30b35dcb",
			Amount:       "0.00000001",
			Token:        "APT",
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: "1",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			transfers, err := GetAptosTransfers(tt.ChainID, tt.TxHash)
			if err != nil {
				t.Errorf("getAptosTransfers() error = %v", err)
			}

			if len(transfers) == 0 {
				t.Fatalf("No transfer transaction for txHash, %v", tt.TxHash)
			}

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

func TestGetAptosTransfersError(t *testing.T) {
	tests := []struct {
		Name        string
		ChainID     string
		TxHash      string
		ExpectError bool
	}{
		{
			Name:        "Invalid ChainID",
			ChainID:     "123",
			TxHash:      "0x0edbfc90e9a5c2d7b2932335088297f34ea42987113b75c752a26cc8978c2869",
			ExpectError: true,
		},
		{
			Name:        "Empty TxHash",
			ChainID:     "11",
			TxHash:      "",
			ExpectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := GetAptosTransfers(tt.ChainID, tt.TxHash)
			if err == nil {
				t.Fatalf("getAptosTransfers() expected error %v, but got nil", tt.Name)
			}

			require.Error(t, err)
		})
	}
}

func TestSendAptosTransaction(t *testing.T) {
	// Test cases
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
			name:          "Validated APT Transaction signature",
			serializedTxn: "752b90c378361f841e10a7c5465e4328566f6f71a95377006f000df870ad835b00000000000000000200000000000000000000000000000000000000000000000000000000000000010d6170746f735f6163636f756e74087472616e736665720002209e8738089f6e847685f7b218f36f568f02a1594f29929e54140362576614b73e086400000000000000400d0300000000006400000000000000ea379b6700000000ab00",
			chainId:       "167",
			publicKey:     "0x7ca662d62370cf09ef3c739e651a597c43c1cf7c5f7560bde1fb69f7c45ff52b",
			signatureHex:  "8f18b8072a5390454e9ea0944315439313deeaac6b3f45f1fa03a5c2e41a6850188743f67055ba29ad989ccdc29d3faf3d58f905879c226c538d7dfd52b87408",
			ExpectError:   false,
		},
		{
			name:          "Validated APT Transaction signature",
			serializedTxn: "752b90c378361f841e10a7c5465e4328566f6f71a95377006f000df870ad835b01000000000000000200000000000000000000000000000000000000000000000000000000000000010d6170746f735f6163636f756e74087472616e736665720002209e8738089f6e847685f7b218f36f568f02a1594f29929e54140362576614b73e086400000000000000400d0300000000006400000000000000cd399b6700000000ab00",
			chainId:       "167",
			publicKey:     "0x7ca662d62370cf09ef3c739e651a597c43c1cf7c5f7560bde1fb69f7c45ff52b",
			signatureHex:  "fa85d4c414bf3a2801cd97c982dea83fb6b8647fac565dd4435b23c7702b4ba315c7fe9c3e33c7b997065ab086896314d01306adccf1537b710bc328b15fec09",
			ExpectError:   false,
		},
		{
			name:          "Invalidated APT Transaction signature",
			serializedTxn: "752b90c378361f841e10a7c5465e4328566f6f71a95377006f000df870ad835b01000000000000000200000000000000000000000000000000000000000000000000000000000000010d6170746f735f6163636f756e74087472616e736665720002209e8738089f6e847685f7b218f36f568f02a1594f29929e54140362576614b73e086400000000000000400d0300000000006400000000000000cd399b6700000000ab00",
			chainId:       "167",
			publicKey:     "0x7ca662d62370cf09ef3c739e651a597c43c1cf7c5f7560bde1fb69f7c45ff52b",
			signatureHex:  "aa85d4c414bf3a2801cd97c982dea83fb6b8647fac565dd4435b23c770234ba315c7fe9c3e33c7b997065ab086896314d01306adccf1537b710bc328b15fec09",
			ExpectError:   true,
			ErrorMessage:  ErrInvalidSignature,
		},
		{
			name:          "Expired Transaction",
			serializedTxn: "1ee0894555a9e029fb411b729e293ffb84bf6dd4d54eee2c2167d50a4061a2cd00000000000000000200000000000000000000000000000000000000000000000000000000000000010d6170746f735f6163636f756e74087472616e736665720002209e8738089f6e847685f7b218f36f568f02a1594f29929e54140362576614b73e086400000000000000400d030000000000640000000000000075469b6700000000ab00",
			chainId:       "167",
			publicKey:     "0x2a6014ea3b423190c0040078b9659fa5ba6435d9f95e7e29982b7bedd0544570",
			signatureHex:  "46790e862b45ca53d24e8f0605b4baf92b3bf924954b3758781b1f2d4e4673f4286817522d8970d03cacfd5678e0e5a7bfaaf56e38c95071824bdc5aa00b0800",
			ExpectError:   true,
			ErrorMessage:  ErrTxExpired,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SendAptosTransaction(tt.serializedTxn, tt.chainId, "ed25519", tt.publicKey, tt.signatureHex)
			if !tt.ExpectError {
				if err != nil {
					require.Error(t, err)
					// Expect nonce error and it shows that the signature is valid for the public key and transaction contents
					require.Equal(t, err, ErrOldNonce)
				} else {
					t.Fatalf("Expected nonce error, %v", tt.name)
				}
			} else {
				if err != nil {
					require.Error(t, err)
					// Expect nil enum value error since failed to set authenticator because of the invalidate signature
					require.Equal(t, err, tt.ErrorMessage)
				} else {
					t.Fatalf("Expected nonce error, %v", tt.name)
				}
			}
		})
	}
}

func TestCheckAptosTransactionConfirmed(t *testing.T) {
	// Test cases
	tests := []struct {
		Name         string
		ChainID      string
		TxHash       string
		Success      bool
		ExpectError  bool
		ErrorMessage error
	}{
		{
			Name:    "Succeed transaction on mainnet",
			ChainID: "11",
			TxHash:  "0x11af2d5770fe893b272fae56a854fafb074efe39bd5ab9cff7a79a3995ff8b53",
			Success: true,
		},
		{
			Name:    "Succeed transaction on devnet",
			ChainID: "167",
			TxHash:  "0x7f6fb178b1d6f2f57ac54bb39b3f778c1ec582ee7afa6004f82e1b6602f7c585",
			Success: true,
		},
		{
			Name:         "Invalid txHash on mainnet",
			ChainID:      "11",
			TxHash:       "0x11af2d5770fe893b272fae56a854fafb074efe39bd5ab9c3f7a79a3995ff8b53",
			ExpectError:  true,
			ErrorMessage: ErrInvalidTxHash,
		},
		{
			Name:         "Invalid chainId",
			ChainID:      "123",
			TxHash:       "0x11af2d5770fe893b272fae56a854fafb074efe39bd5ab9c3f7a79a3995ff8b53",
			ExpectError:  true,
			ErrorMessage: ErrChainNotFound,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			success, err := CheckAptosTransactionConfirmed(tt.ChainID, tt.TxHash)
			if tt.ExpectError {
				require.Error(t, err)
				require.Equal(t, tt.ErrorMessage, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, success, tt.Success)
		})
	}
}
