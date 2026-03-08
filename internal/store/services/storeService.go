package services

import (
	"errors"

	"github.com/google/uuid"

	"multitenancypfe/internal/store/dto"
	"multitenancypfe/internal/store/models"
	"multitenancypfe/internal/store/repo"
)

type StoreService interface {
	Create(tenantID uuid.UUID, req dto.CreateStoreRequest) (*dto.StoreResponse, error)
	GetByID(id uuid.UUID) (*dto.StoreResponse, error)
	GetByTenantID(tenantID uuid.UUID) ([]dto.StoreResponse, error)
	Update(id uuid.UUID, req dto.UpdateStoreRequest) (*dto.StoreResponse, error)
	Delete(id uuid.UUID) error
}

type storeService struct {
	repo repo.StoreRepository
}

func NewStoreService(r repo.StoreRepository) StoreService {
	return &storeService{repo: r}
}

func (s *storeService) Create(tenantID uuid.UUID, req dto.CreateStoreRequest) (*dto.StoreResponse, error) {
	existing, err := s.repo.FindBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("slug already in use")
	}

	store := &models.Store{
		TenantID:  tenantID,
		Name:      req.Name,
		Slug:      req.Slug,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Logo:      req.Logo,
		Currency:  req.Currency,
		Timezone:  req.Timezone,
		Language:  req.Language,
		TaxNumber: req.TaxNumber,
		Status:    "active",
	}
	if err := s.repo.Create(store); err != nil {
		return nil, err
	}
	return toStoreResponse(store), nil
}

func (s *storeService) GetByID(id uuid.UUID) (*dto.StoreResponse, error) {
	store, err := s.findOrFail(id)
	if err != nil {
		return nil, err
	}
	return toStoreResponse(store), nil
}

func (s *storeService) GetByTenantID(tenantID uuid.UUID) ([]dto.StoreResponse, error) {
	stores, err := s.repo.FindByTenantID(tenantID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.StoreResponse, len(stores))
	for i, st := range stores {
		result[i] = *toStoreResponse(&st)
	}
	return result, nil
}

func (s *storeService) Update(id uuid.UUID, req dto.UpdateStoreRequest) (*dto.StoreResponse, error) {
	store, err := s.findOrFail(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		store.Name = *req.Name
	}
	if req.Email != nil {
		store.Email = req.Email
	}
	if req.Phone != nil {
		store.Phone = req.Phone
	}
	if req.Address != nil {
		store.Address = req.Address
	}
	if req.Logo != nil {
		store.Logo = req.Logo
	}
	if req.Currency != nil {
		store.Currency = *req.Currency
	}
	if req.Timezone != nil {
		store.Timezone = *req.Timezone
	}
	if req.Language != nil {
		store.Language = *req.Language
	}
	if req.TaxNumber != nil {
		store.TaxNumber = req.TaxNumber
	}
	if req.Status != nil {
		store.Status = *req.Status
	}

	if err := s.repo.Update(store); err != nil {
		return nil, err
	}
	return toStoreResponse(store), nil
}

func (s *storeService) Delete(id uuid.UUID) error {
	if _, err := s.findOrFail(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func (s *storeService) findOrFail(id uuid.UUID) (*models.Store, error) {
	store, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errors.New("store not found")
	}
	return store, nil
}

func toStoreResponse(s *models.Store) *dto.StoreResponse {
	return &dto.StoreResponse{
		ID:        s.ID.String(),
		TenantID:  s.TenantID.String(),
		Name:      s.Name,
		Slug:      s.Slug,
		Email:     s.Email,
		Phone:     s.Phone,
		Address:   s.Address,
		Logo:      s.Logo,
		Currency:  s.Currency,
		Timezone:  s.Timezone,
		Language:  s.Language,
		TaxNumber: s.TaxNumber,
		Status:    s.Status,
	}
}
