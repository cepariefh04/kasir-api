package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *ProductService) CreateProduct(p *models.Product) error {
	// Here you can add business logic before creating the product
	return s.repo.CreateProduct(p)
}

func (s *ProductService) GetProductById(id int) (*models.Product, error) {
	return s.repo.GetProductById(id)
}