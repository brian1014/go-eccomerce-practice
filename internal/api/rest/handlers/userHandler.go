package handlers

import (
	"errors"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	// svc UserService
	svc service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App

	// create an instace of userService & inject to handler
	svc := service.UserService{
		Repo:   repository.NewUserRepository(rh.DB),
		CRepo:  repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}
	handler := UserHandler{
		svc: svc,
	}

	pubRoutes := app.Group("/")

	// Public endpoints
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.Login)

	pvtRoutes := pubRoutes.Group("/users", rh.Auth.Authorize)

	// Private endpoints
	pvtRoutes.Get("/verify", handler.GetVerificationCode)
	pvtRoutes.Post("/verify", handler.Verify)

	pvtRoutes.Post("/profile", handler.CreateProfile)
	pvtRoutes.Get("/profile", handler.GetProfile)
	pvtRoutes.Patch("/profile", handler.UpdateProfile)

	pvtRoutes.Post("/cart", handler.AddToCart)
	pvtRoutes.Get("/cart", handler.GetCart)

	pvtRoutes.Get("/order", handler.GetOrders)
	pvtRoutes.Get("/order/:id", handler.GetOrder)

	pvtRoutes.Post("/become-seller", handler.BecomeSeller)
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	// to create user
	user := dto.UserSignup{}
	err := ctx.BodyParser(&user)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid inputs",
		})
	}

	token, err := h.svc.SignUp(user)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error on singup",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "register",
		"token":   token,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	loginInput := dto.UserLogin{}
	err := ctx.BodyParser(&loginInput)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid inputs",
		})
	}

	token, err := h.svc.Login(loginInput.Email, loginInput.Password)

	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "please provide correct user id password",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "login",
		"token":   token,
	})
}

func (h *UserHandler) Verify(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	// request
	var req dto.VerificationCodeInput
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide a valid input",
		})
	}

	err := h.svc.VerifyCode(user.ID, req.Code)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Verify succesfully",
	})
}

func (h *UserHandler) GetVerificationCode(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	//create verification code and update to user profile in DB
	err := h.svc.GetVerificationCode(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to generate verification code",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "GetVerificationCode",
	})
}

func (h *UserHandler) CreateProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestError(ctx, "please provide a valid input")
	}

	log.Printf("user %v", user)

	// create profile
	err := h.svc.CreateProfile(user.ID, req)
	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "CreateProfile",
	})
}

func (h *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "GetProfile",
		"profile": profile,
	})
}

func (h *UserHandler) UpdateProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestError(ctx, "please provide a valid input")
	}

	err := h.svc.UpdateProfile(user.ID, req)
	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "profile update succesfully",
	})

}

func (h *UserHandler) AddToCart(ctx *fiber.Ctx) error {
	req := dto.CreateCartRequest{}
	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestError(ctx, "provide a valid product and qty")
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	// call user service and perform create cart
	cartItems, err := h.svc.CreateCart(req, user)
	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "cart created succesfully", cartItems)
}

func (h *UserHandler) GetCart(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	cart, _, err := h.svc.FindCart(user.ID)
	if err != nil {
		return rest.InternalErrorMessage(ctx, errors.New("cart does not exists"))
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "GetCart",
		"cart":    cart,
	})
}

func (h *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	orders, err := h.svc.GetOrders(user)
	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "get orders succesfully", orders)
}

func (h *UserHandler) GetOrder(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	idOrder, _ := strconv.Atoi(ctx.Params("id"))

	order, err := h.svc.GetOrderById(uint(idOrder), user.ID)
	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	return rest.SuccessResponse(ctx, "get order succesfully", order)
}

func (h *UserHandler) BecomeSeller(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.SellerInput{}
	err := ctx.BodyParser(&req)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	token, errBeSeller := h.svc.BecomeSeller(user.ID, req)

	if errBeSeller != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": errBeSeller.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Become seller",
		"token":   token,
	})
}
