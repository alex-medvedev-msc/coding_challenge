package service

import (
	"github.com/messwith/coding_challenge/models"
	"github.com/messwith/coding_challenge/repository"
	"github.com/messwith/coding_challenge/errors"
)

// PaymentService is an interface which allows to abstract from details of storing payments in some storage
type PaymentService interface {
	GetPayments() ([]models.Payment, error)
}

// SqlPaymentService uses sql db as storage, currently it is postgres
// but you can change it with simple import replacement, e.g. _ ".../lib/pq" -> _ ".../lib/mysql"
type SqlPaymentService struct {
	paymentRep *repository.PaymentRepository
}

// NewSqlPaymentService creates ready to use SqlPaymentService instance
func NewSqlPaymentService(
	paymentRep *repository.PaymentRepository) *SqlPaymentService {

	return &SqlPaymentService{
		paymentRep: paymentRep,
	}
}

// GetPayments returns all payments from db without filtering and pagination
func (ps *SqlPaymentService) GetPayments() ([]models.Payment, error) {

	payments, err := ps.paymentRep.GetPayments()
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	return payments, nil
}