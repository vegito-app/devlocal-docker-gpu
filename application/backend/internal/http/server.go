package http

import (
	"context"
	"fmt"
	"net/http"

	backendhttp "github.com/7d4b9/utrade/backend/http"
	"github.com/spf13/viper"
)

var config = viper.New()

const (
	portConfig = "port"
)

func init() {
	config.AutomaticEnv()
	config.SetDefault(portConfig, "8080")
}

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Autoriser les requêtes depuis n'importe quelle origine (en dev)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Autoriser les méthodes HTTP spécifiques
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Autoriser les en-têtes spécifiques
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Gérer la pré-vérification des requêtes OPTIONS (préflight requests)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Passer à l'étape suivante
		next.ServeHTTP(w, r)
	})
}

// StartAPI creates a new instance of apiv1.Service.
func StartAPI(ctx context.Context, mux http.Handler) error {
	port := config.GetString(portConfig)

	if err := backendhttp.ListenAndServe(ctx, "0.0.0.0:"+port, ApplyMiddleware(corsMiddleware(mux), AuditMiddleware)); err != nil {
		return fmt.Errorf("HTTP listenAndServe: %w", err)
	}
	fmt.Print("See you the next time ! Bye")
	return nil
}
