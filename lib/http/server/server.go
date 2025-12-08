package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var instance *Server
var once sync.Once

type Server struct {
	serverHost   string        // example: www.example.com
	serverPort   int           // example: 80
	readTimeout  time.Duration // example: 60 (seconds)
	writeTimeout time.Duration // example: 60 (seconds)
	idleTimeout  time.Duration // example: 60 (seconds)
}

type ServerOption func(*Server) error

// Инициализация инстанса сервера
func GetInstance(options ...ServerOption) (IServer, error) {
	once.Do(func() {
		instance = &Server{
			serverHost:   "127.0.0.1",
			serverPort:   80,
			readTimeout:  60,
			writeTimeout: 60,
			idleTimeout:  60,
		}

		for _, opt := range options {
			if err := opt(instance); err != nil {
				instance = nil
				return
			}
		}
	})

	return instance, nil
}

func ServerHost(serverHost string) ServerOption {
	return func(s *Server) error {
		s.serverHost = serverHost
		return nil
	}
}

func ServerPort(serverPort int) ServerOption {
	return func(s *Server) error {
		if serverPort <= 0 || serverPort > 65535 {
			return fmt.Errorf("invalid port number: %d", serverPort)
		}
		s.serverPort = serverPort
		return nil
	}
}

func ReadTimeout(readTimeout int) ServerOption {
	return func(s *Server) error {
		s.readTimeout = time.Duration(readTimeout) * time.Second
		return nil
	}
}

func WriteTimeout(writeTimeout int) ServerOption {
	return func(s *Server) error {
		s.writeTimeout = time.Duration(writeTimeout) * time.Second
		return nil
	}
}

func IdleTimeout(idleTimeout int) ServerOption {
	return func(s *Server) error {
		s.idleTimeout = time.Duration(idleTimeout) * time.Second
		return nil
	}
}

func (s *Server) Up(router http.Handler) {
	addr := []string{s.serverHost, strconv.Itoa(s.serverPort)}

	server := &http.Server{
		Addr:         strings.Join(addr, ":"),
		Handler:      router,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
	}

	fmt.Println("listen connections...")
	log.Fatal(server.ListenAndServe())
}
