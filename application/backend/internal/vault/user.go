package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/7d4b9/utrade/backend/internal/http/api"
	"github.com/7d4b9/utrade/backend/vault"
	"github.com/spf13/viper"
)

var (
	ErrRecoveryKeyNotFound        = errors.New("recovery key not found")
	ErrRecoveryKeyVersionNotFound = errors.New("recovery key version not found")
	ErrRotationTimestampNotFound  = errors.New("rotation timestamp not found")
)

var config = viper.New()

const (
	addrConfig                               = "addr"
	maxUserRecoveryKeysConfig                = "max_user_recovery_xorkeys"
	minUserRecoveryKeyRotationIntervalConfig = "min_user_recovery_key_rotation_interval"
)

func init() {
	config.AutomaticEnv()
	config.SetEnvPrefix("vault")
	config.SetDefault(minUserRecoveryKeyRotationIntervalConfig, "24h")
	config.SetDefault(maxUserRecoveryKeysConfig, "5")
}

// UserStorage is backend storage for secured userData
type Vault interface {
	DecryptUserRecoveryKey(encryptedUserRecoveryKey []byte) ([]byte, error)
	EncryptUserRecoveryKey(userRecoveryKey []byte) ([]byte, error)
}

// UserStorage is backend storage for secured userData
type UserStorage interface {
	// StoreEncryptedUserRecoveryKey stores the encrypted recovery key at a specific version.
	// Returns an error if storage fails (500).
	StoreEncryptedUserRecoveryKey(userID string, version int, encryptedUserRecoveryKey []byte) error

	// RetrieveEncryptedUserRecoveryKey returns the encrypted recovery key of a given version.
	// Returns (nil, ErrRecoveryKeyNotFound) if not found (404).
	RetrieveEncryptedUserRecoveryKey(userID string, version int) ([]byte, error)

	// RetrieveLatestRecoveryKeyVersion returns the latest version number for the user's recovery key.
	// Returns (0, ErrRecoveryKeyVersionNotFound) if none exists (404).
	RetrieveLatestRecoveryKeyVersion(userID string) (int, error)

	// StoreLatestRecoveryKeyVersion sets the given version as the latest one.
	// Returns an error if it fails to persist (500).
	StoreLatestRecoveryKeyVersion(userID string, version int) error

	// RetrieveLastRotationTimestamp returns the last rotation timestamp.
	// Returns (zero time, ErrRotationTimestampNotFound) if not set yet.
	RetrieveLastRotationTimestamp(userID string) (time.Time, error)

	// StoreLastRotationTimestamp stores the last rotation timestamp.
	StoreLastRotationTimestamp(userID string, timestamp time.Time) error

	// DeleteOldRecoveryKeyVersions deletes old versions and keeps only the latest N.
	// Returns an error if the deletion fails.
	DeleteOldRecoveryKeyVersions(userID string, keepLatestN int) error
}

type Client struct {
	apiClient Vault

	userStorage UserStorage

	maxRecoveryKeysPerUser int

	minRecoveryKeyRotationInterval time.Duration
}

func NewClient(ctx context.Context, userStorage UserStorage) (*Client, error) {

	minUserRecoveryKeyRotationInterval := config.GetDuration(minUserRecoveryKeyRotationIntervalConfig)
	maxRecoveryKeysPerUser := config.GetInt(maxUserRecoveryKeysConfig)

	vaultAPIclient, err := vault.NewAPIclient(ctx)
	if err != nil {
		return nil, fmt.Errorf("new Vault vault: %w", err)
	}
	return &Client{
		apiClient:   vaultAPIclient,
		userStorage: userStorage,

		maxRecoveryKeysPerUser:         maxRecoveryKeysPerUser,
		minRecoveryKeyRotationInterval: minUserRecoveryKeyRotationInterval,
	}, nil
}

// StoreUserRecoveryKey encrypts and stores a new user recovery key with versioning and rotation limits.
func (c *Client) StoreUserRecoveryKey(userID string, recoveryKey []byte) error {

	lastRotation, err := c.userStorage.RetrieveLastRotationTimestamp(userID)
	if err != nil && !errors.Is(err, ErrRotationTimestampNotFound) {
		return fmt.Errorf("failed to retrieve last rotation timestamp: %w", api.ErrRecoveryKeyNotFound)
	}
	if err == nil && time.Since(lastRotation) < c.minRecoveryKeyRotationInterval {
		return fmt.Errorf("recovery key rotation too frequent for user %s", userID)
	}

	latestVersion, err := c.userStorage.RetrieveLatestRecoveryKeyVersion(userID)
	if err != nil && !errors.Is(err, ErrRecoveryKeyVersionNotFound) {
		return fmt.Errorf("failed to retrieve latest recovery key version: %w", err)
	}
	if errors.Is(err, ErrRecoveryKeyVersionNotFound) {
		latestVersion = 0
	}
	newVersion := latestVersion + 1

	encryptedSecret, err := c.apiClient.EncryptUserRecoveryKey(recoveryKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt recoveryKey with Vault Transit: %w", err)
	}

	if err := c.userStorage.StoreEncryptedUserRecoveryKey(userID, newVersion, encryptedSecret); err != nil {
		return fmt.Errorf("failed to store encrypted recovery key: %w", err)
	}

	if err := c.userStorage.StoreLatestRecoveryKeyVersion(userID, newVersion); err != nil {
		return fmt.Errorf("failed to update latest version: %w", err)
	}

	if err := c.userStorage.StoreLastRotationTimestamp(userID, time.Now()); err != nil {
		return fmt.Errorf("failed to update last rotation timestamp: %w", err)
	}

	if err := c.userStorage.DeleteOldRecoveryKeyVersions(userID, c.maxRecoveryKeysPerUser); err != nil {
		return fmt.Errorf("failed to delete old recovery keys: %w", err)
	}

	return nil
}

// RetrieveUserRecoveryKey retrieves and decrypts the latest user recovery key.
func (c *Client) RetrieveUserRecoveryKey(userID string) ([]byte, error) {

	latestVersion, err := c.userStorage.RetrieveLatestRecoveryKeyVersion(userID)
	if err != nil {
		if errors.Is(err, ErrRecoveryKeyVersionNotFound) {
			return nil, api.ErrRecoveryKeyNotFound
		}
		return nil, fmt.Errorf("failed to retrieve latest version: %w", err)
	}

	encryptedRecoveryKey, err := c.userStorage.RetrieveEncryptedUserRecoveryKey(userID, latestVersion)
	if err != nil {
		if errors.Is(err, ErrRecoveryKeyNotFound) {
			return nil, api.ErrRecoveryKeyNotFound
		}
		return nil, fmt.Errorf("failed to retrieve encrypted recovery key: %w", err)
	}

	recoveryKey, err := c.apiClient.DecryptUserRecoveryKey(encryptedRecoveryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt recoveryKey with Vault Transit: %w", err)
	}

	return recoveryKey, nil
}

func (c *Client) GetUserRecoveryKeyVersion(userID string) (int, error) {
	version, err := c.userStorage.RetrieveLatestRecoveryKeyVersion(userID)
	if err != nil {
		if errors.Is(err, ErrRecoveryKeyVersionNotFound) {
			return 0, api.ErrRecoveryKeyNotFound
		}
		return 0, fmt.Errorf("failed to retrieve latest recovery key version: %w", err)
	}
	return version, nil
}

func (c *Client) StoreUserRecoveryKeyVersion(userID string, version int) error {
	return c.userStorage.StoreLatestRecoveryKeyVersion(userID, version)
}
