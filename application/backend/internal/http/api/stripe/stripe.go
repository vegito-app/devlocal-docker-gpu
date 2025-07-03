package stripe

import (
	"fmt"

	"context" // import "context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	v1 "github.com/7d4b9/utrade/backend/internal/http/api/internal"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

var config = viper.New()

const (
	backendEndpointURLConfig = "backend_endpoint_url"
	stripeKeySecretIDConfig  = "stripe_key_secret_id"
)

func init() {
	config.AutomaticEnv()
	config.SetEnvPrefix("stripe")

}

type Client struct {
}

func NewClient() (*Client, error) {
	stripeKeyGSMversion := config.GetString(stripeKeySecretIDConfig)
	if stripeKeyGSMversion == "" {
		return nil, fmt.Errorf("Stripe secret key is not set in the environment variables")
	}
	stripeSecretKey, err := LoadStripeSecretFromGSM(context.Background(), stripeKeyGSMversion)
	if err != nil {
		return nil, fmt.Errorf("failed to load stripe secret key: %w", err)
	}
	stripe.Key = stripeSecretKey

	return &Client{}, nil
}

func (c *Client) CreateCheckoutSession(priceID string) (*v1.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			// {
			// 	Currency: stripe.String(string(stripe.CurrencyEUR)),
			// 	// Amount:   stripe.Int64(priceID),
			// 	Quantity: stripe.Int64(1),
			// },
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://tonapp.com/success"),
		CancelURL:  stripe.String("https://tonapp.com/cancel"),
	}
	s, err := session.New(params)
	if err != nil {
		return nil, err
	}
	if s.ClientReferenceID == "" {
		return nil, fmt.Errorf("ClientReferenceID is required")
	}
	if s.CustomerEmail == "" {
		return nil, fmt.Errorf("CustomerEmail is required")
	}
	if len(s.PaymentMethodTypes) == 0 {
		return nil, fmt.Errorf("PaymentMethodTypes is required")
	}
	sessionV1 := &v1.CheckoutSession{
		// ID:                 s.ID,
		// CancelURL:          s.CancelURL,
		// ClientReferenceID:  s.ClientReferenceID,
		// CustomerEmail:      s.CustomerEmail,
		// Deleted:            s.Deleted,
		// Livemode:           s.Livemode,
		// Locale:             s.Locale,
		// Object:             s.Object,
		// PaymentMethodTypes: s.PaymentMethodTypes,
		// SuccessURL:         s.SuccessURL,
	}
	return sessionV1, nil
}

func LoadStripeSecretFromGSM(ctx context.Context, projectID string) (string, error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg("new stripe secretmanager client")
		}
	}()

	stripeKeySecretID := config.GetString(stripeKeySecretIDConfig)
	stripeKeyAccessSecretVersionRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: stripeKeySecretID,
	}
	stripeKey, err := client.AccessSecretVersion(ctx, stripeKeyAccessSecretVersionRequest)
	if err != nil {
		return "", fmt.Errorf("access stripe secret version: %w", err)
	}
	return string(stripeKey.Payload.Data), nil
}
