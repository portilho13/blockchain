package main

import (
	"fmt"
	"log"

	"github.com/portilho13/blockchain/crypto"

	"github.com/portilho13/blockchain/transaction"
)

func main() {
	mockTx := transaction.Transaction{
		TxId: "abc123def456ghi789", // fake transaction ID
		Input: []transaction.Input{
			{
				PrevTxId:    "prevtx0001",
				OutputIndex: 0,
				ScriptSig:   "3045022100abcdef... [signature]",
			},
		},
		Output: []transaction.Output{
			{
				Value:        1.25,
				ScriptPubKey: "76a91489abcdefabbaabbaabbaabbaabbaabbaabba88ac", // P2PKH-like
			},
			{
				Value:        0.75,
				ScriptPubKey: "76a914abcdef1234567890abcdef1234567890abcd88ac", // change address
			},
		},
	}

	s, err := crypto.HashTransaction(mockTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)

	var t []transaction.Transaction

	t = append(t, mockTx)

	m, _ := crypto.CalculateMerkleRoot(t)
	fmt.Println(m)

}
