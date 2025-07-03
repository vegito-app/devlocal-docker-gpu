package v1

import (
	"encoding/json"
	"net/http"
	nethttp "net/http"

	"github.com/rs/zerolog/log"
)

const (
	paymentWebhookURLconfig = "payment_webhook_url"
)

type CheckoutSession struct {
	// CancelURL          string   `json:"cancel_url"`
	// ClientReferenceID  string   `json:"client_reference_id"`
	// CustomerEmail      string   `json:"customer_email"`
	// Deleted            bool     `json:"deleted"`
	// ID                 string   `json:"id"`
	// Livemode           bool     `json:"livemode"`
	// Locale             string   `json:"locale"`
	// Object             string   `json:"object"`
	// PaymentMethodTypes []string `json:"payment_method_types"`
	// SuccessURL         string   `json:"success_url"`
}

type Payment interface {
	CreateCheckoutSession(priceID string) (CheckoutSession, error)
}

type PaymentService struct {
	payment Payment
}

func NewPaymentService(mux *nethttp.ServeMux, payment Payment) (*PaymentService, error) {
	service := &PaymentService{
		payment: payment,
	}

	mux.Handle("POST /paiement/orders/:orderID/checkout", nethttp.HandlerFunc(service.CreateCheckoutSession))
	mux.Handle("POST /paiement/webhook", nethttp.HandlerFunc(service.HandleStripeWebhook))

	return service, nil
}

func (s *PaymentService) CreateCheckoutSession(w nethttp.ResponseWriter, r *nethttp.Request) {
	// (priceID string) (*stripe.CheckoutSession, error) {
	var payload struct {
		PriceID string `json:"priceId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("create checkout session decode JSON payload")
		http.Error(w, `{"error":"invalid json"}`, nethttp.StatusBadRequest)
		return
	}
	priceID := payload.PriceID
	if priceID == "" {
		log.Error().Msg("create checkout session requires priceId ")
		http.Error(w, `{"error":"priceId is required"}`, nethttp.StatusBadRequest)
		return
	}
	// Create a new checkout session with the provided price ID
	session, err := s.payment.CreateCheckoutSession(priceID)
	if err != nil {
		http.Error(w, "Failed to create checkout session", nethttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(nethttp.StatusOK)
	if err := json.NewEncoder(w).Encode(session); err != nil {
		http.Error(w, "Failed to encode response", nethttp.StatusInternalServerError)
		return
	}
}

func (s *PaymentService) HandleStripeWebhook(w nethttp.ResponseWriter, r *nethttp.Request) {
	// Handle Stripe webhook events here
	// For example, you can verify the event and handle payment success or failure
	w.WriteHeader(nethttp.StatusOK)
}
