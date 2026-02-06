package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}
  
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)")
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		for i := range details {
			details[i].TransactionID = transactionID
			_, err = stmt.Exec(transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
			if err != nil {
				return nil, err
			}
		}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) FetchTransactionReport() (*models.TransactionReport, error) {
	query := `
		WITH best_product AS (
			SELECT 
				p.name,
				SUM(td.quantity) AS total_qty
			FROM transaction_details td
			JOIN products p ON td.product_id = p.id
			JOIN transactions t ON td.transaction_id = t.id
			WHERE DATE(t.created_at) = CURRENT_DATE
			GROUP BY p.id, p.name
			ORDER BY total_qty DESC
			LIMIT 1
		)
		SELECT 
			COALESCE(SUM(t.total_amount), 0) AS total_revenue,
			COALESCE(COUNT(t.id), 0) AS total_transactions,
			COALESCE((SELECT name FROM best_product), '') AS product_name,
			COALESCE((SELECT total_qty FROM best_product), 0) AS product_qty
		FROM transactions t
		WHERE DATE(t.created_at) = CURRENT_DATE
	`

	var report models.TransactionReport
	var productName string
	var productQty int

	err := repo.db.QueryRow(query).Scan(
		&report.TotalRevenue,
		&report.TotalTransactions,
		&productName,
		&productQty,
	)
	if err != nil {
		return nil, err
	}

	report.BestSellingProduct = models.BestSellingProduct{
		Name:     productName,
		Quantity: productQty,
	}

	return &report, nil
}

func (repo *TransactionRepository) FetchTransactionReportByDateRange(startDate, endDate string) (*models.TransactionReport, error) {
	query := `
		WITH best_product AS (
			SELECT 
				p.name,
				SUM(td.quantity) AS total_qty
			FROM transaction_details td
			JOIN products p ON td.product_id = p.id
			JOIN transactions t ON td.transaction_id = t.id
			WHERE DATE(t.created_at) BETWEEN $1 AND $2
			GROUP BY p.id, p.name
			ORDER BY total_qty DESC
			LIMIT 1
		)
		SELECT 
			COALESCE(SUM(t.total_amount), 0) AS total_revenue,
			COALESCE(COUNT(t.id), 0) AS total_transactions,
			COALESCE((SELECT name FROM best_product), '') AS product_name,
			COALESCE((SELECT total_qty FROM best_product), 0) AS product_qty
		FROM transactions t
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
	`

	var report models.TransactionReport
	var productName string
	var productQty int

	err := repo.db.QueryRow(query, startDate, endDate).Scan(
		&report.TotalRevenue,
		&report.TotalTransactions,
		&productName,
		&productQty,
	)
	if err != nil {
		return nil, err
	}

	report.BestSellingProduct = models.BestSellingProduct{
		Name:     productName,
		Quantity: productQty,
	}

	return &report, nil
}