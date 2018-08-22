package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/shopspring/decimal"
	"github.com/messwith/coding_challenge/models"
	"database/sql"
	"errors"
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

func (s *Server) createPayments(tx *sql.Tx, incoming, outgoing *models.Payment) error {
	if err := s.paymentRep.CreatePayment(tx, outgoing); err != nil {
		return err
	}

	if err := s.paymentRep.CreatePayment(tx, incoming); err != nil {
		return err
	}
	return nil
}

func (s *Server) updateAccountBalances(tx *sql.Tx, sender *models.Account, receiver *models.Account) error {
	if err := s.accountRep.UpdateAccountBalance(tx, sender.ID, sender.Balance); err != nil {
		return err
	}

	if err := s.accountRep.UpdateAccountBalance(tx, receiver.ID, receiver.Balance); err != nil {
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

	tx, err := s.paymentRep.BeginTx()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	sender, err := s.accountRep.LockAccount(tx, pr.FromAccount)
	if err == sql.ErrNoRows{
		c.AbortWithError(http.StatusNotFound, err)
		return
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	sender.Balance = sender.Balance.Sub(pr.Amount)
	if sender.Balance.Cmp(decimal.Zero) == -1 {
		c.AbortWithError(http.StatusConflict, errors.New("sender has not enough funds"))
		return
	}

	receiver, err := s.accountRep.LockAccount(tx, pr.ToAccount)
	if err == sql.ErrNoRows{
		c.AbortWithError(http.StatusNotFound, err)
		return
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if sender.Currency != receiver.Currency {
		c.AbortWithError(http.StatusConflict, errors.New("sender and receiver has different currencies, conversion is not supported yet"))
		return
	}
	receiver.Balance = receiver.Balance.Add(pr.Amount)

	if err := s.createPayments(tx, incomingPayment, outgoingPayment); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := s.updateAccountBalances(tx, &sender, &receiver); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
