package dto

import (
	"time"

	"github.com/google/uuid"

	"multitenancypfe/internal/products/models"
)

type CreateProductRequest struct {
	CategoryID  *uuid.UUID               `json:"category_id"`
	Title       string                   `json:"title"        validate:"required,min=1,max=255"`
	Description *string                  `json:"description"`
	Slug        *string                  `json:"slug"`
	Status      models.ProductStatus     `json:"status"       validate:"required,oneof=draft published archived"`
	Visibility  models.ProductVisibility `json:"visibility"   validate:"required,oneof=public private"`
	Price       float64                  `json:"price"        validate:"min=0"`
	Currency    string                   `json:"currency"     validate:"required,len=3"`
	SKU         *string                  `json:"sku"`
	TrackStock  bool                     `json:"track_stock"`
	Stock       int                      `json:"stock"        validate:"min=0"`
	Weight      *float64                 `json:"weight"       validate:"omitempty,min=0"`
	Dimensions  *string                  `json:"dimensions"`
	Brand       *string                  `json:"brand"`
	TaxClass    *string                  `json:"tax_class"`
	PublishedAt *time.Time               `json:"published_at"`
}

type UpdateProductRequest struct {
	CategoryID  *uuid.UUID                `json:"category_id"`
	Title       *string                   `json:"title"        validate:"omitempty,min=1,max=255"`
	Description *string                   `json:"description"`
	Slug        *string                   `json:"slug"`
	Status      *models.ProductStatus     `json:"status"       validate:"omitempty,oneof=draft published archived"`
	Visibility  *models.ProductVisibility `json:"visibility"   validate:"omitempty,oneof=public private"`
	Price       *float64                  `json:"price"        validate:"omitempty,min=0"`
	Currency    *string                   `json:"currency"     validate:"omitempty,len=3"`
	SKU         *string                   `json:"sku"`
	TrackStock  *bool                     `json:"track_stock"`
	Stock       *int                      `json:"stock"        validate:"omitempty,min=0"`
	Weight      *float64                  `json:"weight"       validate:"omitempty,min=0"`
	Dimensions  *string                   `json:"dimensions"`
	Brand       *string                   `json:"brand"`
	TaxClass    *string                   `json:"tax_class"`
	PublishedAt *time.Time                `json:"published_at"`
}

type ProductFilter struct {
	CategoryID *uuid.UUID                `query:"category_id"`
	Status     *models.ProductStatus     `query:"status"`
	Visibility *models.ProductVisibility `query:"visibility"`
	Brand      *string                   `query:"brand"`
	Search     *string                   `query:"search"`
	Page       int                       `query:"page"`
	Limit      int                       `query:"limit"`
}

type ProductResponse struct {
	ID          uuid.UUID                `json:"id"`
	StoreID     uuid.UUID                `json:"store_id"`
	CategoryID  *uuid.UUID               `json:"category_id,omitempty"`
	Title       string                   `json:"title"`
	Description *string                  `json:"description,omitempty"`
	Slug        string                   `json:"slug"`
	Status      models.ProductStatus     `json:"status"`
	Visibility  models.ProductVisibility `json:"visibility"`
	Price       float64                  `json:"price"`
	Currency    string                   `json:"currency"`
	SKU         *string                  `json:"sku,omitempty"`
	TrackStock  bool                     `json:"track_stock"`
	Stock       int                      `json:"stock"`
	Weight      *float64                 `json:"weight,omitempty"`
	Dimensions  *string                  `json:"dimensions,omitempty"`
	Brand       *string                  `json:"brand,omitempty"`
	TaxClass    *string                  `json:"tax_class,omitempty"`
	PublishedAt *time.Time               `json:"published_at,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	Category    *CategoryResponse        `json:"category,omitempty"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}
