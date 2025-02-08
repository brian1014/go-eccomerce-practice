package handlers

import (
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CatalogHandler struct {
	// svc UserService
	svc service.CatalogService
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	// create an instace of userService & inject to handler
	svc := service.CatalogService{
		Repo:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}
	handler := CatalogHandler{
		svc: svc,
	}

	// Public endpoints
	// listing products and categories
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProduct)
	app.Get("/categories", handler.GetCategories)
	app.Get("/categories/:id", handler.GetCategoryById)

	// Private endpoints
	// manage products and categories

	sellerRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)

	// Categories
	sellerRoutes.Post("/categories", handler.CreateCategories)
	sellerRoutes.Patch("/categories/:id", handler.EditCategory)
	sellerRoutes.Delete("/categories/:id", handler.DeleteCategory)

	// Products
	sellerRoutes.Post("/products", handler.CreateProducts)
	sellerRoutes.Get("/products", handler.GetProducts)
	sellerRoutes.Get("/products/:id", handler.GetProduct)
	sellerRoutes.Put("/products/:id", handler.EditProduct)
	sellerRoutes.Patch("/products/:id", handler.UpdateProductStock) // Update stock
	sellerRoutes.Delete("/products/:id", handler.DeleteProduct)
}

func (h CatalogHandler) GetCategories(ctx *fiber.Ctx) error {
	cats, err := h.svc.GetCategories()

	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}
	return rest.SuccessResponse(ctx, "categories", cats)
}

func (h CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	cat, err := h.svc.GetCategory(id)
	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}

	return rest.SuccessResponse(ctx, "category", cat)
}

func (h CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequest{}

	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestError(ctx, "create category request is not valid")
	}

	err = h.svc.CreateCategory(req)

	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "category created succesfully", nil)
}

func (h CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequest{}
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestError(ctx, "update category request is not valid")
	}

	updateCat, err := h.svc.EditCategory(id, req)

	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "update category succesfully", updateCat)
}

func (h CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	err := h.svc.DeleteCategory(id)

	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "category deleted succesfully", nil)
}

func (h CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {
	req := dto.CreateProductRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "create product request is not valid")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)
	err = h.svc.CreateProduct(req, user)
	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "product created succesfully", nil)
}

func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	products, err := h.svc.GetProducts()
	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}

	return rest.SuccessResponse(ctx, "products", products)
}

func (h CatalogHandler) GetProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	product, err := h.svc.GetProductById(id)

	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusNotFound, err)
	}

	return rest.SuccessResponse(ctx, "product", product)
}

func (h CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	req := dto.CreateProductRequest{}
	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestError(ctx, "provide a valid input")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	product, err := h.svc.EditProduct(id, req, user)

	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "edit product", product)
}

func (h CatalogHandler) UpdateProductStock(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	req := dto.UpdateStockRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "update stock request is not valid")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	product := domain.Product{
		ID:     uint(id),
		Stock:  uint(req.Stock),
		UserId: int(user.ID),
	}

	updatedProduct, err := h.svc.UpdateProductStock(product)

	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "update stock", updatedProduct)
}

func (h CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	user := h.svc.Auth.GetCurrentUser(ctx)

	err := h.svc.DeleteProduct(id, user)
	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "delete product succesfully", nil)
}
