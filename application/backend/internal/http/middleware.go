package http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Type middleware : une fonction qui prend un http.Handler et retourne un http.Handler
type Middleware func(http.Handler) http.Handler

// Fonction qui enchaîne les middlewares sur un handler final
func ApplyMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	// On les applique dans l'ordre inverse pour respecter l'ordre d'exécution naturel
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func AuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type requestTimeContextKey string
		const requestTimeKey requestTimeContextKey = "requestTime"
		requestTime := r.Context().Value(requestTimeKey)
		if requestTime == nil {
			// Si le contexte ne contient pas de valeur pour "requestTime", on peut définir une valeur par défaut
			requestTime = r.Context().Value("requestUnixTime")
		}
		// On peut aussi vérifier si requestTime est de type time.Time ou int64 selon le contexte
		if t, ok := requestTime.(time.Time); ok {
			// Si c'est un time.Time, on peut l'utiliser directement
			r = r.WithContext(context.WithValue(r.Context(), requestTimeKey, t.Unix()))
		} else if t, ok := requestTime.(int64); ok {
			// Si c'est un int64, on peut l'utiliser directement
			r = r.WithContext(context.WithValue(r.Context(), requestTimeKey, t))
		} else {
			// Si ce n'est ni l'un ni l'autre, on peut définir une valeur par défaut
			r = r.WithContext(context.WithValue(r.Context(), requestTimeKey, time.Now().Unix()))
		}
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read request body")
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		r = setRequestBodyInContext(r, requestBody)
		// Logique d'audit ici, par exemple enregistrer l'URL et la méthode
		auditedFields := map[string]any{
			"method":      r.Method,
			"url":         r.URL.Path,
			"requestTime": requestTime,
			"requestBody": string(requestBody),
		}
		requestStartTime := time.Now()
		defer func() {
			auditedFields["duration"] = time.Since(requestStartTime)
			log.Info().Fields(auditedFields).Msg("Request received")
		}()
		// Appel du handler suivant
		next.ServeHTTP(w, r)
	})
}

const requestBodyKey = "requestBody"

func setRequestBodyInContext(r *http.Request, body []byte) *http.Request {
	ctx := r.Context()
	return r.WithContext(context.WithValue(ctx, requestBodyKey, body))
}

func RequestBodyFromContext(r *http.Request) ([]byte, bool) {
	ctx := r.Context()
	body, ok := ctx.Value(requestBodyKey).([]byte)
	if !ok {
		return nil, false
	}
	return body, true
}
