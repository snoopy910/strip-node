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
	accountFrom := solana.MustPublicKeyFromBase58("oz2g94bsqEcgHKaDtsiN9Gi2DkMqJpZU6FZUy87GcUX")
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

	sig, _ := base58.Decode("3SCUJAbErrWQK7AZvYMn3dbZyjqGPAPoBRZgrDc7vrvSVC7ZVZffvSE8HixNKJctAVJuSffob7EeVduiawLoY6pK")
	signature := solana.SignatureFromBytes(sig)

	// tx.Signatures = append(tx.Signatures, signature)
	// err = tx.VerifySignatures()

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Signatures verified")
	// }

	// after signing marshal it again
	_msg, err := tx.ToBase64()
	if err != nil {
		panic(err)
	}

	_msgBytes, _ := base64.StdEncoding.DecodeString(_msg)
	_msgBase58 := base58.Encode(_msgBytes)

	fmt.Println(_msgBase58)

	// solana.TransactionFromDecoder(bin.NewBinDecoder(msg))

	// encoded := base64.StdEncoding.EncodeToString(msg)
	// fmt.Println(encoded)

	decodedTransactionData, err := base58.Decode(_msgBase58)
	if err != nil {
		fmt.Println("Error decoding transaction data:", err)
		return
	}

	// var decodedTransaction types.Transaction
	// err = decodedTransaction.UnmarshalBinary(decodedTransactionData)
	// if err != nil {
	// 	fmt.Println("Error unmarshalling transaction data:", err)
	// 	return
	// }

	// signature
	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		panic(err)
	}

	_tx.Signatures = append(_tx.Signatures, signature)

	// fmt.Println(_tx)

	err = _tx.VerifySignatures()

	if err != nil {
		fmt.Println("error during verification")
		fmt.Println(err)
	} else {
		fmt.Println("Signatures verified")
	}

	// tx.Message.MarshalBinary() gives the hash to sign
	// tx.ToBase64() is the actual transaction data

}
