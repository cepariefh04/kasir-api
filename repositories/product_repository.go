package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAllProducts(name string) ([]models.ProductWithCategory, error) {
	query := `SELECT 
							p.id, 
							p.name, 
							p.price, 
							p.stock, 
							p.id_category,
							COALESCE(c.name, '') as category_name
						FROM products p
						LEFT JOIN categories c ON p.id_category = c.id`

	var args []interface{}
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}
	
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.ProductWithCategory, 0)
	for rows.Next() {
		var p models.ProductWithCategory
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.IDCategory, &p.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) CreateProduct(p *models.Product) error {
	query := "INSERT INTO products (name, price, stock, id_category) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, p.Name, p.Price, p.Stock, p.IDCategory).Scan(&p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) GetProductById(id int) (*models.ProductWithCategory, error) {
	query := `SELECT 
							p.id, 
							p.name, 
							p.price, 
							p.stock, 
							p.id_category,
							COALESCE(c.name, '') as category_name
						FROM products p
						LEFT JOIN categories c ON p.id_category = c.id
						WHERE p.id = $1`
	
	var p models.ProductWithCategory
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.IDCategory, &p.CategoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) GetProductByIdCategory(id_category int) (*models.ProductWithCategory, error) {
	query := `SELECT 
							p.id, 
							p.name, 
							p.price, 
							p.stock, 
							p.id_category,
							COALESCE(c.name, '') as category_name
						FROM products p
						LEFT JOIN categories c ON p.id_category = c.id
						WHERE p.id_category = $1`
	
	var p models.ProductWithCategory
	err := repo.db.QueryRow(query, id_category).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.IDCategory, &p.CategoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) UpdateProduct(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}