package storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/7d4b9/utrade/backend/internal/vault"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecoveryKeyStorage struct {
	firestore *firestore.Client
}

func (r *RecoveryKeyStorage) StoreEncryptedUserRecoveryKey(userID string, version int, encryptedUserRecoveryKey []byte) error {
	ctx := context.Background()
	doc := r.firestore.Collection("users").Doc(userID).Collection("recoveryKeys").Doc(fmt.Sprintf("%d", version))
	_, err := doc.Set(ctx, map[string]any{
		"recoveryKey": encryptedUserRecoveryKey,
	})
	if err != nil {
		return fmt.Errorf("failed to store encrypted user recovery key version %d: %w", version, err)
	}
	return nil
}

func (r *RecoveryKeyStorage) RetrieveEncryptedUserRecoveryKey(userID string, version int) ([]byte, error) {
	ctx := context.Background()
	doc := r.firestore.Collection("users").Doc(userID).Collection("recoveryKeys").Doc(fmt.Sprintf("%d", version))
	snapshot, err := doc.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve encrypted user recovery key version %d: %w", version, err)
	}
	data := snapshot.Data()
	recoveryKeyRaw, ok := data["recoveryKey"]
	if !ok {
		return nil, fmt.Errorf("recoveryKey not found for user %s version %d", userID, version)
	}
	recoveryKeyBytes, ok := recoveryKeyRaw.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid type for recoveryKey for user %s version %d", userID, version)
	}
	return recoveryKeyBytes, nil
}

func (r *RecoveryKeyStorage) RetrieveLatestRecoveryKeyVersion(userID string) (int, error) {
	ctx := context.Background()
	doc := r.firestore.Collection("users").Doc(userID).Collection("recoveryKeys_metadata").Doc("metadata")
	snapshot, err := doc.Get(ctx)
	if errFromError, ok := status.FromError(err); ok {
		if errFromError.Code() == codes.NotFound {
			return 0, vault.ErrRecoveryKeyVersionNotFound
		}
	}
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve latest version: %w", err)
	}
	data := snapshot.Data()
	versionRaw, ok := data["latestVersion"]
	if !ok {
		return 0, fmt.Errorf("latestVersion not found for user %s", userID)
	}
	version, ok := versionRaw.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid type for latestVersion for user %s", userID)
	}
	return int(version), nil
}

func (r *RecoveryKeyStorage) StoreLatestRecoveryKeyVersion(userID string, version int) error {
	ctx := context.Background()
	doc := r.firestore.Collection("users").Doc(userID).Collection("recoveryKeys_metadata").Doc("metadata")
	_, err := doc.Set(ctx, map[string]any{
		"latestVersion": version,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to store latest version: %w", err)
	}
	return nil
}

func (r *RecoveryKeyStorage) RetrieveLastRotationTimestamp(userID string) (time.Time, error) {
	ctx := context.Background()
	doc := r.firestore.Collection("users").Doc(userID).Collection("recoveryKeys_metadata").Doc("metadata")
	snapshot, err := doc.Get(ctx)
	if errFromError, ok := status.FromError(err); ok {
		if errFromError.Code() == codes.NotFound {
			return time.Time{}, vault.ErrRotationTimestampNotFound
		}
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to retrieve last rotation timestamp: %w", err)
	}
	data := snapshot.Data()
	timestampRaw, ok := data["lastRotationAt"]
	if !ok {
		return time.Time{}, vault.ErrRotationTimestampNotFound
	}
	timestamp, ok := timestampRaw.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid type for lastRotationAt for user %s", userID)
	}
	return timestamp, nil
}

func (r *RecoveryKeyStorage) StoreLastRotationTimestamp(userID string, timestamp time.Time) error {
	ctx := context.Background()
	doc := r.firestore.Collection("users").Doc(userID).Collection("recoveryKeys_metadata").Doc("metadata")
	_, err := doc.Set(ctx, map[string]any{
		"lastRotationAt": timestamp,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to store last rotation timestamp: %w", err)
	}
	return nil
}

func (r *RecoveryKeyStorage) DeleteOldRecoveryKeyVersions(userID string, keepLatestN int) error {
	ctx := context.Background()
	collection := r.firestore.Collection("users").Doc(userID).Collection("recoveryKeys")
	snapshots, err := collection.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to list recovery key versions: %w", err)
	}
	if len(snapshots) <= keepLatestN {
		return nil
	}
	toDelete := snapshots[:len(snapshots)-keepLatestN]
	for _, snap := range toDelete {
		_, err := snap.Ref.Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete old recovery key version: %w", err)
		}
	}
	return nil
}
