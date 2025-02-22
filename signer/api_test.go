// package signer

// import (
// 	"encoding/json"
// 	"testing"

// 	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
// 	"github.com/bnb-chain/tss-lib/v2/tss"
// 	"github.com/decred/dcrd/dcrec/edwards"
// 	"github.com/stellar/go/strkey"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/test-go/testify/require"
// )

// func TestAddressEndpoint(t *testing.T) {
// 	identity := "0xE9D46772E8441671718bf8c1521BB9F717EDB7Fd"
// 	identityCurve := "ecdsa"
// 	keyCurve := "stellar_eddsa"

// 	InitialiseDB("localhost:5432", "postgres", "postgres", "password")

// 	keyShare, err := GetKeyShare(identity, identityCurve, keyCurve)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, keyShare)
// 	address := StellarGenerateAddress(t, keyShare)
// 	assert.Equal(t, "GBDYDX3STKB6O36JFMOSGGMGQH7MWTKFRSDGIF4UYCAC6IC6NRC3WRVX", address)

// 	// Verify address checksum using Stellar SDK
// 	_, err = strkey.Decode(strkey.VersionByteAccountID, address)
// 	require.NoError(t, err, "Invalid Stellar address checksum")
// }

// func StellarGenerateAddress(t *testing.T, keyShare string) string {
// 	var rawKeyEddsa *eddsaKeygen.LocalPartySaveData
// 	json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

// 	pk := edwards.PublicKey{
// 		Curve: tss.Edwards(),
// 		X:     rawKeyEddsa.EDDSAPub.X(),
// 		Y:     rawKeyEddsa.EDDSAPub.Y(),
// 	}

// 	// Get the public key bytes
// 	pkBytes := pk.Serialize()

// 	// Stellar StrKey format:
// 	if len(pkBytes) != 32 {
// 		return ""
// 	}

// 	// Version byte for ED25519 public key in Stellar
// 	versionByte := strkey.VersionByteAccountID // 6 << 3, or 48

// 	// Use Stellar SDK's strkey package to encode
// 	address, err := strkey.Encode(versionByte, pkBytes)
// 	if err != nil {
// 		return ""
// 	}

//		return address
//	}
package signer
