package models

import (
	"time"
)

type BlockHeader struct {
	Version          string
	PrevHash         string
	MerkleRoot       string
	CreatedAt        time.Time
	DifficultyTarget int
	Nonce            int
	Hash             string
}

type BlockBody struct {
	Transactions []ConfirmedTransaction
}

type Block struct {
	BlockHeader BlockHeader
	BlockBody   BlockBody
}
