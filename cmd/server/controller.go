package server

import (
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
	challenge, err := pow.GenerateChallenge()
	if err != nil {
		log.Print(err)
		return
	}
	challengeMsg := fmt.Sprintf("%s%s%s%d\n", pow.ChallengePrefix, challenge, pow.Delimiter, difficulty)
	if _, err := conn.Write([]byte(challengeMsg)); err != nil {
		log.Print(err)
		return
	}
	response, err := readResponse(conn)
	if err != nil {
		log.Print(err)
		return
	}
	valid, err := pow.ValidatePoW(challenge, difficulty, response)
	if err != nil {
		conn.Write([]byte(pow.ErrorPrefix + err.Error() + "\n"))
		return
	}
	if !valid {
		conn.Write([]byte(pow.FailPrefix + "Invalid"))
		return
	}
	random := rm.GetRandomResponse()
	conn.Write([]byte(pow.WowPrefix + random + "\n"))
}

func readResponse(conn net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buffer[:n])), nil
}
