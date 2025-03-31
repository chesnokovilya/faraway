package response

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

type ResponseManager struct {
	responses []string
}

func NewResponseManager(filePath string) (*ResponseManager, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var responses []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		responses = append(responses, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	return &ResponseManager{responses: responses}, nil
}

func (rm *ResponseManager) GetRandomResponse() string {
	if len(rm.responses) == 0 {
		return "Welcome to the service!"
	}
	return rm.responses[rand.Intn(len(rm.responses))]
}
