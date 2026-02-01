package models

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	IDCategory int    `json:"id_category"`

	Category *Category `json:"category,omitempty"`
}

type ProductWithCategory struct {
	Product
	CategoryName string `json:"category_name"`
}