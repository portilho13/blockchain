package block

import (
	"time"

	"github.com/portilho13/blockchain/transaction"
)

type BlockHeader struct {
	Version         string
	PrevHash        string
	MerkleRoot      string
	CreatedAt       time.Time
	DiffcultyTarget string
	Nonce           string
}

type BlockBody struct {
	Transactions []transaction.Transaction
}

type Block struct {
	BlockHeader BlockHeader
	BlockBody   BlockBody
}
