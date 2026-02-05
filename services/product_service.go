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

func (s *ProductService) GetAllProducts(name string) ([]models.ProductWithCategory, error) {
	return s.repo.GetAllProducts(name)
}

func (s *ProductService) CreateProduct(p *models.Product) error {
	// Here you can add business logic before creating the product
	return s.repo.CreateProduct(p)
}

func (s *ProductService) GetProductById(id int) (*models.ProductWithCategory, error) {
	return s.repo.GetProductById(id)
}

func (s *ProductService) GetProductByIdCategory(id_category int) (*models.ProductWithCategory, error) {
	return s.repo.GetProductByIdCategory(id_category)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.repo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}