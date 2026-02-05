package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	TotalAmount int       `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	Details		 []TransactionDetail `json:"details,omitempty"`
}

type TransactionDetail struct {
	ID            int `json:"id"`
	TransactionID int `json:"transaction_id"`
	ProductID     int `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int `json:"quantity"`
	Subtotal         int `json:"price"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}