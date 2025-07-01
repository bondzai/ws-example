package handlers

import (
	"api-gateway/internal/entities"
	"api-gateway/internal/usecases"
	"api-gateway/pkg/errs"

	"github.com/gofiber/fiber/v2"
)

type (
	RestHandler interface {
		HealthCheck(c *fiber.Ctx) error

		CreateMerchant(c *fiber.Ctx) error
		GetMerchants(c *fiber.Ctx) error
		GetMerchantById(c *fiber.Ctx) error
		UpdateMerchant(c *fiber.Ctx) error
		DeleteMerchant(c *fiber.Ctx) error

		CreateCryptoChain(c *fiber.Ctx) error
		GetCryptoChains(c *fiber.Ctx) error
		GetCryptoChainById(c *fiber.Ctx) error
		UpdateCryptoChain(c *fiber.Ctx) error

		CreateMerchantCategory(c *fiber.Ctx) error
		GetMerchantCategories(c *fiber.Ctx) error
		GetMerchantCategoryById(c *fiber.Ctx) error
		UpdateMerchantCategory(c *fiber.Ctx) error
		DeleteMerchantCategory(c *fiber.Ctx) error
	}

	restHandler struct {
		merchantUseCase         usecases.MerchantUseCase
		merchantCategoryUseCase usecases.MerchantCategoryUseCase
		cryptoChainUseCase      usecases.CryptoChainUseCase
	}
)

func NewRestHandler(
	merchantUseCase usecases.MerchantUseCase,
	merchantCategoryUseCase usecases.MerchantCategoryUseCase,
	cryptoChainUseCase usecases.CryptoChainUseCase,
) RestHandler {
	return &restHandler{
		merchantUseCase:         merchantUseCase,
		merchantCategoryUseCase: merchantCategoryUseCase,
		cryptoChainUseCase:      cryptoChainUseCase,
	}
}

func (h *restHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Server is running")
}

func (h *restHandler) CreateMerchant(c *fiber.Ctx) error {
	var req entities.MerchantCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.HandleFiberError(c, err)
	}

	res, err := h.merchantUseCase.CreateMerchant(req)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *restHandler) GetMerchants(c *fiber.Ctx) error {
	res, err := h.merchantUseCase.GetMerchants()
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) GetMerchantById(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.merchantUseCase.GetMerchantById(id)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) UpdateMerchant(c *fiber.Ctx) error {
	id := c.Params("id")
	var req entities.MerchantUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.HandleFiberError(c, err)
	}

	res, err := h.merchantUseCase.UpdateMerchant(id, req)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) DeleteMerchant(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.merchantUseCase.DeleteMerchant(id)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *restHandler) CreateCryptoChain(c *fiber.Ctx) error {
	var req entities.CryptoChainCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.HandleFiberError(c, err)
	}

	res, err := h.cryptoChainUseCase.CreateCryptoChain(req)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *restHandler) GetCryptoChains(c *fiber.Ctx) error {
	res, err := h.cryptoChainUseCase.GetCryptoChains()
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) GetCryptoChainById(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.cryptoChainUseCase.GetCryptoChainById(id)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) UpdateCryptoChain(c *fiber.Ctx) error {
	id := c.Params("id")
	var req entities.CryptoChainUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.HandleFiberError(c, err)
	}

	res, err := h.cryptoChainUseCase.UpdateCryptoChain(id, &req)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) CreateMerchantCategory(c *fiber.Ctx) error {
	var req entities.MerchantCategoryCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.HandleFiberError(c, err)
	}

	res, err := h.merchantCategoryUseCase.CreateMerchantCategory(req)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *restHandler) GetMerchantCategories(c *fiber.Ctx) error {
	res, err := h.merchantCategoryUseCase.GetMerchantCategories()
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) GetMerchantCategoryById(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.merchantCategoryUseCase.GetMerchantCategoryById(id)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) UpdateMerchantCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var req entities.MerchantCategoryUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.HandleFiberError(c, err)
	}

	res, err := h.merchantCategoryUseCase.UpdateMerchantCategory(id, req)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.JSON(res)
}

func (h *restHandler) DeleteMerchantCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.merchantCategoryUseCase.DeleteMerchantCategory(id)
	if err != nil {
		return errs.HandleFiberError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
