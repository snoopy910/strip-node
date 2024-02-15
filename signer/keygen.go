package signer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	tssCommon "github.com/Silent-Protocol/go-sio/common"
	"github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/mr-tron/base58"
)

func updateKeygen(networkId string, partyKey string, from int, bz []byte, isBroadcast bool, to int) {
	instance := getSignerHubContract(
		networks[networkId].RPC_URL,
		networks[networkId].SignerHubContractAddress,
	)

	startKey, err := instance.StartKey(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	contractStartKey := startKey.Int64()
	contractTotalSigners, _ := instance.NextIndex(&bind.CallOpts{})

	parties, _ := getParties(int(contractTotalSigners.Int64()), int(contractStartKey))

	//wait for 1 minute for party to be ready
	for i := 0; i < 6; i++ {
		if !partyProcesses[networkId][partyKey].Exists {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	if !partyProcesses[networkId][partyKey].Exists {
		return
	}

	party := *partyProcesses[networkId][partyKey].Party

	if to != -1 && to != party.PartyID().Index {
		return
	}

	if party.PartyID().Index == from {
		return
	}

	pMsg, err := tss.ParseWireMessage(bz, parties[from], isBroadcast)
	if err != nil {
		panic(err)
	}

	ok, err := party.Update(pMsg)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("processed keygen message with status: ", ok)
}

func generateKeygen(networkId string) {
	instance := getSignerHubContract(
		networks[networkId].RPC_URL,
		networks[networkId].SignerHubContractAddress,
	)

	contractStartKey, err := instance.StartKey(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
		return
	}
	contractTotalSigners, err := instance.NextIndex(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
		return
	}
	contractThreshold, err := instance.CurrentThreshold(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
		return
	}
	contractIndex, err := instance.Signers(&bind.CallOpts{}, tssCommon.PublicKeyStrToBytes32(nodeKey.PublicKey))
	if err != nil {
		fmt.Println(err)
		return
	}

	if networks[networkId].StartKeyInt == 0 || contractStartKey.Int64() != int64(networks[networkId].StartKeyInt) {
		fmt.Println("started keygen generation")
		delete(partyProcesses[networkId], "keygen")
		parties, partiesIds := getParties(int(contractTotalSigners.Int64()), int(contractStartKey.Int64()))

		ctx := tss.NewPeerContext(parties)

		outChanKeygen := make(chan tss.Message)
		saveChan := make(chan *keygen.LocalPartySaveData)

		params := tss.NewParameters(tss.Edwards(), ctx, partiesIds[contractIndex.Index.Int64()], len(parties), int(contractThreshold.Int64()))

		localParty := keygen.NewLocalParty(params, outChanKeygen, saveChan)
		partyProcesses[networkId]["keygen"] = PartyProcess{&localParty, true}
		go localParty.Start()

		completed := false
		for !completed {
			select {
			case msg := <-outChanKeygen:
				dest := msg.GetTo()
				bytes, _, _ := msg.WireBytes()
				to := 0
				if dest == nil {
					to = -1
				} else {
					to = dest[0].Index
				}

				message := Message{
					Type:        MESSAGE_TYPE_GENERATE_KEYGEN,
					From:        msg.GetFrom().Index,
					To:          to,
					Message:     bytes,
					IsBroadcast: msg.IsBroadcast(),
					NetworkId:   networkId,
				}

				go broadcast(message)

			case save := <-saveChan:

				_i := save.EDDSAPub.IsOnCurve()

				fmt.Println("is on curve: ", _i)
				fmt.Println("saving key")

				pk := edwards.PublicKey{
					Curve: save.EDDSAPub.Curve(),
					X:     save.EDDSAPub.X(),
					Y:     save.EDDSAPub.Y(),
				}

				publicKeyStr := base58.Encode(pk.Serialize())

				fmt.Println("new TSS Address is: ", publicKeyStr)

				networkData := NetworkData{
					StartKey:     uint(contractStartKey.Int64()),
					TotalSigners: uint(contractTotalSigners.Int64()),
					Threshold:    uint(contractThreshold.Int64()),
					Index:        uint(contractIndex.Index.Int64()),
					Address:      publicKeyStr,
				}

				networkDataJSON, err := json.Marshal(networkData)
				if err != nil {
					fmt.Println(err)
				}

				err = ioutil.WriteFile(networks[networkId].NetworkDataFilePath, networkDataJSON, 0644)

				if err != nil {
					fmt.Println(err)
				}

				_network := networks[networkId]

				_network.Key = save
				out, err := json.Marshal(_network.Key)
				if err != nil {
					fmt.Println(err)
				}

				err = ioutil.WriteFile(_network.KeyFilePath, out, 0644)
				if err != nil {
					fmt.Println(err)
				}

				_network.StartKeyInt = int(contractStartKey.Int64())
				_network.Index = int(contractIndex.Index.Int64())
				_network.TotalSigners = int(contractTotalSigners.Int64())
				_network.Threshold = int(contractThreshold.Int64())
				networks[networkId] = _network
				completed = true

				fmt.Println("completed saving of new keygen ", publicKeyStr)
			}
		}
	} else {
		fmt.Println("cannot start keygen generation", networks[networkId].StartKeyInt, contractStartKey.Int64())
	}
}
