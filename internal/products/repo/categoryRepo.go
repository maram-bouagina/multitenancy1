package repo

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"multitenancypfe/internal/products/models"
)

type CategoryRepository interface {
	Create(c *models.Category, db *gorm.DB) error
	FindByID(id, storeID uuid.UUID, db *gorm.DB) (*models.Category, error)
	FindRoots(storeID uuid.UUID, db *gorm.DB) ([]models.Category, error)
	Update(c *models.Category, db *gorm.DB) error
	Delete(id, storeID uuid.UUID, db *gorm.DB) error
	HasProducts(id uuid.UUID, db *gorm.DB) (bool, error)
	SlugExists(slug string, storeID uuid.UUID, excludeID *uuid.UUID, db *gorm.DB) (bool, error)
}

type categoryRepository struct{ db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(c *models.Category, db *gorm.DB) error {
	return db.Create(c).Error
}

func (r *categoryRepository) FindByID(id, storeID uuid.UUID, db *gorm.DB) (*models.Category, error) {
	var c models.Category
	err := db.
		Preload("Children").
		Where("id = ? AND store_id = ?", id, storeID).
		First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

// FindRoots returns top-level categories (parent_id IS NULL) with their direct children.
func (r *categoryRepository) FindRoots(storeID uuid.UUID, db *gorm.DB) ([]models.Category, error) {
	var cats []models.Category
	err := db.
		Preload("Children").
		Where("store_id = ? AND parent_id IS NULL", storeID).
		Order("name ASC").
		Find(&cats).Error
	return cats, err
}

func (r *categoryRepository) Update(c *models.Category, db *gorm.DB) error {
	return db.Save(c).Error
}

func (r *categoryRepository) Delete(id, storeID uuid.UUID, db *gorm.DB) error {
	result := db.Where("id = ? AND store_id = ?", id, storeID).Delete(&models.Category{})
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return result.Error
}

func (r *categoryRepository) HasProducts(id uuid.UUID, db *gorm.DB) (bool, error) {
	var count int64
	err := db.Model(&models.Product{}).
		Where("category_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	return count > 0, err
}

func (r *categoryRepository) SlugExists(slug string, storeID uuid.UUID, excludeID *uuid.UUID, db *gorm.DB) (bool, error) {
	var count int64
	query := db.Model(&models.Category{}).
		Where("slug = ? AND store_id = ?", slug, storeID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	return count > 0, query.Count(&count).Error
}
