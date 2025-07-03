package api

import (
	"context"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/7d4b9/utrade/backend/firebase"
)

type FirebaseAuthkey string

const firebaseAuthTokenUID FirebaseAuthkey = "firebase_auth_token_uid"
const firebaseAuthIsAnonymous FirebaseAuthkey = "firebase_auth_is_anonymous"

func FirebaseAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || len(auth) < 8 || auth[:7] != "Bearer " {
			http.Error(w, `{"error":"missing or invalid Authorization header"}`, http.StatusUnauthorized)
			return
		}
		idToken := auth[7:]

		ctx := r.Context()
		app, _ := firebase.NewApp(ctx)
		client, _ := app.Auth(ctx)
		token, err := client.VerifyIDToken(ctx, idToken)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		r = setRequiredUserID(r, token)
		next.ServeHTTP(w, r)
	})
}

func authenticatedUserID(r *http.Request) string {
	ctx := r.Context()
	a := ctx.Value(firebaseAuthTokenUID)
	return a.(string)
}

func setRequiredUserID(r *http.Request, token *auth.Token) *http.Request {
	ctx := r.Context()
	checkAnonymousAuth(r, token)
	return r.WithContext(context.WithValue(ctx, firebaseAuthTokenUID, token.UID))
}

func checkAnonymousAuth(r *http.Request, token *auth.Token) *http.Request {
	ctx := r.Context()
	isAnonymous := token.Firebase.SignInProvider == "anonymous"
	return r.WithContext(context.WithValue(ctx, firebaseAuthIsAnonymous, isAnonymous))
}
