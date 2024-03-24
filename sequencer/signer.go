package sequencer

import "fmt"

type Signer struct {
	PublicKey string
	URL       string
}

// TODO: This list will be fetched from SC by the sequencer
var Signers = []Signer{
	{
		PublicKey: "0x0226d1556a83c01a9d2b1cce29b32cb520238efc602f86481d2d0b9af8a2fff0cf",
		URL:       "http://localhost:8080",
	},
	{
		PublicKey: "0x0354455a1f7f4244ef645ac62baa8bd90af0cc18cdb0eae369766b7b58134edf35",
		URL:       "http://localhost:8081",
	},
}

func GetSigner(publicKey string) (Signer, error) {
	for _, signer := range Signers {
		if signer.PublicKey == publicKey {
			return signer, nil
		}
	}
	return Signer{}, fmt.Errorf("signer not found")
}
