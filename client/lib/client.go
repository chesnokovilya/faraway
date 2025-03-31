package client

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/chesnokovilya/faraway/cmd/server"
	"github.com/chesnokovilya/faraway/pow"
)

type PoWClient struct {
	serverAddress string
	timeout       time.Duration
}

func NewClient(serverAddress string) *PoWClient {
	return &PoWClient{
		serverAddress: serverAddress,
		timeout:       pow.Timeout,
	}
}

func (c *PoWClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.serverAddress, c.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(c.timeout))

	challenge, difficulty, err := c.receiveChallenge(conn)
	if err != nil {
		return err
	}

	fmt.Printf("Received challenge: %s (difficulty: %d)\n", challenge, difficulty)

	startTime := time.Now()
	solution := c.solveChallenge(challenge, difficulty)
	elapsed := time.Since(startTime)

	fmt.Printf("Solved in %v: nonce=%s hash=%s\n", elapsed, solution.Nonce, solution.Hash)

	if err := c.sendSolution(conn, solution); err != nil {
		return err
	}

	if err := c.readServerResponse(conn); err != nil {
		return err
	}

	return nil
}

func (c *PoWClient) receiveChallenge(conn net.Conn) (string, int, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", 0, err
	}

	msg := strings.TrimSpace(string(buf[:n]))
	if !strings.HasPrefix(msg, pow.ChallengePrefix) {
		return "", 0, errors.New("invalid challenge format")
	}

	parts := strings.Split(msg, pow.Delimiter)
	if len(parts) != 3 {
		return "", 0, errors.New("invalid challenge parts")
	}

	difficulty, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", 0, fmt.Errorf("invalid difficulty: %w", err)
	}

	return parts[1], difficulty, nil
}

func (c *PoWClient) solveChallenge(challenge string, difficulty int) pow.Solution {
	var nonce uint64
	for {
		nonceStr := strconv.FormatUint(nonce, 10)
		hash := sha256.Sum256([]byte(challenge + nonceStr))
		if server.CheckLeadingZeroBits(hash[:], difficulty) {
			return pow.Solution{
				Nonce: nonceStr,
				Hash:  hex.EncodeToString(hash[:]),
			}
		}
		nonce++
	}
}

func (c *PoWClient) sendSolution(conn net.Conn, solution pow.Solution) error {
	response := fmt.Sprintf("%s%s%s", solution.Nonce, pow.Delimiter, solution.Hash)
	_, err := conn.Write([]byte(response + "\n"))
	return err
}

func (c *PoWClient) readServerResponse(conn net.Conn) error {
	buf := make([]byte, 10000)
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			return errors.New("connection closed by server")
		}
		return err
	}

	msg := strings.TrimSpace(string(buf[:n]))

	switch {
	case strings.HasPrefix(msg, pow.WowPrefix):
		fmt.Println(strings.TrimPrefix(msg, pow.WowPrefix))
	default:
		fmt.Println("Unknown response:", msg)
	}

	return nil
}
