package block

import "time"

type BlockHeader struct {
	Version         string
	PrevHash        string
	MerkleRoot      string
	CreatedAt       time.Time
	DiffcultyTarget string
	Nonce           string
}
