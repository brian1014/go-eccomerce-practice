package handlers

import (
	"encoding/json"
	"errors"
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"go-ecommerce-app/pkg/payment"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionHanlder struct {
	svc           service.TransactionService
	userSvc       service.UserService
	paymentClient payment.PaymentClient
	config        config.AppConfig
}

func initializeTransactionService(db *gorm.DB, auth helper.Auth) service.TransactionService {
	return service.TransactionService{
		Repo: repository.NewTransactionRepository(db),
		Auth: auth,
	}
}

func SetupTransactionRoutes(rh *rest.RestHandler) {
	app := rh.App
	svc := initializeTransactionService(rh.DB, rh.Auth)

	userSvc := service.UserService{
		Repo:   repository.NewUserRepository(rh.DB),
		CRepo:  repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}

	handler := TransactionHanlder{
		svc:           svc,
		paymentClient: rh.Pc,
		userSvc:       userSvc,
		config:        rh.Config,
	}

	secRoute := app.Group("/buyer", rh.Auth.Authorize)
	secRoute.Get("/payment", handler.MakePayment)
	secRoute.Get("/verify", handler.VerifyPayment)

	sellerRoute := app.Group("/seller", rh.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHanlder) MakePayment(ctx *fiber.Ctx) error {
	// grab authorized user
	user := h.svc.Auth.GetCurrentUser(ctx)
	pubKey := h.config.StripePublicKey

	// 2. check if payment session active
	activePayment, _ := h.svc.GetActivePayment(int(user.ID))
	if activePayment.ID > 0 {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "create payment",
			"pubKey":  pubKey,
			"secret":  activePayment.ClientSecret,
		})
	}

	// 3. call user service and get the cart data to agregate the total amount and collect payment
	_, amount, _ := h.userSvc.FindCart(user.ID)

	orderId, err := helper.RandomNumbers(8)
	if err != nil {
		return rest.InternalErrorMessage(ctx, errors.New("error genereting id order"))
	}

	// 4. Create new payment session
	paymentResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderId)
	if err != nil {
		return rest.ErrorMessage(ctx, http.StatusBadRequest, err)
	}

	// 5. Store the payment in Database to create and validate order
	_ = h.svc.StoreCreatedPayment(dto.CreatePaymentRequest{
		UserId:       user.ID,
		Amount:       amount,
		ClientSecret: paymentResult.ClientSecret,
		PaymentId:    paymentResult.ID,
		OrderId:      orderId,
	})
	// if err != nil {
	// 	return rest.ErrorMessage(ctx, http.StatusInternalServerError, err)
	// }

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "create payment",
		"pubKey":  pubKey,
		"secret":  paymentResult.ClientSecret,
	})
}

func (h *TransactionHanlder) VerifyPayment(ctx *fiber.Ctx) error {
	// grab authorized user
	user := h.svc.Auth.GetCurrentUser(ctx)

	// do we have a session active to verify?
	activePayment, err := h.svc.GetActivePayment(int(user.ID))
	if err != nil || activePayment.ID == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(errors.New("no active payment exists"))
	}

	// fetch payment status from stripe
	paymentRes, _ := h.paymentClient.GetPaymentStatus(activePayment.PaymentId)
	paymentJson, _ := json.Marshal(paymentRes)
	paymentLogs := string(paymentJson)
	paymentStatus := "failed"

	//if payment then create order
	if paymentRes.Status == "succeeded" {
		// create Order
		paymentStatus = "success"
		err = h.userSvc.CreateOrder(user.ID, activePayment.OrderId, activePayment.PaymentId, activePayment.Amount)
	}

	if err != nil {
		return rest.InternalErrorMessage(ctx, err)
	}

	// update payment status
	h.svc.UpdatePayment(user.ID, paymentStatus, paymentLogs)

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":  "create payment",
		"response": paymentRes,
	})
}

func (h *TransactionHanlder) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}

func (h *TransactionHanlder) GetOrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}
