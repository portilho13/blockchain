package block

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/portilho13/blockchain/crypto"
	"github.com/portilho13/blockchain/models"
	"github.com/portilho13/blockchain/transaction"
)

const GENESIS_DIFFICULTY = 4

type Blockchain struct {
	Blocks []models.Block
}

func (b *Blockchain) AddBlock(t transaction.Transaction) error {
	var prevHash string

	var difficulty int

	var new_transactions []transaction.Transaction

	if len(b.Blocks) == 0 {
		prevHash = strings.Repeat("0", 64)
		difficulty = GENESIS_DIFFICULTY
	} else {
		prevHash = b.Blocks[len(b.Blocks)-1].BlockHeader.Hash // Last block hash
		difficulty = GENESIS_DIFFICULTY                       // Implement dynamic difficulty value later

		new_transactions = b.Blocks[len(b.Blocks)-1].BlockBody.Transactions
	}

	new_transactions = append(new_transactions, t)

	merkle_root, err := crypto.CalculateMerkleRoot(new_transactions)
	if err != nil {
		return err
	}

	bh := models.BlockHeader{
		Version:         "1",
		PrevHash:        prevHash,
		MerkleRoot:      merkle_root,
		CreatedAt:       time.Now().UTC(),
		DiffcultyTarget: difficulty,
	}

	err = crypto.HashBlockHeaderWithNonce(&bh)
	if err != nil {
		return err
	}

	bb := models.BlockBody{
		Transactions: new_transactions,
	}

	block := models.Block{
		BlockHeader: bh,
		BlockBody:   bb,
	}

	b.Blocks = append(b.Blocks, block)

	return nil
}

func (bc *Blockchain) PrintBlockchain() {
	for i, block := range bc.Blocks {
		fmt.Printf("\n🧱 Block #%d\n", i)
		fmt.Println("──────────────────────────────────────────────")
		fmt.Printf("Version:         %s\n", block.BlockHeader.Version)
		fmt.Printf("Timestamp:       %s\n", block.BlockHeader.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Previous Hash:   %s\n", block.BlockHeader.PrevHash)
		fmt.Printf("Merkle Root:     %s\n", block.BlockHeader.MerkleRoot)
		fmt.Printf("Difficulty:      %d\n", block.BlockHeader.DiffcultyTarget)
		fmt.Printf("Nonce:           %d\n", block.BlockHeader.Nonce)
		fmt.Printf("Hash:            %s\n", block.BlockHeader.Hash)

		fmt.Println("📦 Transactions:")
		if len(block.BlockBody.Transactions) == 0 {
			fmt.Println("  (none)")
		} else {
			for j, tx := range block.BlockBody.Transactions {
				txJSON, _ := json.MarshalIndent(tx, "    ", "  ")
				fmt.Printf("  Tx #%d:\n%s\n", j+1, string(txJSON))
			}
		}
		fmt.Println("──────────────────────────────────────────────")
	}
}
