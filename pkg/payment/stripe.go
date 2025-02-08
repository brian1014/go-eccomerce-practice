package payment

import (
	"errors"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error)
	GetPaymentStatus(pId string) (*stripe.PaymentIntent, error)
}

type payment struct {
	stripeSecretKey string
	successUrl      string
	cancelUrl       string
}

// CreatePayment implements PaymentClient.
func (p *payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.stripeSecretKey
	amountInCents := amount * 100

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(amountInCents)),
		Currency:           stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	params.AddMetadata("order_id", fmt.Sprintf("%v", orderId))
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))

	pi, err := paymentintent.New(params)

	if err != nil {
		log.Printf("Error creating intent payment: %v", err)
		return nil, errors.New("create payment intent failed")
	}
	log.Println(pi)
	return pi, nil
}

// GetPaymentStatus implements PaymentClient.
func (p *payment) GetPaymentStatus(pId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.stripeSecretKey
	params := &stripe.PaymentIntentParams{}

	result, err := paymentintent.Get(pId, params)

	if err != nil {
		log.Printf("Error getting intent payment: %v", err)
		return nil, errors.New("payment get intent failed")
	}

	return result, nil
}

func NewPaymentClient(stripeSecretKey, succesUrl, cancelUrl string) PaymentClient {
	return &payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      succesUrl,
		cancelUrl:       cancelUrl,
	}
}

// func (p *payment) GetPaymentStatus(pId string) (*stripe.CheckoutSession, error) {
// 	stripe.Key = p.stripeSecretKey
// 	session, err := session.Get(pId, nil)
// 	if err != nil {
// 		log.Printf("Error getting session payment: %v", err)
// 		return nil, errors.New("payment get session failed")
// 	}
// 	return session, nil
// }

// params := &stripe.CheckoutSessionParams{
// 	PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
// 	LineItems: []*stripe.CheckoutSessionLineItemParams{
// 		{
// 			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
// 				UnitAmount: stripe.Int64(int64(amountInCents)),
// 				Currency:   stripe.String("usd"),
// 				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
// 					Name: stripe.String("Electronics"),
// 				},
// 			},
// 			Quantity: stripe.Int64(1),
// 		},
// 	},
// 	Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
// 	SuccessURL: stripe.String(p.successUrl),
// 	CancelURL:  stripe.String(p.cancelUrl),
// }
// params.AddMetadata("order_id", fmt.Sprintf("%v", orderId))
// params.AddMetadata("user_id", fmt.Sprintf("%d", userId))

// session, err := session.New(params)
