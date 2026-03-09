package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"multitenancypfe/internal/products/dto"
	"multitenancypfe/internal/products/models"
	"multitenancypfe/internal/products/repo"
)

type ProductService interface {
	Create(storeID uuid.UUID, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetByID(id, storeID uuid.UUID) (*dto.ProductResponse, error)
	GetAll(storeID uuid.UUID, filter dto.ProductFilter) (*dto.ProductListResponse, error)
	Update(id, storeID uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	Delete(id, storeID uuid.UUID) error
}

type productService struct {
	repo repo.ProductRepository
}

func NewProductService(r repo.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (s *productService) Create(storeID uuid.UUID, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	slug := resolveSlug(req.Slug, req.Title)

	exists, err := s.repo.SlugExists(slug, storeID, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("slug already in use")
	}

	product := &models.Product{
		StoreID:     storeID,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		Slug:        slug,
		Status:      req.Status,
		Visibility:  req.Visibility,
		Price:       req.Price,
		Currency:    req.Currency,
		SKU:         req.SKU,
		TrackStock:  req.TrackStock,
		Stock:       req.Stock,
		Weight:      req.Weight,
		Dimensions:  req.Dimensions,
		Brand:       req.Brand,
		TaxClass:    req.TaxClass,
		PublishedAt: req.PublishedAt,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return toProductResponse(product), nil
}

func (s *productService) GetByID(id, storeID uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.findOrFail(id, storeID)
	if err != nil {
		return nil, err
	}
	return toProductResponse(product), nil
}

func (s *productService) GetAll(storeID uuid.UUID, filter dto.ProductFilter) (*dto.ProductListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}

	products, err := s.repo.FindAll(storeID, filter)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		result[i] = *toProductResponse(&p)
	}
	return &dto.ProductListResponse{Products: result, Page: filter.Page, Limit: filter.Limit}, nil
}

func (s *productService) Update(id, storeID uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.findOrFail(id, storeID)
	if err != nil {
		return nil, err
	}

	if req.Slug != nil {
		exists, err := s.repo.SlugExists(*req.Slug, storeID, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("slug already in use")
		}
		product.Slug = *req.Slug
	}

	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.Title != nil {
		product.Title = *req.Title
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.Visibility != nil {
		product.Visibility = *req.Visibility
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Currency != nil {
		product.Currency = *req.Currency
	}
	if req.SKU != nil {
		product.SKU = req.SKU
	}
	if req.TrackStock != nil {
		product.TrackStock = *req.TrackStock
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Weight != nil {
		product.Weight = req.Weight
	}
	if req.Dimensions != nil {
		product.Dimensions = req.Dimensions
	}
	if req.Brand != nil {
		product.Brand = req.Brand
	}
	if req.TaxClass != nil {
		product.TaxClass = req.TaxClass
	}
	if req.PublishedAt != nil {
		product.PublishedAt = req.PublishedAt
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return toProductResponse(product), nil
}

func (s *productService) Delete(id, storeID uuid.UUID) error {
	if _, err := s.findOrFail(id, storeID); err != nil {
		return err
	}
	return s.repo.Delete(id, storeID)
}

func (s *productService) findOrFail(id, storeID uuid.UUID) (*models.Product, error) {
	p, err := s.repo.FindByID(id, storeID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("product not found")
	}
	return p, nil
}

// ── shared mappers (used by category & collection services too) ───────────────

func resolveSlug(slug *string, title string) string {
	if slug != nil && *slug != "" {
		return *slug
	}
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(title), " ", "-"))
}

func toProductResponse(p *models.Product) *dto.ProductResponse {
	resp := &dto.ProductResponse{
		ID:          p.ID,
		StoreID:     p.StoreID,
		CategoryID:  p.CategoryID,
		Title:       p.Title,
		Description: p.Description,
		Slug:        p.Slug,
		Status:      p.Status,
		Visibility:  p.Visibility,
		Price:       p.Price,
		Currency:    p.Currency,
		SKU:         p.SKU,
		TrackStock:  p.TrackStock,
		Stock:       p.Stock,
		Weight:      p.Weight,
		Dimensions:  p.Dimensions,
		Brand:       p.Brand,
		TaxClass:    p.TaxClass,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Category != nil {
		cat := toCategoryResponse(p.Category)
		resp.Category = &cat
	}
	return resp
}

func toCategoryResponse(c *models.Category) dto.CategoryResponse {
	resp := dto.CategoryResponse{
		ID:          c.ID,
		StoreID:     c.StoreID,
		ParentID:    c.ParentID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Visibility:  c.Visibility,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
	resp.Children = make([]dto.CategoryResponse, len(c.Children))
	for i, ch := range c.Children {
		resp.Children[i] = toCategoryResponse(&ch)
	}
	return resp
}

func toCollectionResponse(col *models.Collection) dto.CollectionResponse {
	return dto.CollectionResponse{
		ID:        col.ID,
		StoreID:   col.StoreID,
		Name:      col.Name,
		Slug:      col.Slug,
		Type:      col.Type,
		Rule:      col.Rule,
		CreatedAt: col.CreatedAt,
		UpdatedAt: col.UpdatedAt,
	}
}
