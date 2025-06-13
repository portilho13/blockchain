package main

import (
	"github.com/portilho13/blockchain/block"
	"github.com/portilho13/blockchain/transaction"
)

const adjustmentInterval = 5 // Match this with your difficulty logic

func main() {
	mockTx := transaction.Transaction{
		TxId: "abc123def456ghi789",
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
				ScriptPubKey: "76a91489abcdefabbaabbaabbaabbaabbaabbaabba88ac",
			},
			{
				Value:        0.75,
				ScriptPubKey: "76a914abcdef1234567890abcdef1234567890abcd88ac",
			},
		},
	}

	bc := block.Blockchain{}
	bc.AddBlock(mockTx)

	bc.PrintBlockchain()

}
