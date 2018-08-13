package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/shopspring/decimal"
	"github.com/messwith/coding_challenge/models"
	"database/sql"
	"errors"
)

func (s *Server) GetPayments(c *gin.Context) {
	payments, err := s.paymentRep.GetPayments()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, payments)
}

type PaymentRequest struct {
	FromAccount string `json:"from_account"`
	ToAccount string `json:"to_account"`
	Amount decimal.Decimal `json:"amount"`
}

func (s *Server) CreatePayment(c *gin.Context) {
	pr := PaymentRequest{}
	if err := c.BindJSON(&pr); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	incomingPayment := models.NewIncomingPayment(pr.FromAccount, pr.ToAccount, pr.Amount)
	outgoingPayment := models.NewOutgoingPayment(pr.FromAccount, pr.ToAccount, pr.Amount)

	if err := incomingPayment.Validate(); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := outgoingPayment.Validate(); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx, err := s.paymentRep.BeginTx()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	sender, err := s.accountRep.LockAccount(pr.FromAccount)
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

	receiver, err := s.accountRep.LockAccount(pr.ToAccount)
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

	if err := s.paymentRep.CreatePayment(tx, outgoingPayment); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := s.paymentRep.CreatePayment(tx, incomingPayment); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := s.accountRep.UpdateAccountBalance(sender.Id, sender.Balance); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := s.accountRep.UpdateAccountBalance(receiver.Id, receiver.Balance); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
