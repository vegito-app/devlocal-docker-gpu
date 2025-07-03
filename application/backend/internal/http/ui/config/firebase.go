package config

import (
	"net/http"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/rs/zerolog/log"
)

func Firebase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("new firebase config secretmanager client")
		return
	}
	defer func() {
		err = client.Close()
		if err != nil {
			log.Error().Err(err).Msg("new firebase secretmanager client")
			return
		}
	}()
	firebaseConfigSecretID := config.GetString(firebaseSecretSecretIDConfig)
	firebaseConfigSecretVersionRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: firebaseConfigSecretID,
	}
	firebaseConfigSecretVersion, err := client.AccessSecretVersion(ctx, firebaseConfigSecretVersionRequest)
	if err != nil {
		log.Error().Err(err).Msg("firebase secret version access")
		return
	}
	firebaseConfigSecret := firebaseConfigSecretVersion.GetPayload()
	if firebaseConfigSecret == nil {
		log.Error().Err(err).Msg("firebase config get payload from secret version")
		return
	}
	firebaseConfigJSON := firebaseConfigSecret.Data
	_, err = w.Write(firebaseConfigJSON)
	if err != nil {
		log.Error().Err(err).Msg("write firebase decoded secret JSON returned payload")
		return
	}
}
