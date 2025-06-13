package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/portilho13/blockchain/models"
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

func HashBlockHeader(bh models.BlockHeader) (string, error) {
	data, err := json.Marshal(bh)
	if err != nil {
		return "", err
	}

	hash := sha256.New()

	hash.Write(data)

	bs := hash.Sum(nil)

	return hex.EncodeToString(bs[:]), nil
}

func CheckDifficulty(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)

	return strings.HasPrefix(hash, prefix)
}

func HashBlockHeaderWithNonce(b *models.BlockHeader) error {
	nonce := 0
	for {
		b.Nonce = nonce
		hash, err := HashBlockHeader(*b)
		if err != nil {
			return err
		}

		b.Hash = hash

		if CheckDifficulty(hash, b.DifficultyTarget) {
			break
		}
		nonce += 1

	}
	return nil
}
