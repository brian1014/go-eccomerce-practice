package api

import (
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/handlers"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/pkg/payment"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Database connection error %v", err)
	}

	log.Println("Database connected!!!")

	// run migrations
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Address{},
		&domain.BackAccount{},
		&domain.Category{},
		&domain.Product{},
		&domain.Cart{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	)
	if err != nil {
		log.Fatalf("error on runing migration %v", err.Error())
	}
	log.Println("Migrations succesfully!!!")

	// cors configuration
	cors := cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	})

	app.Use(cors)

	app.Get("/", func(c *fiber.Ctx) error {
		return rest.SuccessResponse(c, "Healthy", &fiber.Map{
			"status": "ok",
		})
	})

	auth := helper.SetupAuth(config.AppSecret)
	paymentClient := payment.NewPaymentClient(config.StripeSecretKey)

	restHandler := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
		Pc:     paymentClient,
	}

	setupRoutes(restHandler)

	app.Listen(config.ServerPort)
}

func setupRoutes(rh *rest.RestHandler) {
	// catalog handler
	handlers.SetupCatalogRoutes(rh)

	// user handler
	handlers.SetupUserRoutes(rh)

	// transaction handler
	handlers.SetupTransactionRoutes(rh)

}
