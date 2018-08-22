package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/shopspring/decimal"
	"github.com/messwith/coding_challenge/models"
	"github.com/messwith/coding_challenge/errors"
	errors2 "errors"
)

// GetPayments is an endpoint for getting all the payments in system without any filtering
func (s *Server) GetPayments(c *gin.Context) {
	payments, err := s.paymentRep.GetPayments()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, payments)
}

// PaymentRequest describes request format for POST /payments endpoint
type PaymentRequest struct {
	FromAccount string `json:"from_account"`
	ToAccount string `json:"to_account"`
	Amount decimal.Decimal `json:"amount"`
}

func (s *Server) validatePayments(incoming, outgoing *models.Payment) error {
	if err := incoming.Validate(); err != nil {
		return err
	}

	if err := outgoing.Validate(); err != nil {
		return err
	}
	return nil
}

// CreatePayment is an endpoint for POST /payments which allow you to transfer funds from "from_account" to "to_account"
func (s *Server) CreatePayment(c *gin.Context) {
	pr := PaymentRequest{}
	if err := c.BindJSON(&pr); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	incomingPayment := models.NewIncomingPayment(pr.FromAccount, pr.ToAccount, pr.Amount)
	outgoingPayment := models.NewOutgoingPayment(pr.FromAccount, pr.ToAccount, pr.Amount)

	if err := s.validatePayments(incomingPayment, outgoingPayment); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := s.transactioner.Do(incomingPayment, outgoingPayment); err != nil {
		s.HandleTransactionError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) HandleTransactionError(c *gin.Context, err error) {
	switch convertedErr := err.(type) {
	case *errors.DataError:
		c.AbortWithError(http.StatusConflict, convertedErr)
		return
	case *errors.AccountNotFoundError:
		c.AbortWithError(http.StatusNotFound, convertedErr)
		return
	case *errors.InternalError:
		s.logger.Println("internal server error: "+convertedErr.Error())
		c.AbortWithError(http.StatusInternalServerError, errors2.New("internal server error"))
		return
	default:
		s.logger.Println("unknown error: " + err.Error())
		c.AbortWithError(http.StatusInternalServerError, errors2.New("internal server error"))
		return
	}
}
