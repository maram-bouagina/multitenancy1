package handlers

import (
	"github.com/gofiber/fiber/v2"

	"multitenancypfe/internal/auth/dto"
	"multitenancypfe/internal/auth/services"
	"multitenancypfe/internal/helpers"
)

// ─── CRUD Handler ─────────────────────────────────────────────────────────────

type TenantHandler struct {
	svc services.TenantService
}

func NewTenantHandler(svc services.TenantService) *TenantHandler {
	return &TenantHandler{svc: svc}
}

func (h *TenantHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateTenantRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Create(req)
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *TenantHandler) GetAll(c *fiber.Ctx) error {
	tenants, err := h.svc.GetAll()
	if err != nil {
		return helpers.Fail(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(tenants)
}

func (h *TenantHandler) GetByID(c *fiber.Ctx) error {
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	tenant, err := h.svc.GetByID(id)
	if err != nil {
		return helpers.Fail(c, fiber.StatusNotFound, err)
	}
	return c.JSON(tenant)
}

func (h *TenantHandler) Update(c *fiber.Ctx) error {
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	var req dto.UpdateTenantRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Update(id, req)
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	return c.JSON(resp)
}

func (h *TenantHandler) Delete(c *fiber.Ctx) error {
	id, err := helpers.ParseID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(id); err != nil {
		return helpers.Fail(c, fiber.StatusNotFound, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ─── Auth Handler ─────────────────────────────────────────────────────────────

type TenantAuthHandler struct {
	svc services.TenantAuthService
}

func NewTenantAuthHandler(svc services.TenantAuthService) *TenantAuthHandler {
	return &TenantAuthHandler{svc: svc}
}

// POST /api/auth/tenant/login
func (h *TenantAuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}
	resp, err := h.svc.Login(req)
	if err != nil {
		return helpers.Fail(c, fiber.StatusUnauthorized, err)
	}
	return c.JSON(resp)
}
