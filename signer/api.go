package signer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/mr-tron/base58"
)

var messageChan = make(map[string]chan (Message))

func generateKeygenMessage(networkId string) {
	message := Message{
		Type:      MESSAGE_TYPE_GENERATE_START_KEYGEN,
		NetworkId: networkId,
	}

	broadcast(message)
}

func generateSignatureMessage(networkId string, msg []byte) {
	message := Message{
		Type:      MESSAGE_TYPE_START_SIGN,
		Hash:      msg,
		NetworkId: networkId,
	}

	broadcast(message)
}

func startHTTPServer(port string) {
	http.HandleFunc("/keygen", func(w http.ResponseWriter, r *http.Request) {
		networkId := r.URL.Query().Get("networkId")
		go generateKeygenMessage(networkId)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		networkId := r.URL.Query().Get("networkId")

		if networks[networkId].Key == nil {
			return
		}

		pk := edwards.PublicKey{
			Curve: tss.Edwards(),
			X:     networks[networkId].Key.EDDSAPub.X(),
			Y:     networks[networkId].Key.EDDSAPub.Y(),
		}

		publicKeyStr := base58.Encode(pk.Serialize())

		fmt.Fprintf(w, "%s", publicKeyStr)
	})

	http.HandleFunc("/signature", func(w http.ResponseWriter, r *http.Request) {
		c := rpc.New("https://api.devnet.solana.com")
		// pubKeyByte, _ := base58.Decode("4dEgqPG9FtjCiwW1HUdReraBozwq6qUCcDFXD8BnUn9Z")
		accountFrom := solana.MustPublicKeyFromBase58("4dEgqPG9FtjCiwW1HUdReraBozwq6qUCcDFXD8BnUn9Z")
		accountTo := solana.MustPublicKeyFromBase58("5oNDL3swdJJF1g9DzJiZ4ynHXgszjAEpUkxVYejchzrY")
		amount := uint64(1)

		// _, err := c.RequestAirdrop(context.Background(), accountFrom, solana.LAMPORTS_PER_SOL, rpc.CommitmentConfirmed)

		recentHash, err := c.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					amount,
					accountFrom,
					accountTo,
				).Build(),
			},
			solana.MustHashFromBase58(recentHash.Value.Blockhash.String()),
			solana.TransactionPayer(accountFrom),
		)

		if err != nil {
			panic(err)
		}

		msg, err := tx.Message.MarshalBinary()
		if err != nil {
			panic(err)
		}

		hash := string(msg)
		networkId := r.URL.Query().Get("networkId")
		go generateSignatureMessage(networkId, msg)

		messageChan[hash] = make(chan Message)

		sig := <-messageChan[hash]

		fmt.Println(len(sig.Message))
		newSig, err := edwards.ParseSignature(sig.Message)
		if err != nil {
			println("new sig error, ", err.Error())
		}

		fmt.Println(len(newSig.R.Bytes()), len(newSig.S.Bytes()))
		signingSigBytes := make([]byte, 88)
		copy(signingSigBytes[:32], newSig.R.Bytes())
		copy(signingSigBytes[32:], newSig.S.Bytes())

		signature := solana.SignatureFromBytes(sig.Message)

		tx.Signatures = append(tx.Signatures, signature)
		err = tx.VerifySignatures()

		if err != nil {
			fmt.Println(err)
		}

		wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
		if err != nil {
			panic(err)
		}

		rpcClient := rpc.New(rpc.DevNet_RPC)

		sig1, err := confirm.SendAndConfirmTransaction(
			context.Background(),
			rpcClient,
			wsClient,
			tx,
		)
		if err != nil {
			panic(err)
		}

		fmt.Println(sig1)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode("{\"signature\":\"" + string(sig.Message) + "\",\"address\":\"" + sig.Address + "\"}")
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		}

		delete(messageChan, hash)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
