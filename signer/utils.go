package signer

import (
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/insight-chain/inb-go/crypto/sha3"
)

func getParties(totalSigners int) (tss.SortedPartyIDs, []*tss.PartyID) {
	partiesIds := []*tss.PartyID{}

	for i := 0; i < totalSigners; i++ {
		party := tss.NewPartyID(strconv.Itoa(i+1), "", big.NewInt(int64(i+1)))
		partiesIds = append(partiesIds, party)
	}

	parties := tss.SortPartyIDs(partiesIds)

	return parties, partiesIds
}

func publicKeyToAddress(pubkey []byte) string {
	var buf []byte
	_hash := sha3.NewKeccak256()
	_hash.Write(pubkey[1:])
	buf = _hash.Sum(nil)
	publicAddress := hexutil.Encode(buf[12:])
	return publicAddress
}

func publicKeyToBitcoinAddresses(pubkey []byte) (string, string, string) {
	log.Println("pubkey", pubkey) // NOTE: don't remove this log
	mainnetPubkey, err := btcutil.NewAddressPubKey(pubkey, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", ""
	}
	fmt.Println("mainnetPubkey: ", mainnetPubkey)

	testnetPubkey, err := btcutil.NewAddressPubKey(pubkey, &chaincfg.TestNet3Params)
	if err != nil {
		return mainnetPubkey.EncodeAddress(), "", ""
	}
	fmt.Println("testnetPubkey: ", testnetPubkey)

	regtestPubkey, err := btcutil.NewAddressPubKey(pubkey, &chaincfg.RegressionNetParams)
	if err != nil {
		return mainnetPubkey.EncodeAddress(), testnetPubkey.EncodeAddress(), ""
	}
	fmt.Println("regtestPubkey: ", regtestPubkey)

	return mainnetPubkey.EncodeAddress(), testnetPubkey.EncodeAddress(), regtestPubkey.EncodeAddress()
}

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %x or upper case
}

func SliceContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceIndexOfString(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func CalculateThreshold(totalSigners int) uint {
	if totalSigners == 1 || totalSigners == 2 {
		return 1
	} else {
		return uint((totalSigners / 2) + 1)
	}
}
