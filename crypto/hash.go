package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/portilho13/blockchain/transaction"
)

func HashTransaction(transactions ...transaction.Transaction) ([]string, error) {
	var hashes []string

	for _, t := range transactions {
		data, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}

		hash := sha256.Sum256(data)
		hashes = append(hashes, hex.EncodeToString(hash[:]))
	}

	return hashes, nil
}

func HashPair(left, right string) string {
	leftBytes, _ := hex.DecodeString(left)
	rightBytes, _ := hex.DecodeString(right)
	combined := append(leftBytes, rightBytes...)
	hash := sha256.Sum256(combined)
	return hex.EncodeToString(hash[:])
}
