package main

import (
	"fmt"
	"math/big"
	"strconv"

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

// toHexInt returns a 64-character hexadecimal string of n, padded with leading zeros.
func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%064x", n)
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
