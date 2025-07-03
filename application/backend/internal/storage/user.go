package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/7d4b9/utrade/backend/firebase"
	"github.com/7d4b9/utrade/backend/internal/http/api"
)

type UserStorage struct {
	firestore *firestore.Client
}

func NewUserStorage(firestore *firestore.Client) *UserStorage {
	return &UserStorage{
		firestore: firestore,
	}
}

func (s *UserStorage) StoreUser(ctx context.Context, userID string, o api.User) error {
	if o.ID == "" {
		return fmt.Errorf("missing user ID")
	}
	_, err := s.firestore.Collection("Users").Doc(o.ID).Set(ctx, o)
	return err
}

func (s *UserStorage) GetUser(ctx context.Context, id string) (*api.User, error) {
	doc := s.firestore.Collection("Users").Doc(id)
	snap, err := doc.Get(ctx)
	if firebase.FirestoreIsNotFound(err) {
		return nil, api.ErrUserNotFound
	}

	var o api.User
	if err := snap.DataTo(&o); err != nil {
		return nil, err
	}

	return &o, nil
}

func (s *UserStorage) ListUsers(ctx context.Context, userID string) ([]*api.User, error) {
	docs, err := s.firestore.Collection("Users").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var users []*api.User
	for _, doc := range docs {
		var u api.User
		if err := doc.DataTo(&u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (s *UserStorage) DeleteUser(ctx context.Context, userID, id string) error {
	doc := s.firestore.Collection("Users").Doc(id)
	snap, err := doc.Get(ctx)
	if firebase.FirestoreIsNotFound(err) {
		return api.ErrUserNotFound
	}
	var o api.User
	if err := snap.DataTo(&o); err != nil {
		return err
	}
	_, err = doc.Delete(ctx)
	return err
}

func (s *UserStorage) UpdateUserStatus(ctx context.Context, userID, UserID, status string) error {
	doc := s.firestore.Collection("Users").Doc(UserID)
	snap, err := doc.Get(ctx)
	if firebase.FirestoreIsNotFound(err) {
		return api.ErrUserNotFound
	}
	if err != nil {
		return err
	}
	var o api.User
	if err := snap.DataTo(&o); err != nil {
		return err
	}
	_, err = doc.Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
	})
	return err
}
