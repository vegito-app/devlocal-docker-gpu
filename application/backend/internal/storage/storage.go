package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/7d4b9/utrade/backend/firebase"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var config = viper.New()

const appCreationTimeoutConfig = "new_application_timeout"

func init() {
	config.AutomaticEnv()
	config.SetEnvPrefix("firebase")
	config.SetDefault(appCreationTimeoutConfig, 1*time.Minute)
}

type Storage struct {
	*VegetableStorage
	*RecoveryKeyStorage
	*OrderStorage
	*UserStorage
	firestore *firestore.Client
	app       *firebase.App
}

func NewStorage(app *firebase.App, opts ...option.ClientOption) (*Storage, error) {

	timeout := config.GetDuration(appCreationTimeoutConfig)
	ctx, cancelAppCreation := context.WithTimeout(context.Background(), timeout)
	defer cancelAppCreation()

	firestore, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase failed to create client: %w", err)
	}
	recoveryKeyStorage := &RecoveryKeyStorage{
		firestore: firestore,
	}
	vegetableStorage := &VegetableStorage{
		firestore: firestore,
	}
	orderStorage := &OrderStorage{
		firestore: firestore,
	}
	userStorage := &UserStorage{
		firestore: firestore,
	}
	return &Storage{
		app:                app,
		firestore:          firestore,
		OrderStorage:       orderStorage,
		RecoveryKeyStorage: recoveryKeyStorage,
		UserStorage:        userStorage,
		VegetableStorage:   vegetableStorage,
	}, nil
}

func (c *Storage) Close() {
	if err := c.firestore.Close(); err != nil {
		log.Println("firebase close firestore, error:", err.Error())
	}
}
