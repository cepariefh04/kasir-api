package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
		return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetTransactionReport() (*models.TransactionReport, error) {
	return s.repo.FetchTransactionReport()
}

func (s *TransactionService) GetTransactionReportByDateRange(startDate, endDate string) (*models.TransactionReport, error) {
	return s.repo.FetchTransactionReportByDateRange(startDate, endDate)
}