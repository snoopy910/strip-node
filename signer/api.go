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

		//0100010335db7213e45b7498a259b23793399c2821d56df35856f436d87c932ae396034c474f7335d5399e496566fcfe89b06ddf2f9df31fab601aafbf9afe5574a596ad000000000000000000000000000000000000000000000000000000000000000028ef4fdffb91cdfc80b5c397a086cd81d6a171dfd0ade48fd448243d8bc3686801020200010c020000000100000000000000
		// fmt.Println(hex.EncodeToString(msg))

		//[1 0 1 3 53 219 114 19 228 91 116 152 162 89 178 55 147 57 156 40 33 213 109 243 88 86 244 54 216 124 147 42 227 150 3 76 71 79 115 53 213 57 158 73 101 102 252 254 137 176 109 223 47 157 243 31 171 96 26 175 191 154 254 85 116 165 150 173 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 200 196 8 222 236 199 174 150 78 150 17 5 177 231 191 139 183 168 69 6
		fmt.Println(msg)

		hash := string(msg)
		fmt.Println(hash)
		networkId := r.URL.Query().Get("networkId")
		go generateSignatureMessage(networkId, msg)

		messageChan[hash] = make(chan Message)

		sig := <-messageChan[hash]
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode("{\"signature\":\"" + string(sig.Message) + "\",\"address\":\"" + sig.Address + "\"}")
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		}

		delete(messageChan, hash)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
