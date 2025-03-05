package sequencer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// createWallet creates a new wallet with the specified identity and identity curve.
// Note: It selects a list of signers, ensuring the number of signers does not exceed MaximumSigners.
// The function performs the following steps:
// 1. Selects a random subset of signers if the total number exceeds MaximumSigners.
// 2. Constructs a CreateWalletRequest for both "eddsa" and "ecdsa" key curves.
// 3. Sends HTTP requests to the first signer's URL to generate keys.
// 4. Retrieves the generated addresses for both key curves.
// 5. Constructs a WalletSchema and adds it to the wallet store.
//
// Parameters:
// - identity: A string representing the identity for the wallet.
// - identityCurve: A string representing the curve type for the identity.
//
// Returns:
// - error: An error if any step in the wallet creation process fails.
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

	// create the wallet whose keycurve is eddsa here
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

	// create the wallet whose keycurve is sui_eddsa here
	createWalletRequest = CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "sui_eddsa",
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
		return fmt.Errorf("failed to create sui wallet: %v", err)
	}

	// create the wallet whose keycurve is aptos_eddsa here
	createWalletRequest = CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "aptos_eddsa",
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

	// create the wallet whose keycurve is ecdsa here
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

	// create the wallet whose keycurve is dogecoin_ecdsa here
	createWalletRequest = CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "dogecoin_ecdsa",
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

	// create the wallet whose keycurve is bitcoin_ecdsa here
	createWalletRequest = CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "bitcoin_ecdsa",
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

	// create the wallet whose keycurve is algorand_eddsa here
	createWalletRequest = CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "algorand_eddsa",
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

	// get the address of the wallet whose keycurve is eddsa here
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

	// get the address of the wallet whose keycurve is sui_eddsa here
	resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=sui_eddsa")
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

	suiAddress := getAddressResponse.Address

	// get the address of the wallet whose keycurve is aptos_eddsa here
	resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=aptos_eddsa")
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

	aptosEddsaAddress := getAddressResponse.Address

	// get the address of the wallet whose keycurve is ecdsa here
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

	// get the address of the wallet whose keycurve is bitcoin_ecdsa here
	resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=bitcoin_ecdsa")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var getBitcoinAddressesResponse GetBitcoinAddressesResponse

	err = json.Unmarshal(body, &getBitcoinAddressesResponse)
	if err != nil {
		return err
	}

	bitcoinMainnetAddress := getBitcoinAddressesResponse.MainnetAddress
	bitcoinTestnetAddress := getBitcoinAddressesResponse.TestnetAddress
	bitcoinRegtestAddress := getBitcoinAddressesResponse.RegtestAddress

	// get the address of the wallet whose keycurve is dogecoin_ecdsa here
	resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=dogecoin_ecdsa")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var getDogecoinAddressesResponse GetDogecoinAddressesResponse

	err = json.Unmarshal(body, &getDogecoinAddressesResponse)
	if err != nil {
		return err
	}

	dogecoinMainnetAddress := getDogecoinAddressesResponse.MainnetAddress
	dogecoinTestnetAddress := getDogecoinAddressesResponse.TestnetAddress

	// get the address of the wallet whose keycurve is algorand_eddsa here
	resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=algorand_eddsa")
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

	algorandEddsaAddress := getAddressResponse.Address

	// create the wallet whose keycurve is stellar_eddsa here
	createWalletRequest = CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "stellar_eddsa",
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

	// get the address of the wallet whose keycurve is stellar_eddsa here
	resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=stellar_eddsa")
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

	stellarAddress := getAddressResponse.Address

	// create the wallet whose keycurve is stellar_eddsa here
	createWalletRequest = CreateWalletRequest{
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      "ripple_eddsa",
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

	// get the address of the wallet whose keycurve is stellar_eddsa here
	resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=ripple_eddsa")
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

	rippleAddress := getAddressResponse.Address

	// add created wallet to the store
	wallet := WalletSchema{
		Identity:                 identity,
		IdentityCurve:            identityCurve,
		Signers:                  strings.Join(signersPublicKeyList, ","),
		EDDSAPublicKey:           eddsaAddress,
		AptosEDDSAPublicKey:      aptosEddsaAddress,
		ECDSAPublicKey:           ecdsaAddress,
		BitcoinMainnetPublicKey:  bitcoinMainnetAddress,
		BitcoinTestnetPublicKey:  bitcoinTestnetAddress,
		BitcoinRegtestPublicKey:  bitcoinRegtestAddress,
		DogecoinMainnetPublicKey: dogecoinMainnetAddress,
		DogecoinTestnetPublicKey: dogecoinTestnetAddress,
		SuiPublicKey:             suiAddress,
		StellarPublicKey:         stellarAddress,
		AlgorandEDDSAPublicKey:   algorandEddsaAddress,
		RippleEDDSAPublicKey:     rippleAddress,
	}

	_, err = AddWallet(&wallet)
	if err != nil {
		return err
	}

	return nil
}
