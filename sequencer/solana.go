package sequencer

import (
	"encoding/base64"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/mr-tron/base58"
)

func TestBuildSolana() {
	// pubKeyByte, _ := base58.Decode("4dEgqPG9FtjCiwW1HUdReraBozwq6qUCcDFXD8BnUn9Z")

	// c := rpc.New("https://api.devnet.solana.com")
	// recentHash, err := c.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(recentHash.Value.Blockhash.String())

	accountFrom := solana.MustPublicKeyFromBase58("DpZqkyDKkVv2S7Lhbd5dUVcVCPJz2Lypr4W5Cru2sHr7")
	accountTo := solana.MustPublicKeyFromBase58("5oNDL3swdJJF1g9DzJiZ4ynHXgszjAEpUkxVYejchzrY")
	amount := uint64(1)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				accountFrom,
				accountTo,
			).Build(),
		},
		solana.MustHashFromBase58("CBLp4VEPu9T9W2uzURoawLGqgAQ65LvmUwDYRHymgwbd"),
		solana.TransactionPayer(accountFrom),
	)

	if err != nil {
		panic(err)
	}

	// msg, err := tx.Message.MarshalBinary()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Input to sign: ", base58.Encode(msg))

	// fmt.Println("Input to sign: ", base58.Encode(msg))

	sig, _ := base58.Decode("5jLFtNTCAnHA9uurWhyNNqzwHLwWCaSNrZBWG48AANMGkreX1DYGbkHL2VWNNt2Kz327QwzzsAacJj2YFdSsfkwN")
	signature := solana.SignatureFromBytes(sig)

	_msg, err := tx.ToBase64()
	if err != nil {
		panic(err)
	}

	_msgBytes, _ := base64.StdEncoding.DecodeString(_msg)
	_msgBase58 := base58.Encode(_msgBytes)

	fmt.Println(_msgBase58)

	decodedTransactionData, err := base58.Decode(_msgBase58)
	if err != nil {
		fmt.Println("Error decoding transaction data:", err)
		return
	}

	// // signature
	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		panic(err)
	}

	_tx.Signatures = append(_tx.Signatures, signature)

	err = _tx.VerifySignatures()

	if err != nil {
		fmt.Println("error during verification")
		fmt.Println(err)
	} else {
		fmt.Println("Signatures verified")
	}
}
