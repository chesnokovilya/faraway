package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

const (
	ChallengePrefix = "CHALLENGE:"
	FailPrefix      = "FAIL:"
	ErrorPrefix     = "ERROR:"
	WowPrefix       = "WOW:"
	Delimiter       = ":"
	Difficulty      = 24
	Timeout         = 10 * time.Second
)

type Challenge struct {
	Value      string
	Difficulty int
}

type Solution struct {
	Nonce string
	Hash  string
}

func CheckLeadingZeroBits(hash []byte, difficulty int) bool {
	requiredBytes := difficulty / 8
	remainingBits := difficulty % 8
	for i := 0; i < requiredBytes; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	if remainingBits > 0 && requiredBytes < len(hash) {
		mask := byte(0xFF) << (8 - remainingBits)
		return (hash[requiredBytes] & mask) == 0
	}
	return true
}

func ValidatePoW(challenge string, difficulty int, response string) (bool, error) {
	parts := strings.Split(response, Delimiter)
	if len(parts) != 2 {
		return false, errors.New("invalid response format")
	}
	nonce := parts[0]
	hash := parts[1]
	data := challenge + nonce
	computedHash := sha256.Sum256([]byte(data))
	computedHashStr := hex.EncodeToString(computedHash[:])
	if computedHashStr != hash {
		return false, nil
	}
	return CheckLeadingZeroBits(computedHash[:], difficulty), nil
}

func GenerateChallenge() (string, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
