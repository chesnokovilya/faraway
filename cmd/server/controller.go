package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/chesnokovilya/faraway/cmd/response"
	"github.com/chesnokovilya/faraway/pow"
)

func handleConnection(conn net.Conn, difficulty int, rm *response.ResponseManager) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(pow.Timeout))

	challenge, err := generateChallenge()
	if err != nil {
		log.Print(err)
		return
	}

	challengeMsg := fmt.Sprintf("%s%s%s%d\n", pow.ChallengePrefix, challenge, pow.Delimiter, difficulty)
	if _, err := conn.Write([]byte(challengeMsg)); err != nil {
		log.Print(err)
		return
	}
	fmt.Println(challengeMsg)

	response, err := readResponse(conn)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println(response)

	if valid, err := validatePoW(challenge, difficulty, response); !valid {
		if err != nil {
			conn.Write([]byte(pow.ErrorPrefix + err.Error() + "\n"))
		} else {
			conn.Write([]byte(pow.FailPrefix + "Invalid proof of work"))
		}
		return
	}
	random := rm.GetRandomResponse()

	conn.Write([]byte(pow.WowPrefix + random + "\n"))
}

func generateChallenge() (string, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func readResponse(conn net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buffer[:n])), nil
}

func validatePoW(challenge string, difficulty int, response string) (bool, error) {
	parts := strings.Split(response, pow.Delimiter)
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
