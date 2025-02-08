package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort            string
	Dsn                   string
	AppSecret             string
	TwilioAccountSID      string
	TwilioAuthToken       string
	TwilioFromPhoneNumber string
	StripeSecretKey       string
	StripePublicKey       string
}

func SetupEnv() (cfg AppConfig, err error) {
	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}
	// godotenv.Load()

	httpPort := os.Getenv("SERVER_PORT")

	Dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	appSecret := os.Getenv("APP_SECRET")

	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")

	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")

	twilioFromPhoneNumber := os.Getenv("TWILIO_FROM_PHONE_NUMBER")

	stripePublicKey := os.Getenv("STRIPE_PUBLIC_KEY")

	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")

	return AppConfig{
		ServerPort:            httpPort,
		Dsn:                   Dsn,
		AppSecret:             appSecret,
		TwilioAccountSID:      twilioAccountSID,
		TwilioAuthToken:       twilioAuthToken,
		TwilioFromPhoneNumber: twilioFromPhoneNumber,
		StripePublicKey:       stripePublicKey,
		StripeSecretKey:       stripeSecretKey,
	}, nil
}
