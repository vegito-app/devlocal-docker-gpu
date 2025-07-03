package api

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"strings"

	"github.com/7d4b9/utrade/backend/internal/http"
)

var ErrOrderNotFound = errors.New("order not found")

type Order struct {
	ID           string           `json:"id,omitempty"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	SaleType     string           `json:"saleType"`
	WeightGrams  int              `json:"weightGrams"`
	PriceCents   int              `json:"priceCents"`
	Images       []VegetableImage `json:"images"`
	CreatedAt    int64            `json:"createdAt"`
	Status       string           `json:"status"`
	VegetableIds []string         `json:"vegetableIds"`
}

// OrderStorage defines persistent operations for managing orders.
// All implementations must return ErrOrderNotFound if the order is not found.
type OrderStorage interface {
	// StoreOrder stores a new order associated with the given user.
	StoreOrder(ctx context.Context, userID string, o Order) error

	// GetOrder retrieves a specific order by its ID for the given user.
	// Returns ErrOrderNotFound if the order does not exist.
	GetOrder(ctx context.Context, userID, id string) (*Order, error)

	// ListOrders returns all orders visible to the given user.
	ListOrders(ctx context.Context, userID string) ([]*Order, error)

	// DeleteOrder deletes an existing order.
	// Returns ErrOrderNotFound if the order does not exist.
	DeleteOrder(ctx context.Context, userID, id string) error

	// UpdateOrderStatus updates the status of an existing order.
	// Returns ErrOrderNotFound if the order does not exist.
	UpdateOrderStatus(ctx context.Context, userID, orderID, status string) error
}

// OrderService defines all routes handled by OrderService
type OrderService struct {
	storage OrderStorage
}

func NewOrderService(mux *nethttp.ServeMux, storage OrderStorage) (*OrderService, error) {
	service := &OrderService{
		storage: storage,
	}

	mux.Handle("POST /order", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.CreateOrder),
		FirebaseAuthMiddleware))

	mux.Handle("POST /orders/search", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.ListOrdersByVegetableIDs),
		FirebaseAuthMiddleware))

	mux.Handle("GET /orders/client/{id}", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.ListOrdersByClientID),
		FirebaseAuthMiddleware))

	mux.Handle("PUT /orders/", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.UpdateOrderStatus),
		FirebaseAuthMiddleware))

	mux.Handle("GET /orders/", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.GetOrder),
		FirebaseAuthMiddleware))

	mux.Handle("DELETE /orders/", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.DeleteOrder),
		FirebaseAuthMiddleware))

	return service, nil
}

func (s *OrderService) CreateOrder(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := r.Header.Get("X-User-ID")
	var o Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusBadRequest)
		return
	}
	if err := s.storage.StoreOrder(ctx, userID, o); err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(nethttp.StatusCreated)
}

func (s *OrderService) ListOrdersByVegetableIDs(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := r.Header.Get("X-User-ID")

	var payload struct {
		VegetableIDs []string `json:"vegetableIds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		nethttp.Error(w, "invalid JSON payload", nethttp.StatusBadRequest)
		return
	}
	if len(payload.VegetableIDs) == 0 {
		nethttp.Error(w, "vegetableIds list is required", nethttp.StatusBadRequest)
		return
	}

	idSet := make(map[string]struct{})
	for _, id := range payload.VegetableIDs {
		idSet[id] = struct{}{}
	}

	orders, err := s.storage.ListOrders(ctx, userID)
	if err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}

	filtered := make([]*Order, 0, len(orders))
	for _, o := range orders {
		if _, ok := idSet[o.ID]; ok {
			filtered = append(filtered, o)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filtered); err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(nethttp.StatusOK)
}

func (s *OrderService) ListOrdersByClientID(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()

	clientID := strings.TrimPrefix(r.URL.Path, "/orders/client/")
	if clientID == "" {
		nethttp.Error(w, "missing client ID", nethttp.StatusBadRequest)
		return
	}
	userID := authenticatedUserID(r)
	if clientID != userID {
		nethttp.Error(w, "forbidden", nethttp.StatusForbidden)
		return
	}

	orders, err := s.storage.ListOrders(ctx, clientID)
	if err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}

	if orders == nil {
		orders = []*Order{}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(nethttp.StatusOK)
}

func (s *OrderService) UpdateOrderStatus(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")

	if orderID == "" {
		nethttp.Error(w, "missing order ID", nethttp.StatusBadRequest)
		return
	}

	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		nethttp.Error(w, "invalid JSON", nethttp.StatusBadRequest)
		return
	}
	if payload.Status == "" {
		nethttp.Error(w, "missing status", nethttp.StatusBadRequest)
		return
	}

	if err := s.storage.UpdateOrderStatus(ctx, userID, orderID, payload.Status); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			nethttp.Error(w, "order not found", nethttp.StatusNotFound)
			return
		}
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func (s *OrderService) GetOrder(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
	if orderID == "" {
		nethttp.Error(w, "missing order ID", nethttp.StatusBadRequest)
		return
	}

	order, err := s.storage.GetOrder(ctx, userID, orderID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			nethttp.Error(w, "order not found", nethttp.StatusNotFound)
			return
		}
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(nethttp.StatusOK)
}

func (s *OrderService) DeleteOrder(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
	if orderID == "" {
		nethttp.Error(w, "missing order ID", nethttp.StatusBadRequest)
		return
	}

	if err := s.storage.DeleteOrder(ctx, userID, orderID); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			nethttp.Error(w, "order not found", nethttp.StatusNotFound)
			return
		}
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}
