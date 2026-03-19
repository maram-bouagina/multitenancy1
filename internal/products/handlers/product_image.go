package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"multitenancypfe/internal/helpers"
	"multitenancypfe/internal/middleware"
	"multitenancypfe/internal/products/dto"
	"multitenancypfe/internal/products/repo"
	"multitenancypfe/internal/products/services"
)

type ImageHandler struct {
	validationSvc  services.ImageValidationService
	pricingSvc     services.PricingService
	publicationSvc services.PublicationValidationService
}

func NewImageHandler(
	validationSvc services.ImageValidationService,
	pricingSvc services.PricingService,
	publicationSvc services.PublicationValidationService,
) *ImageHandler {
	return &ImageHandler{
		validationSvc:  validationSvc,
		pricingSvc:     pricingSvc,
		publicationSvc: publicationSvc,
	}
}

// getTenantImageService creates a service with the tenant-scoped database
func (h *ImageHandler) getTenantImageService(c *fiber.Ctx) services.ProductImageService {
	tenantDB := middleware.GetTenantDB(c)
	return services.NewProductImageService(repo.NewProductImageRepository(tenantDB))
}

// POST /api/stores/:storeId/products/:productId/images
func (h *ImageHandler) Create(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	productID, err := uuid.Parse(c.Params("productId"))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	var req dto.CreateProductImageRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}

	// Validate image before creating
	if err := h.validationSvc.ValidateCreateImageRequest(req, 10); err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	svc := h.getTenantImageService(c)
	resp, err := svc.Create(storeID, productID, req)
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GET /api/stores/:storeId/products/:productId/images
func (h *ImageHandler) GetByProductID(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	productID, err := uuid.Parse(c.Params("productId"))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	svc := h.getTenantImageService(c)
	resp, err := svc.GetByProductID(storeID, productID)
	if err != nil {
		return helpers.Fail(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(resp)
}

// PUT /api/stores/:storeId/products/:productId/images/:imageId
func (h *ImageHandler) Update(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	productID, err := uuid.Parse(c.Params("productId"))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	imageID, err := uuid.Parse(c.Params("imageId"))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	var req dto.UpdateProductImageRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}

	svc := h.getTenantImageService(c)
	resp, err := svc.Update(storeID, productID, imageID, req)
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	return c.JSON(resp)
}

// DELETE /api/stores/:storeId/products/:productId/images/:imageId
func (h *ImageHandler) Delete(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	productID, err := uuid.Parse(c.Params("productId"))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}
	imageID, err := uuid.Parse(c.Params("imageId"))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	svc := h.getTenantImageService(c)
	if err := svc.Delete(storeID, productID, imageID); err != nil {
		return helpers.Fail(c, fiber.StatusNotFound, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// POST /api/stores/:storeId/products/:productId/images/reorder
func (h *ImageHandler) Reorder(c *fiber.Ctx) error {
	storeID, err := parseStoreID(c)
	if err != nil {
		return err
	}
	productID, err := uuid.Parse(c.Params("productId"))
	if err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	var req dto.ReorderImagesRequest
	if err := helpers.ParseBody(c, &req); err != nil {
		return err
	}

	// Convert images to position map
	imagePositions := make(map[uuid.UUID]int)
	for _, img := range req.Images {
		imagePositions[img.ID] = img.Position
	}

	svc := h.getTenantImageService(c)
	if err := svc.Reorder(storeID, productID, imagePositions); err != nil {
		return helpers.Fail(c, fiber.StatusBadRequest, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
