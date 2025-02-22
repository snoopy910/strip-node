package ripple

import (
	"math/big"
	"testing"

	"github.com/bnb-chain/tss-lib/v2/crypto"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/stretchr/testify/require"
)

func TestIsValidRippleAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "Valid address",
			address: "rKmqeMFEDmqGdbrdaTu4oUWPM3rbQGmSdj",
			want:    true,
		},
		{
			name:    "Invalid prefix",
			address: "xBJHj7tQASabS2yb8XHE2by9APxZfQpFGF",
			want:    false,
		},
		{
			name:    "Too short",
			address: "rBJHj7",
			want:    false,
		},
		{
			name:    "Too long",
			address: "rBJHj7tQASabS2yb8XHE2by9APxZfQpFGFrBJHj7tQASabS2yb8XHE2by9APxZfQpFGF",
			want:    false,
		},
		{
			name:    "Invalid characters",
			address: "rBJHj7tQASabS2yb8XHE2by9APxZfQpFG0",
			want:    false,
		},
		{
			name:    "Empty string",
			address: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidRippleAddress(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

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
		name        string
		rawKeyEddsa *eddsaKeygen.LocalPartySaveData
		want        string
	}{
		{
			name:        "Valid public key to Ripple address",
			rawKeyEddsa: &rawKeyEddsa,
			want:        "rGGasCecEGuD39ag5S1cgKHdMxMyn6nfDh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPubKey := PublicKeyToAddress(tt.rawKeyEddsa)

			require.Equal(t, tt.want, gotPubKey)
			require.True(t, len(gotPubKey) == 66)  // 33 bytes in hex
			require.True(t, gotPubKey[:2] == "ED") // Should start with ED for Ed25519
			require.True(t, IsValidRippleAddress(gotPubKey))
		})
	}
}
