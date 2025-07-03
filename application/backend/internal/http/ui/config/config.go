package config

import "github.com/spf13/viper"

var config = viper.New()

const (
	firebaseSecretSecretIDConfig = "firebase_secret_id"
	firebaseConfig               = "firebase_config_web"
	googlemapsSecretIDConfig     = "googlemaps_secret_id"
	googleMapsAPIkey             = "google_maps_api_key"
)

func init() {
	config.AutomaticEnv()
	config.SetEnvPrefix("ui_config")
}
