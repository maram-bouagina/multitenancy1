package handlers

import (
	"github.com/gofiber/fiber/v2"

	"multitenancypfe/internal/helpers"
	"multitenancypfe/internal/middleware"
	"multitenancypfe/internal/products/dto"
	"multitenancypfe/internal/products/services"
)

type CategoryHandler struct{ svc services.CategoryService }

func NewCategoryHandler(svc services.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// POST /api/stores/:storeId/categories
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	var req dto.CreateCategoryRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Create(storeID, req, middleware.GetTenantDB(c))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GET /api/stores/:storeId/categories — returns full tree
func (h *CategoryHandler) GetTree(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.GetTree(storeID, middleware.GetTenantDB(c))
	if err != nil {
		return helpers.Fail(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(resp)
}

// GET /api/stores/:storeId/categories/:id
func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.GetByID(id, storeID, middleware.GetTenantDB(c))
	if err != nil {
		return helpers.Fail(c, fiber.StatusNotFound, err)
	}
	return c.JSON(resp)
}

// PUT /api/stores/:storeId/categories/:id
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	var req dto.UpdateCategoryRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Update(id, storeID, req, middleware.GetTenantDB(c))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	return c.JSON(resp)
}

// DELETE /api/stores/:storeId/categories/:id
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(id, storeID, middleware.GetTenantDB(c)); err != nil {
		return helpers.Fail(c, fiber.StatusConflict, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
