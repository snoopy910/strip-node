package libs

import (
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs/blockchains"
)

type CreateWalletRequest struct {
	Identity      string       `json:"identity"`
	IdentityCurve common.Curve `json:"identityCurve"`
	KeyCurve      common.Curve `json:"keyCurve"`
	Signers       []string     `json:"signers"`
}

type AddressesResponse struct {
	Addresses    map[blockchains.BlockchainID]map[blockchains.NetworkType]string `json:"addresses"`
	EDDSAAddress string                                                          `json:"eddsaAddress"`
	ECDSAAddress string                                                          `json:"ecdsaAddress"`
}
