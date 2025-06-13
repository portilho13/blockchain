package crypto

import "github.com/portilho13/blockchain/transaction"

func CalculateMerkleRoot(transactions []transaction.Transaction) (string, error) {
	level, err := HashTransaction(transactions...)
	if err != nil {
		return "", err
	}

	for len(level) > 1 {
		var nextLevel []string
		for i := 0; i < len(level); i += 2 {
			left := level[i]
			var right string

			if i+1 < len(level) {
				right = level[i+1]
			} else {
				right = left
			}

			combinedHash := HashPair(left, right)
			nextLevel = append(nextLevel, combinedHash)
		}
		level = nextLevel
	}

	return level[0], nil
}
