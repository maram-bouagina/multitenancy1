package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"multitenancypfe/internal/helpers"
	"multitenancypfe/internal/products/dto"
	"multitenancypfe/internal/products/services"
)

type ProductHandler struct {
	svc services.ProductService
}

func NewProductHandler(svc services.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// POST /api/stores/:storeId/products
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	var req dto.CreateProductRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Create(storeID, req)
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GET /api/stores/:storeId/products
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	var filter dto.ProductFilter
	if err := c.QueryParser(&filter); err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	resp, err := h.svc.GetAll(storeID, filter)
	if err != nil {
		return helpers.Fail(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(resp)
}

// GET /api/stores/:storeId/products/:id
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.GetByID(id, storeID)
	if err != nil {
		return helpers.Fail(c, fiber.StatusNotFound, err)
	}
	return c.JSON(resp)
}

// PUT /api/stores/:storeId/products/:id
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	var req dto.UpdateProductRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Update(id, storeID, req)
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	return c.JSON(resp)
}

// DELETE /api/stores/:storeId/products/:id
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(id, storeID); err != nil {
		return helpers.Fail(c, fiber.StatusNotFound, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseStoreID(c *fiber.Ctx) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params("storeId"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store id"})
	}
	return id, err
}
