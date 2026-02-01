package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.GetAllCategories()
}

func (s *CategoryService) CreateCategory(c *models.Category) error {
	// Here you can add business logic before creating the category
	return s.repo.CreateCategory(c)
}

func (s *CategoryService) GetCategoryById(id int) (*models.Category, error) {
	return s.repo.GetCategoryById(id)
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
	return s.repo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.DeleteCategory(id)
}