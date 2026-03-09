package services

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"multitenancypfe/internal/products/dto"
	"multitenancypfe/internal/products/models"
	"multitenancypfe/internal/products/repo"
)

type CategoryService interface {
	Create(storeID uuid.UUID, req dto.CreateCategoryRequest, db *gorm.DB) (*dto.CategoryResponse, error)
	GetByID(id, storeID uuid.UUID, db *gorm.DB) (*dto.CategoryResponse, error)
	GetTree(storeID uuid.UUID, db *gorm.DB) ([]dto.CategoryResponse, error)
	Update(id, storeID uuid.UUID, req dto.UpdateCategoryRequest, db *gorm.DB) (*dto.CategoryResponse, error)
	Delete(id, storeID uuid.UUID, db *gorm.DB) error
}

type categoryService struct {
	repo repo.CategoryRepository
}

func NewCategoryService(r repo.CategoryRepository) CategoryService {
	return &categoryService{repo: r}
}

func (s *categoryService) Create(storeID uuid.UUID, req dto.CreateCategoryRequest, db *gorm.DB) (*dto.CategoryResponse, error) {
	slug := resolveSlug(req.Slug, req.Name)

	exists, err := s.repo.SlugExists(slug, storeID, nil, db)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("slug already in use")
	}

	c := &models.Category{
		StoreID:     storeID,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Visibility:  req.Visibility,
	}
	if err := s.repo.Create(c, db); err != nil {
		return nil, err
	}
	resp := toCategoryResponse(c)
	return &resp, nil
}

func (s *categoryService) GetByID(id, storeID uuid.UUID, db *gorm.DB) (*dto.CategoryResponse, error) {
	c, err := s.findOrFail(id, storeID, db)
	if err != nil {
		return nil, err
	}
	resp := toCategoryResponse(c)
	return &resp, nil
}

// GetTree returns all root categories with nested children (arbitrary depth).
func (s *categoryService) GetTree(storeID uuid.UUID, db *gorm.DB) ([]dto.CategoryResponse, error) {
	roots, err := s.repo.FindRoots(storeID, db)
	if err != nil {
		return nil, err
	}
	result := make([]dto.CategoryResponse, len(roots))
	for i, c := range roots {
		result[i] = toCategoryResponse(&c)
	}
	return result, nil
}

func (s *categoryService) Update(id, storeID uuid.UUID, req dto.UpdateCategoryRequest, db *gorm.DB) (*dto.CategoryResponse, error) {
	c, err := s.findOrFail(id, storeID, db)
	if err != nil {
		return nil, err
	}

	if req.Slug != nil {
		exists, err := s.repo.SlugExists(*req.Slug, storeID, &id, db)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("slug already in use")
		}
		c.Slug = *req.Slug
	}

	if req.ParentID != nil {
		c.ParentID = req.ParentID
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.Visibility != nil {
		c.Visibility = *req.Visibility
	}

	if err := s.repo.Update(c, db); err != nil {
		return nil, err
	}
	resp := toCategoryResponse(c)
	return &resp, nil
}

// Delete blocks deletion if the category still has products assigned to it.
func (s *categoryService) Delete(id, storeID uuid.UUID, db *gorm.DB) error {
	if _, err := s.findOrFail(id, storeID, db); err != nil {
		return err
	}
	hasProducts, err := s.repo.HasProducts(id, db)
	if err != nil {
		return err
	}
	if hasProducts {
		return errors.New("cannot delete category with products")
	}
	return s.repo.Delete(id, storeID, db)
}

func (s *categoryService) findOrFail(id, storeID uuid.UUID, db *gorm.DB) (*models.Category, error) {
	c, err := s.repo.FindByID(id, storeID, db)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("category not found")
	}
	return c, nil
}
