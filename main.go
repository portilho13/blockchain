package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/portilho13/blockchain/block"
	"github.com/portilho13/blockchain/conn"
	"github.com/portilho13/blockchain/models"
)

const MAIN_DOMAIN = "localhost:8000"

func main() {
	mockTx := &models.PendingTransaction{
		Tx: &models.Transaction{
			TxId: "abc123def456ghi789",
			Input: []models.Input{
				{
					PrevTxId:    "prevtx0001",
					OutputIndex: 0,
					ScriptSig:   "3045022100abcdef... [signature]",
				},
			},
			Output: []models.Output{
				{
					Value:        1.25,
					ScriptPubKey: "76a91489abcdefabbaabbaabbaabbaabbaabbaabba88ac",
				},
				{
					Value:        0.75,
					ScriptPubKey: "76a914abcdef1234567890abcdef1234567890abcd88ac",
				},
			},
		},
		Received: time.Now(),
		Fee:      1000,
	}

	bc := block.Blockchain{}
	mem := &models.Mempool{Tx: &[]models.PendingTransaction{*mockTx}}

	args := os.Args
	conn := conn.Connection{BlockChain: &bc}

	go conn.Start(args[1])

	block, err := bc.MintBlock(mem)
	if err != nil {
		panic(err)
	}

	err = conn.BroadcastBlockToAll(*block)
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Println("Running forever...")
		time.Sleep(2 * time.Second)

		//bc.PrintBlockchain()
	}

}
