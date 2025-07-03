package storage

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/7d4b9/utrade/backend/firebase"
	"github.com/7d4b9/utrade/backend/internal/http/api"
	"google.golang.org/api/iterator"
)

type OrderStorage struct {
	firestore *firestore.Client
}

func NewOrderStorage(firestore *firestore.Client) *OrderStorage {
	return &OrderStorage{
		firestore: firestore,
	}
}

func (s *OrderStorage) StoreOrder(ctx context.Context, userID string, o api.Order) error {
	// Implementation goes here
	return nil
}

func (s *OrderStorage) GetOrder(ctx context.Context, userID, id string) (*api.Order, error) {
	doc := s.firestore.Collection("orders").Doc(id)
	snap, err := doc.Get(ctx)
	if firebase.FirestoreIsNotFound(err) {
		return nil, api.ErrOrderNotFound
	}

	var p persistentOrder
	if err := snap.DataTo(&p); err != nil {
		return nil, err
	}

	return &p.Order, nil
}

func (s *OrderStorage) ListOrders(ctx context.Context, userID string) ([]*api.Order, error) {
	// Implementation goes here

	return nil, nil
}

type persistentOrder struct {
	api.Order
	Deleted bool `firestore:"deleted,omitempty"`
}

func (s *OrderStorage) DeleteOrder(ctx context.Context, userID, id string) error {
	doc := s.firestore.Collection("orders").Doc(id)

	return s.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(doc)
		if firebase.FirestoreIsNotFound(err) {
			return api.ErrOrderNotFound
		}
		if err != nil {
			return err
		}

		var p persistentOrder
		if err := snap.DataTo(&p); err != nil {
			return err
		}
		o := p.Order

		if o.Status != "draft" {
			// Soft delete
			if err := tx.Update(doc, []firestore.Update{{Path: "deleted", Value: true}}); err != nil {
				return err
			}
			return nil
		}

		// Hard delete
		if err := tx.Delete(doc); err != nil {
			return err
		}

		// Conditional cleanup of associated vegetables
		for _, vegetableID := range o.VegetableIds {
			vegDoc := s.firestore.Collection("vegetables").Doc(vegetableID)
			vegSnap, err := tx.Get(vegDoc)
			if err != nil {
				if firebase.FirestoreIsNotFound(err) {
					continue
				}
				return err
			}

			var veg persistentVegetable
			if err := vegSnap.DataTo(&veg); err != nil {
				return err
			}

			if !veg.Deleted {
				// Only consider vegetables that are already soft deleted
				continue
			}

			orderIter := s.firestore.Collection("orders").
				Where("vegetableIds", "array-contains", vegetableID).
				Where("deleted", "!=", true).
				Limit(1).Documents(ctx)

			exists := false
			if _, err := orderIter.Next(); err != nil && err != iterator.Done {
				return err
			} else if err == nil {
				exists = true
			}

			if !exists {
				// Delete subcollection "images" in the same transaction
				imageIter := vegDoc.Collection("images").Documents(ctx)
				for {
					imageDoc, err := imageIter.Next()
					if err != nil {
						if err == iterator.Done {
							break
						}
						return err
					}
					if err := tx.Delete(imageDoc.Ref); err != nil {
						return err
					}
				}
				if err := tx.Delete(vegDoc); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *OrderStorage) UpdateOrderStatus(ctx context.Context, userID, orderID, status string) error {
	doc := s.firestore.Collection("orders").Doc(orderID)
	snap, err := doc.Get(ctx)
	if firebase.FirestoreIsNotFound(err) {
		return api.ErrOrderNotFound
	}
	if err != nil {
		return err
	}
	var p persistentOrder
	if err := snap.DataTo(&p); err != nil {
		return err
	}
	_, err = doc.Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
	})
	return err
}
