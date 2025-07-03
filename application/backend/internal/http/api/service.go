package api

import (
	"context"
	"fmt"
	nethttp "net/http"

	"github.com/7d4b9/utrade/backend/btc"
	"github.com/7d4b9/utrade/backend/internal/http"
	"github.com/7d4b9/utrade/backend/internal/http/ui"
	uiconfig "github.com/7d4b9/utrade/backend/internal/http/ui/config"
	"github.com/spf13/viper"
)

var config = viper.New()

const (
	uiBuildDirConfig            = "frontend_build_dir"
	isValidatedImagesCDNEnabled = "enable_validated_images_cdn"
)

func init() {
	config.AutomaticEnv()
	config.SetDefault(uiBuildDirConfig, "../frontend/build")
}

// Service defines all routes handled by Service
type Service struct {
	nethttp.Handler
	// authUser is used to verify user authentication
	authUser AuthUser
	// user is used to manage users
	user *UserService
}

type Storage interface {
	UserStorage
}

type AuthUser interface {
	VerifyIDToken(ctx context.Context, idToken string) (string, error)
}

var (
	ErrRecoveryKeyNotFound = fmt.Errorf("recovery key not found")
)

func NewService(authUser AuthUser, storage Storage, btcService *btc.BTC) (*Service, error) {
	mux := nethttp.NewServeMux()

	frontendDir := config.GetString(uiBuildDirConfig)
	uiServe, err := ui.NewUI(frontendDir)
	if err != nil {
		return nil, fmt.Errorf("http start api, server side ui render: %w", err)
	}
	mux.Handle("GET /ui/config/firebase", nethttp.HandlerFunc(uiconfig.Firebase))
	mux.Handle("GET /ui/config/googlemaps", nethttp.HandlerFunc(uiconfig.GoogleMaps))
	mux.Handle("GET /ui/public", nethttp.StripPrefix("/ui", nethttp.FileServer(nethttp.Dir(frontendDir))))
	mux.Handle("GET /ui", uiServe)
	mux.Handle("GET /price", nethttp.HandlerFunc(btcService.Price))
	mux.Handle("GET /cgv", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		fmt.Fprintf(w, "Welcome to the local BackEnd!")
	}))
	mux.Handle("GET /confidentiality", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		fmt.Fprintf(w, "Welcome to the local BackEnd!")
	}))
	mux.Handle("GET /info", nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		fmt.Fprintf(w, "Welcome to the local BackEnd!")
	}))
	mux.Handle("GET /", nethttp.FileServer(nethttp.Dir(frontendDir)))

	userService, err := NewUserService(mux, storage)
	if err != nil {
		return nil, fmt.Errorf("http api v1 new user service: %w", err)
	}
	serviceV1 := &Service{
		user:     userService,
		authUser: authUser,
		Handler:  mux,
	}
	mux.HandleFunc("POST /run", serviceV1.run)

	mux.Handle("GET /api/auth-check", http.ApplyMiddleware(
		nethttp.HandlerFunc(serviceV1.authCheck),
		FirebaseAuthMiddleware))

	return serviceV1, nil
}

func (s *Service) run(w nethttp.ResponseWriter, r *nethttp.Request) {

	w.WriteHeader(nethttp.StatusOK)
}

func (s *Service) Status(w nethttp.ResponseWriter, r *nethttp.Request) {

	w.WriteHeader(nethttp.StatusOK)
}

func (s *Service) authCheck(w nethttp.ResponseWriter, r *nethttp.Request) {

	userID := requestUserID(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"authenticated", "uid":"%s"}`, userID)))
}
