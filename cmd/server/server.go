package server

import (
	"fmt"
	"log"
	"net"
	"path/filepath"
	"runtime"

	"github.com/chesnokovilya/faraway/cmd/response"
	"github.com/chesnokovilya/faraway/pow"
)

type TCPServer struct {
	address         string
	difficulty      int
	connections     int
	listener        net.Listener
	responseManager *response.ResponseManager
	done            chan struct{}
}

func NewServer(address string) (*TCPServer, error) {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(filename))
	log.Println("rootDit=" + rootDir)
	responsePath := filepath.Join(rootDir, "../", "lorem-ipsum.txt")
	rm, err := response.NewResponseManager(responsePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load responses: %w", err)
	}
	return &TCPServer{
		address:         address,
		difficulty:      pow.Difficulty,
		done:            make(chan struct{}),
		responseManager: rm,
	}, nil
}

func (s *TCPServer) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	log.Printf("Server started on %s", s.address)
	go s.acceptConnections()
	return nil
}

func (s *TCPServer) Stop() {
	close(s.done)
	s.listener.Close()
}

func (s *TCPServer) acceptConnections() {
	for {
		select {
		case <-s.done:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
			s.connections++
			go handleConnection(conn, s.difficulty, s.responseManager)
		}
	}
}
