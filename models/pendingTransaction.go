package models

import "time"

type PendingTransaction struct {
	Tx       *Transaction
	Received time.Time
	Fee      uint64
}
