package api

import (
	"github.com/messwith/coding_challenge/repository"
	"log"
)

// Server is an object which stores all dependencies for api methods
type Server struct {
	logger *log.Logger
	accountRep *repository.AccountRepository
	paymentRep *repository.PaymentRepository
}

// NewServer creates ready to use server object with specified dependencies
func NewServer(accountRep *repository.AccountRepository, paymentRep *repository.PaymentRepository, logger *log.Logger) *Server {
	return &Server{
		logger: logger,
		accountRep: accountRep,
		paymentRep: paymentRep,
	}
}

// Run runs the server at the specified port and blocks forever if there is no start error
func (s *Server) Run(port int) error {
	s.logger.Println("Test run on port ", port)
	return nil
}
