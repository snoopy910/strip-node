package cardano

import (
	"math/big"
	"testing"

	"github.com/bnb-chain/tss-lib/v2/crypto"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/stretchr/testify/require"
)

func TestPublicKeyToAddress(t *testing.T) {
	rawKeyEddsa := eddsaKeygen.NewLocalPartySaveData(1)

	X, ok := new(big.Int).SetString("15112221349535400772501151409588531511454012693041857206046113283949847762202", 10)
	require.True(t, ok)

	Y, ok := new(big.Int).SetString("46316835694926478169428394003475163141307993866256225615783033603165251855960", 10)
	require.True(t, ok)

	ecPoint, err := crypto.NewECPoint(
		tss.Edwards(),
		X,
		Y,
	)
	require.NoError(t, err)
	rawKeyEddsa.EDDSAPub = ecPoint

	tests := []struct {
		name    string
		keyData *eddsaKeygen.LocalPartySaveData
		chainId string
		address string
		wantErr bool
	}{
		{
			name:    "Valid public key to mainnet address",
			keyData: &rawKeyEddsa,
			chainId: "1005",
			wantErr: false,
			address: "addr1vx8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49cgg73d3",
		},
		{
			name:    "Valid public key to testnet address",
			keyData: &rawKeyEddsa,
			chainId: "1006",
			wantErr: false,
			address: "addr_test1vz8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49cnq2dz5",
		},
		{
			name:    "Nil public key",
			keyData: nil,
			chainId: "1005",
			wantErr: true,
			address: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, gotTestnetAddress, err := PublicKeyToAddress(tt.keyData)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			t.Logf("Generated address: %s", gotAddress)
			require.True(t, IsValidAddress(gotAddress))

			// Check that testnet addresses start with "addr_test" and mainnet with "addr"
			if tt.chainId == "1006" || tt.chainId == "0" || tt.chainId == "1097" {
				require.Contains(t, gotTestnetAddress, "addr_test")
			} else {
				require.Contains(t, gotAddress, "addr")
				require.NotContains(t, gotAddress, "addr_test")
			}
			require.Equal(t, tt.address, gotAddress)
		})
	}
}

func TestAddressBech32Decoding(t *testing.T) {
	validAddresses := []string{
		"addr1vx8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49lj64t7r0vduj",
		"addr1q9ufw2v8p4mec986eqvjdhw6myn95kkfleqmf0ts5qrks9m5fgs5u9ke7pqyzmya8h4zye8qwgezehyvau59z93u4mlsnc2yy9",
		"addr_test1vx8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49cz766vy7q908m",
		"addr_test1qz2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzer3jcu5d8ps7zex2k2xt3uqxgjqnnj83ws8lhrn648jjxtwq2ytjqp",
	}

	for _, addr := range validAddresses {
		t.Run(addr, func(t *testing.T) {
			hrp, decoded, err := bech32.DecodeNoLimit(addr)
			if err != nil {
				t.Logf("Failed to decode address %s with error: %v", addr, err)
			} else {
				t.Logf("Successfully decoded address %s with HRP %s and %d decoded bytes", addr, hrp, len(decoded))
			}

			// Also check our IsValidAddress function
			isValid := IsValidAddress(addr)
			t.Logf("IsValidAddress result for %s: %v", addr, isValid)
		})
	}
}

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		// Valid mainnet addresses
		{
			name:    "Valid mainnet address",
			address: "addr1vx8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49lj64t7r0vduj",
			want:    true,
		},
		{
			name:    "Valid mainnet address 2",
			address: "addr1q9ufw2v8p4mec986eqvjdhw6myn95kkfleqmf0ts5qrks9m5fgs5u9ke7pqyzmya8h4zye8qwgezehyvau59z93u4mlsnc2yy9",
			want:    true,
		},
		// Valid testnet addresses
		{
			name:    "Valid testnet address",
			address: "addr_test1vx8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49cz766vy7q908m",
			want:    true,
		},
		{
			name:    "Valid testnet address 2",
			address: "addr_test1qz2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzer3jcu5d8ps7zex2k2xt3uqxgjqnnj83ws8lhrn648jjxtwq2ytjqp",
			want:    true,
		},
		// Invalid addresses
		{
			name:    "Empty address",
			address: "",
			want:    false,
		},
		{
			name:    "Wrong prefix",
			address: "xddr1vz8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49lj64t7x7eruvh",
			want:    false,
		},
		{
			name:    "Invalid characters",
			address: "addr1vz8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49lj64t7x$eruvh",
			want:    false,
		},
		{
			name:    "Too short",
			address: "addr1vz8huzmq",
			want:    false,
		},
		{
			name:    "Modified checksum",
			address: "addr1vz8huzmqryfxf65e8f6mv6q87ce2thfsgp20pg96ea3x49lj64t7x7eruve",
			want:    false,
		},
		{
			name:    "Bitcoin address",
			address: "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2",
			want:    false,
		},
		{
			name:    "Ethereum address",
			address: "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidAddress(tt.address)
			if got != tt.want {
				t.Errorf("IsValidAddress() = %v, want %v for address %s", got, tt.want, tt.address)
			}
		})
	}
}
