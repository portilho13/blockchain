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

const initialDifficulty = 4
const adjustmentInterval = 5                                // Adjust difficulty every 5 blocks
const targetBlockTime = 60                                  // Target: 60 seconds per block
const targetInterval = adjustmentInterval * targetBlockTime // 300 seconds total
const maxDifficultyChange = 7                               // Maximum difficulty change per adjustment

type Blockchain struct {
	Blocks []models.Block
}

func (b *Blockchain) AddBlock(t transaction.Transaction) error {
	var prevHash string
	var difficulty int

	if len(b.Blocks) == 0 {
		prevHash = strings.Repeat("0", 64)
		difficulty = initialDifficulty
	} else {
		prevHash = b.Blocks[len(b.Blocks)-1].BlockHeader.Hash
		difficulty = b.CheckDifficulty()
	}

	new_transactions := []transaction.Transaction{t}

	merkle_root, err := crypto.CalculateMerkleRoot(new_transactions)
	if err != nil {
		return err
	}

	bh := models.BlockHeader{
		Version:          "1",
		PrevHash:         prevHash,
		MerkleRoot:       merkle_root,
		CreatedAt:        time.Now().UTC(),
		DifficultyTarget: difficulty,
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

func (b *Blockchain) CheckDifficulty() int {
	blocks := b.Blocks
	n := len(blocks)

	if n < adjustmentInterval {
		return initialDifficulty
	}

	if n%adjustmentInterval != 0 {
		return blocks[n-1].BlockHeader.DifficultyTarget
	}

	latestBlock := blocks[n-1]
	startBlock := blocks[n-adjustmentInterval]

	actualTime := latestBlock.BlockHeader.CreatedAt.Sub(startBlock.BlockHeader.CreatedAt).Seconds()
	currentDifficulty := blocks[n-1].BlockHeader.DifficultyTarget

	adjustmentFactor := float64(targetInterval) / actualTime

	newDifficulty := int(float64(currentDifficulty) * adjustmentFactor)

	if newDifficulty < 1 {
		newDifficulty = 1
	}

	maxIncrease := currentDifficulty + maxDifficultyChange
	maxDecrease := currentDifficulty - maxDifficultyChange
	if maxDecrease < 1 {
		maxDecrease = 1
	}

	if newDifficulty > maxIncrease {
		newDifficulty = maxIncrease
	} else if newDifficulty < maxDecrease {
		newDifficulty = maxDecrease
	}

	return newDifficulty
}

func (bc *Blockchain) PrintBlockchain() {
	for i, block := range bc.Blocks {
		fmt.Printf("\nBlock #%d\n", i)
		fmt.Println("──────────────────────────────────────────────")
		fmt.Printf("Version:         %s\n", block.BlockHeader.Version)
		fmt.Printf("Timestamp:       %s\n", block.BlockHeader.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Previous Hash:   %s\n", block.BlockHeader.PrevHash)
		fmt.Printf("Merkle Root:     %s\n", block.BlockHeader.MerkleRoot)
		fmt.Printf("Difficulty:      %d\n", block.BlockHeader.DifficultyTarget)
		fmt.Printf("Nonce:           %d\n", block.BlockHeader.Nonce)
		fmt.Printf("Hash:            %s\n", block.BlockHeader.Hash)

		fmt.Println("Transactions:")
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
