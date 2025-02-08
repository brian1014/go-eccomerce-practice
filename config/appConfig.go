package config

import (
	"errors"
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
	StripeSuccessUrl      string
	StripeCancelUrl       string
}

func SetupEnv() (cfg AppConfig, err error) {
	// if os.Getenv("APP_ENV") == "dev" {
	// 	godotenv.Load()
	// }
	godotenv.Load()

	httpPort := os.Getenv("HTTP_PORT")
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("env variable not found")
	}

	Dsn := os.Getenv("DSN")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("env variable not found")
	}

	appSecret := os.Getenv("APP_SECRET")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("APP_SECRET variable not found")
	}

	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("TWILIO_ACCOUNT_SID variable not found")
	}

	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("TWILIO_AUTH_TOKEN variable not found")
	}

	twilioFromPhoneNumber := os.Getenv("TWILIO_FROM_PHONE_NUMBER")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("TWILIO_FROM_PHONE_NUMBER variable not found")
	}

	stripePublicKey := os.Getenv("STRIPE_PUBLIC_KEY")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("STRIPE_PUBLIC_KEY variable not found")
	}

	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("STRIPE_SECRET_KEY variable not found")
	}

	stripeSuccessUrl := os.Getenv("STRIPE_SUCCESS_URL")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("STRIPE_SUCCESS_URL variable not found")
	}

	stripeCancelUrl := os.Getenv("STRIPE_CANCEL_URL")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("STRIPE_CANCEL_URL variable not found")
	}

	return AppConfig{
		ServerPort:            httpPort,
		Dsn:                   Dsn,
		AppSecret:             appSecret,
		TwilioAccountSID:      twilioAccountSID,
		TwilioAuthToken:       twilioAuthToken,
		TwilioFromPhoneNumber: twilioFromPhoneNumber,
		StripePublicKey:       stripePublicKey,
		StripeSecretKey:       stripeSecretKey,
		StripeSuccessUrl:      stripeSuccessUrl,
		StripeCancelUrl:       stripeCancelUrl,
	}, nil
}
