package signer

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Silent-Protocol/go-sio/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type NodeKey struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

var nodeKey NodeKey

func loadKey(path string, privateKey string, publicKey string) {
	if privateKey != "" && publicKey != "" {
		nodeKey.PrivateKey = privateKey
		nodeKey.PublicKey = publicKey

		return
	}

	filePath := path + "/" + "identity.json"
	exists, err := common.FileExists(filePath)
	if err != nil {
		panic(err)
	}

	if exists {
		jsonFile, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &nodeKey)

		fmt.Println("Loaded node key: ", nodeKey.PublicKey)
	} else {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		nodeKey.PrivateKey = hexutil.Encode(privateKeyBytes)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}

		publicKeyBytes := crypto.CompressPubkey(publicKeyECDSA)

		privateKeyStr := hexutil.Encode(privateKeyBytes)
		publicKeyStr := hexutil.Encode(publicKeyBytes)

		nodeKey.PublicKey = publicKeyStr
		nodeKey.PrivateKey = privateKeyStr

		nodeKeyJSON, err := json.Marshal(nodeKey)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(filePath, nodeKeyJSON, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Generated node key: ", publicKeyStr)
	}
}
