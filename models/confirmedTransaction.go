package models

import "time"

type ConfirmedTransaction struct {
	Tx          *Transaction
	BlockHash   string
	BlockHeight uint32
	ConfirmedAt time.Time
}
