package api

import (
	"github.com/messwith/coding_challenge/repository"
	"log"
	"github.com/gin-gonic/gin"
	"errors"
	"fmt"
)

// Server is an object which stores all dependencies for api methods
type Server struct {
	logger *log.Logger
	accountRep *repository.AccountRepository
	paymentRep *repository.PaymentRepository
	router *gin.Engine
}

// NewServer creates ready to use server object with specified dependencies
func NewServer(accountRep *repository.AccountRepository, paymentRep *repository.PaymentRepository, logger *log.Logger) *Server {
	server := &Server{
		logger: logger,
		accountRep: accountRep,
		paymentRep: paymentRep,
	}
	router := gin.New()
	router.GET("/accounts", server.GetAccounts)
	router.GET("/payments", server.GetPayments)
	router.POST("/payments", server.CreatePayment)
	server.router = router
	return server
}

// Run runs the server at the specified port and blocks forever if there is no start error
func (s *Server) Run(port int) error {
	if port < 0 || port > 65535 {
		return errors.New("invalid port")
	}
	s.logger.Println("Starting server on port ", port)
	return s.router.Run(fmt.Sprintf(":%d", port))
}
