package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string
type ProductVisibility string

const (
	StatusDraft     ProductStatus = "draft"
	StatusPublished ProductStatus = "published"
	StatusArchived  ProductStatus = "archived"

	VisibilityPublic  ProductVisibility = "public"
	VisibilityPrivate ProductVisibility = "private"
)

type Product struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoreID     uuid.UUID         `gorm:"type:uuid;not null;index"                       json:"store_id"`
	CategoryID  *uuid.UUID        `gorm:"type:uuid;index"                                json:"category_id,omitempty"`
	Title       string            `gorm:"type:varchar(255);not null"                     json:"title"`
	Description *string           `gorm:"type:text"                                      json:"description,omitempty"`
	Slug        string            `gorm:"type:varchar(255);not null;uniqueIndex:idx_slug_store" json:"slug"`
	Status      ProductStatus     `gorm:"type:varchar(20);not null;default:'draft'"      json:"status"`
	Visibility  ProductVisibility `gorm:"type:varchar(20);not null;default:'public'"     json:"visibility"`
	Price       float64           `gorm:"type:decimal(12,2);not null;default:0"          json:"price"`
	Currency    string            `gorm:"type:char(3);not null;default:'EUR'"            json:"currency"`
	SKU         *string           `gorm:"type:varchar(100)"                              json:"sku,omitempty"`
	TrackStock  bool              `gorm:"not null;default:false"                         json:"track_stock"`
	Stock       int               `gorm:"not null;default:0"                             json:"stock"`
	Weight      *float64          `gorm:"type:decimal(10,2)"                             json:"weight,omitempty"`
	Dimensions  *string           `gorm:"type:varchar(100)"                              json:"dimensions,omitempty"`
	Brand       *string           `gorm:"type:varchar(255)"                              json:"brand,omitempty"`
	TaxClass    *string           `gorm:"type:varchar(100)"                              json:"tax_class,omitempty"`
	PublishedAt *time.Time        `gorm:"type:timestamptz"                               json:"published_at,omitempty"`
	CreatedAt   time.Time         `gorm:"type:timestamptz;autoCreateTime"                json:"created_at"`
	UpdatedAt   time.Time         `gorm:"type:timestamptz;autoUpdateTime"                json:"updated_at"`
	DeletedAt   *time.Time        `gorm:"type:timestamptz;index"                         json:"-"`

	// Relations
	Category    *Category    `gorm:"foreignKey:CategoryID"          json:"category,omitempty"`
	Collections []Collection `gorm:"many2many:collection_products;" json:"collections,omitempty"`
}
