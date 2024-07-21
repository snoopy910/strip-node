package sequencer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func createWallet(identity string, identityCurve string) error {
	// select a list of nodes.
	// If length of selected nodes is more than maximum nodes then use maximum nodes length as signers.
	// If length of selected nodes is less than maximum nodes then use all nodes as signers.

	signers := SignersList()

	if len(signers) > MaximumSigners {
		// select random number of max signers
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(signers), func(i, j int) { signers[i], signers[j] = signers[j], signers[i] })
		signers = signers[:MaximumSigners]
	}

	signersPublicKeyList := make([]string, len(signers))
	for i, signer := range signers {
		signersPublicKeyList[i] = signer.PublicKey
	}

	// now create the wallet here
	createWalletRequest := CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "eddsa",
		Signers:       signersPublicKeyList,
	}

	marshalled, err := json.Marshal(createWalletRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", signers[0].URL+"/keygen", bytes.NewReader(marshalled))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 3 * time.Minute}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	createWalletRequest = CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "ecdsa",
		Signers:       signersPublicKeyList,
	}

	marshalled, err = json.Marshal(createWalletRequest)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("GET", signers[0].URL+"/keygen", bytes.NewReader(marshalled))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client = http.Client{Timeout: 3 * time.Minute}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	resp, err := http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=eddsa")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var getAddressResponse GetAddressResponse
	err = json.Unmarshal(body, &getAddressResponse)
	if err != nil {
		return err
	}

	eddsaAddress := getAddressResponse.Address

	resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=ecdsa")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err

	}

	err = json.Unmarshal(body, &getAddressResponse)
	if err != nil {
		return err
	}

	ecdsaAddress := getAddressResponse.Address

	wallet := WalletSchema{
		Identity:       identity,
		IdentityCurve:  identityCurve,
		Signers:        strings.Join(signersPublicKeyList, ","),
		EDDSAPublicKey: eddsaAddress,
		ECDSAPublicKey: ecdsaAddress,
	}

	_, err = AddWallet(&wallet)
	if err != nil {
		return err
	}

	return nil
}
