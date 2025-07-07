package block

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/portilho13/blockchain/crypto"
	"github.com/portilho13/blockchain/models"
)

const initialDifficulty = 4
const adjustmentInterval = 5                                // Adjust difficulty every 5 blocks
const targetBlockTime = 60                                  // Target: 60 seconds per block
const targetInterval = adjustmentInterval * targetBlockTime // 300 seconds total
const maxDifficultyChange = 7                               // Maximum difficulty change per adjustment
const transactionsPerBlock = 5

type Blockchain struct {
	Blocks []models.Block
}

func (b *Blockchain) MintBlock(m *models.Mempool) (*models.Block, error) {
	var pendingTransactions []models.PendingTransaction
	if len(*m.Tx) > transactionsPerBlock {
		pendingTransactions = (*m.Tx)[:transactionsPerBlock]
	} else {
		pendingTransactions = *m.Tx
	}

	var prevHash string
	var difficulty int

	if len(b.Blocks) == 0 { // If is first block
		prevHash = strings.Repeat("0", 64)
		difficulty = initialDifficulty
	} else {
		prevHash = b.Blocks[len(b.Blocks)-1].BlockHeader.Hash
		difficulty = b.CheckDifficulty()
	}

	var txs []models.Transaction
	for _, tx := range pendingTransactions {
		txs = append(txs, *tx.Tx)
	}

	merkle_root, err := crypto.CalculateMerkleRoot(txs)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	var confirmedTransactions []models.ConfirmedTransaction

	for idx, tx := range pendingTransactions {
		confirmedTx := models.ConfirmedTransaction{
			Tx:          tx.Tx,
			BlockHash:   bh.Hash,
			BlockHeight: uint32(idx + 1),
			ConfirmedAt: time.Now().UTC(),
		}

		confirmedTransactions = append(confirmedTransactions, confirmedTx)
	}

	bb := models.BlockBody{
		Transactions: confirmedTransactions,
	}

	block := models.Block{
		BlockHeader: bh,
		BlockBody:   bb,
	}

	b.Blocks = append(b.Blocks, block)

	if len(*m.Tx) > transactionsPerBlock {
		*m.Tx = (*m.Tx)[transactionsPerBlock:] // Reduce mempool
	} else {
		*m.Tx = nil // Mempool is empty
	}

	return &block, nil
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

func (b *Blockchain) Addblock(block models.Block) {
	b.Blocks = append(b.Blocks, block)
}
