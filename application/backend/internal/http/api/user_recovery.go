package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type UserRecoveryKeyVault interface {
	// StoreUserRecoveryKey stores a user's recovery key in the vault.
	// Returns an error if the storage fails (500).
	StoreUserRecoveryKey(userID string, recoveryKey []byte) error

	// RetrieveUserRecoveryKey retrieves a user's recovery key.
	// Returns (nil, ErrRecoveryKeyNotFound) if the key does not exist (404).
	RetrieveUserRecoveryKey(userID string) ([]byte, error)

	// GetUserRecoveryKeyVersion returns the latest version of the user's recovery key.
	// Returns (0, ErrRecoveryKeyNotFound) if not found (404).
	GetUserRecoveryKeyVersion(userID string) (int, error)

	// StoreUserRecoveryKeyVersion stores the version for a user's recovery key.
	// Returns an error if the storage fails (500).
	StoreUserRecoveryKeyVersion(userID string, version int) error
}

// UserRecoveryKeyService defines all routes handled by UserRecoveryKeyService
type UserRecoveryKeyService struct {
	vault UserRecoveryKeyVault
}

type StoreUserRecoveryKeyRequestBody struct {
	UserID      string `json:"userId"`
	RecoveryKey []byte `json:"recoveryKey"`
}

func (s *UserRecoveryKeyService) retrieveUserRecoveryKey(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	recoveryKey, err := s.vault.RetrieveUserRecoveryKey(reqBody.UserID)
	if err != nil {
		if errors.Is(err, ErrRecoveryKeyNotFound) {
			http.Error(w, `{"error":"recovery key not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"recoveryKey": string(recoveryKey),
	})
}

func (s *UserRecoveryKeyService) storeUserRecoveryKey(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	var reqBody StoreUserRecoveryKeyRequestBody

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := s.vault.StoreUserRecoveryKey(reqBody.UserID, reqBody.RecoveryKey); err != nil {
		http.Error(w, `{"error":"failed to store UserRecoveryKey in vault"}`, http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to store UserRecoveryKey in vault")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"recoveryKey stored successfully"}`))
}

// getUserRecoveryKeyVersion handles getting the current version of the recovery key for a user.
func (s *UserRecoveryKeyService) getUserRecoveryKeyVersion(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}
	version, err := s.vault.GetUserRecoveryKeyVersion(req.UserID)
	if err != nil {
		if errors.Is(err, ErrRecoveryKeyNotFound) {
			http.Error(w, `{"error":"recovery key version not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"version": version})
}
