package utils

import (
	"ekira-backend/app/models"
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/charge"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/refund"
	"math"
)

const STRIPE_COMMISION_PERCENTAGE = 2.9 // 2.9%
const STRIPE_COMMISION_FIXED_TRY = 6.29 // 0.30 USD (Fixed) -> 6.29 TRY

func GetPriceWithCommission(price float64) float64 {
	// calculate net price + stripe commission
	x := (price + STRIPE_COMMISION_FIXED_TRY) / (1 - (STRIPE_COMMISION_PERCENTAGE / 100))
	return math.Round(x*100) / 100
}

func GetReceiptURL(chargeID string) (string, error) {
	chargeInfo, err := charge.Get(chargeID, nil)
	if err != nil {
		return "", err
	}
	return chargeInfo.ReceiptURL, nil
}

func GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func CreateCustomer(user models.User) (*stripe.Customer, error) {
	if user.StripeCustomerID == nil {
		customerInfo, err := customer.New(&stripe.CustomerParams{
			Name:             stripe.String(user.FullName()),
			Email:            stripe.String(user.Email),
			Phone:            stripe.String(user.PhoneNumber),
			Description:      stripe.String(user.ID.String()),
			PreferredLocales: stripe.StringSlice([]string{"tr", "en"}),
		})
		if err != nil {
			return nil, err
		}
		return customerInfo, nil
	} else {
		customerInfo, err := customer.Get(*user.StripeCustomerID, nil)
		if err != nil {
			return nil, err
		}
		if customerInfo.Email != user.Email || customerInfo.Name != user.FullName() || customerInfo.Phone != user.PhoneNumber {
			customerInfo, err = customer.Update(*user.StripeCustomerID, &stripe.CustomerParams{
				Email: stripe.String(user.Email),
				Name:  stripe.String(user.FullName()),
				Phone: stripe.String(user.PhoneNumber),
			})
			if err != nil {
				return nil, err
			}
		}
		return customerInfo, nil
	}
}

func CreateAPaymentIntent(customer *stripe.Customer, description string, amount float64, currency string) (*stripe.PaymentIntent, error) {
	// Create a PaymentIntent with the order amount and currency
	fmt.Println(customer)
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)),
		Currency: stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		ReceiptEmail: stripe.String(customer.Email),
		Customer:     stripe.String(customer.ID),
		Description:  stripe.String(description),
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func RefundCharge(chargeID string) (*stripe.Refund, error) {
	chargeInfo, err := charge.Get(chargeID, nil)
	if err != nil {
		return nil, err
	}
	params := &stripe.RefundParams{
		Charge: stripe.String(chargeInfo.ID),
	}
	refundInfo, err := refund.New(params)
	if err != nil {
		return nil, err
	}
	return refundInfo, nil
}
