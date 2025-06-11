package crypto

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/portilho13/blockchain/transaction"
)

func HashTransaction(t transaction.Transaction) (string, error) {
	h := sha256.New()

	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	h.Write(data)

	bs := h.Sum(nil)

	s := fmt.Sprintf("%x", bs)

	return s, nil
}
