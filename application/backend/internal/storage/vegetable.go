package storage

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/7d4b9/utrade/backend/firebase"
	"github.com/7d4b9/utrade/backend/internal/http/api"
	"google.golang.org/api/iterator"
)

type persistentVegetable struct {
	api.Vegetable
	Deleted bool `firestore:"deleted,omitempty"`
}

type VegetableStorage struct {
	firestore *firestore.Client
}

func NewVegetableStorage(firestore *firestore.Client) *VegetableStorage {
	return &VegetableStorage{
		firestore: firestore,
	}
}

func (s *VegetableStorage) StoreVegetable(ctx context.Context, userID string, v api.Vegetable) error {
	data := map[string]any{
		"name":              v.Name,
		"description":       v.Description,
		"saleType":          v.SaleType,
		"priceCents":        v.PriceCents,
		"ownerId":           userID,
		"createdAt":         firestore.ServerTimestamp,
		"userCreatedAt":     v.UserCreatedAt.UTC(),
		"active":            v.Active,
		"availabilityType":  v.AvailabilityType,
		"availabilityDate":  v.AvailabilityDate,
		"quantityAvailable": v.QuantityAvailable,
		"deleted":           false,

		// Champs de géolocalisation
		"latitude":         v.Latitude,
		"longitude":        v.Longitude,
		"deliveryRadiusKm": v.DeliveryRadiusKm,
	}

	_, err := s.firestore.Collection("vegetables").Doc(v.ID).Set(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to store vegetable %q: %w", v.ID, err)
	}

	for i, img := range v.Images {
		imageDoc := s.firestore.Collection("vegetables").
			Doc(v.ID).
			Collection("images").
			Doc(strconv.Itoa(i))

		data := map[string]interface{}{
			"path":   img.Path,
			"status": img.Status,
		}
		if img.DownloadToken != nil {
			data["downloadToken"] = *img.DownloadToken
		}
		if _, err := imageDoc.Set(ctx, data, firestore.MergeAll); err != nil {
			return fmt.Errorf("failed to store image %q for vegetable %q: %w", strconv.Itoa(i), v.ID, err)
		}
	}

	return nil
}

func (f *Storage) GetVegetable(ctx context.Context, userID, id string) (*api.Vegetable, error) {
	doc, err := f.firestore.Collection("vegetables").Doc(id).Get(ctx)
	if err != nil {
		if firebase.FirestoreIsNotFound(err) {
			return nil, api.ErrVegetableNotFound
		}
		return nil, fmt.Errorf("failed to get vegetable %q: %w", id, err)
	}
	data := doc.Data()
	if data["ownerId"] != userID {
		return nil, fmt.Errorf("unauthorized access to vegetable %q", id)
	}
	var pv persistentVegetable
	if err := doc.DataTo(&pv); err != nil {
		return nil, fmt.Errorf("failed to convert document %q to Vegetable: %w", doc.Ref.ID, err)
	}
	v := pv.Vegetable
	if active, ok := data["active"].(bool); ok {
		v.Active = active
	} else {
		v.Active = true // fallback pour rétrocompatibilité
	}

	imagesIter := doc.Ref.Collection("images").
		OrderBy(firestore.DocumentID, firestore.Asc).
		Documents(ctx)
	defer imagesIter.Stop()
	var images []api.VegetableImage
	for {
		docSnap, err := imagesIter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return nil, fmt.Errorf("failed to iterate vegetable images: %w", err)
		}
		var img api.VegetableImage
		if err := docSnap.DataTo(&img); err != nil {
			return nil, fmt.Errorf("failed to decode vegetable image: %w", err)
		}
		if token, ok := docSnap.Data()["downloadToken"].(string); ok {
			img.DownloadToken = &token
		}
		images = append(images, img)
	}
	v.Images = images
	v.ID = doc.Ref.ID // Ensure the ID is set from the document reference
	return &v, nil
}

func (f *Storage) ListVegetables(ctx context.Context, userID string) ([]*api.Vegetable, error) {
	iter := f.firestore.Collection("vegetables").
		Where("ownerId", "==", userID).
		Where("deleted", "!=", true).
		Documents(ctx)
	defer iter.Stop()

	var list []*api.Vegetable
	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return nil, fmt.Errorf("error while listing vegetables: %w", err)
		}
		var pv persistentVegetable
		if err := doc.DataTo(&pv); err != nil {
			return nil, fmt.Errorf("failed to convert document %q to Vegetable: %w", doc.Ref.ID, err)
		}
		v := pv.Vegetable
		if active, ok := doc.Data()["active"].(bool); ok {
			v.Active = active
		} else {
			v.Active = true // fallback pour rétrocompatibilité
		}
		v.ID = doc.Ref.ID // Ensure the ID is set from the document reference

		imagesIter := doc.Ref.Collection("images").
			OrderBy(firestore.DocumentID, firestore.Asc).
			Documents(ctx)
		defer imagesIter.Stop()
		var images []api.VegetableImage
		for {
			docSnap, err := imagesIter.Next()
			if err != nil {
				if errors.Is(err, iterator.Done) {
					break
				}
				return nil, fmt.Errorf("failed to iterate vegetable images: %w", err)
			}
			var img api.VegetableImage
			if err := docSnap.DataTo(&img); err != nil {
				return nil, fmt.Errorf("failed to decode vegetable image: %w", err)
			}
			if token, ok := docSnap.Data()["downloadToken"].(string); ok {
				img.DownloadToken = &token
			}
			// img.ID = docSnap.Ref.ID
			images = append(images, img)
		}
		v.Images = images

		list = append(list, &v)
	}
	return list, nil
}

func (f *Storage) DeleteVegetable(ctx context.Context, userID, id string) error {
	doc := f.firestore.Collection("vegetables").Doc(id)
	snap, err := doc.Get(ctx)
	if err != nil {
		if firebase.FirestoreIsNotFound(err) {
			return fmt.Errorf("%w: %s", api.ErrVegetableNotFound, id)
		}
		return fmt.Errorf("failed to find vegetable %q: %w", id, err)
	}
	if snap.Data()["ownerId"] != userID {
		return fmt.Errorf("unauthorized delete attempt on vegetable %q", id)
	}

	return f.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		orders := f.firestore.Collection("orders")
		orderIter := orders.Where("vegetableIds", "array-contains", id).Limit(1).Documents(ctx)
		defer orderIter.Stop()

		exists := false
		if _, err := orderIter.Next(); err != nil && !errors.Is(err, iterator.Done) {
			return fmt.Errorf("failed to query orders for vegetable %q: %w", id, err)
		} else if err == nil {
			exists = true
		}

		if exists {
			if err := tx.Update(doc, []firestore.Update{{Path: "deleted", Value: true}}); err != nil {
				return fmt.Errorf("failed to soft delete vegetable %q: %w", id, err)
			}
			return nil
		}

		imagesIter := doc.Collection("images").Documents(ctx)
		for {
			imageDoc, err := imagesIter.Next()
			if err != nil {
				if errors.Is(err, iterator.Done) {
					break
				}
				return fmt.Errorf("failed to iterate images for deletion: %w", err)
			}
			if err := tx.Delete(imageDoc.Ref); err != nil {
				return fmt.Errorf("failed to delete image %q: %w", imageDoc.Ref.ID, err)
			}
		}

		if err := tx.Delete(doc); err != nil {
			return fmt.Errorf("failed to hard delete vegetable %q: %w", id, err)
		}
		return nil
	})
}

func (s *Storage) SetVegetableImageUploaded(ctx context.Context, vegetableID string, imageIndex int, imagePath string) error {
	err := s.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		vegDoc := s.firestore.Collection("vegetables").Doc(vegetableID)
		snap, err := tx.Get(vegDoc)
		if err != nil {
			if firebase.FirestoreIsNotFound(err) {
				return fmt.Errorf("cannot update image: vegetable %q not found", vegetableID)
			}
			return fmt.Errorf("failed to get vegetable %q: %w", vegetableID, err)
		}
		if !snap.Exists() {
			return fmt.Errorf("cannot update image: vegetable %q does not exist", vegetableID)
		}

		imageDoc := vegDoc.Collection("images").Doc(strconv.Itoa(imageIndex))
		data := map[string]interface{}{
			"status": api.VegetableImageStatusUploaded,
		}
		if err := tx.Set(imageDoc, data, firestore.MergeAll); err != nil {
			return fmt.Errorf("failed to set image status for vegetable %q image %d: %w", vegetableID, imageIndex, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to set image data for vegetable %q image %q: %w", vegetableID, imageIndex, err)
	}
	return nil
}

func (s *Storage) UpdateMainImage(ctx context.Context, userID, vegetableID string, mainImageCurrentIndex int) error {
	err := s.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		vegDoc := s.firestore.Collection("vegetables").Doc(vegetableID)
		snap, err := tx.Get(vegDoc)
		if err != nil {
			if firebase.FirestoreIsNotFound(err) {
				return api.ErrVegetableNotFound
			}
			return fmt.Errorf("failed to get vegetable %q: %w", vegetableID, err)
		}
		if snap.Data()["ownerId"] != userID {
			return fmt.Errorf("unauthorized main image update attempt on vegetable %q", vegetableID)
		}
		imagesColl := vegDoc.Collection("images")
		imageDocs, err := imagesColl.Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("failed to fetch images for vegetable %q: %w", vegetableID, err)
		}
		if mainImageCurrentIndex < 0 || mainImageCurrentIndex >= len(imageDocs) {
			return fmt.Errorf("invalid image index %d for vegetable %q", mainImageCurrentIndex, vegetableID)
		}

		type imageData struct {
			Path          string
			Status        string
			DownloadToken *string
		}

		var imageDataList []imageData
		for _, doc := range imageDocs {
			d := doc.Data()
			img := imageData{
				Path:   fmt.Sprintf("%v", d["path"]),
				Status: fmt.Sprintf("%v", d["status"]),
			}
			if t, ok := d["downloadToken"].(string); ok {
				img.DownloadToken = &t
			}
			imageDataList = append(imageDataList, img)
		}

		main := imageDataList[mainImageCurrentIndex]
		others := append(imageDataList[:mainImageCurrentIndex], imageDataList[mainImageCurrentIndex+1:]...)
		reordered := append([]imageData{main}, others...)

		for _, doc := range imageDocs {
			if err := tx.Delete(doc.Ref); err != nil {
				return fmt.Errorf("failed to delete image doc %q: %w", doc.Ref.ID, err)
			}
		}

		for i, img := range reordered {
			data := map[string]interface{}{
				"path":   img.Path,
				"status": img.Status,
			}
			if img.DownloadToken != nil {
				data["downloadToken"] = *img.DownloadToken
			}
			if err := tx.Set(imagesColl.Doc(strconv.Itoa(i)), data); err != nil {
				return fmt.Errorf("failed to reinsert image %d: %w", i, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update main image for vegetable %q: %w", vegetableID, err)
	}
	return nil
}

func (s *Storage) ListAvailableVegetables(ctx context.Context, lat float64, lon float64, radiusKm float64, keyword *string) ([]*api.Vegetable, error) {
	radiusDeg := radiusKm / 111.0

	iter := s.firestore.Collection("vegetables").
		Where("availabilityType", "==", "available").
		Where("deleted", "!=", true).
		Where("latitude", ">=", lat-radiusDeg).
		Where("latitude", "<=", lat+radiusDeg).
		Where("longitude", ">=", lon-radiusDeg).
		Where("longitude", "<=", lon+radiusDeg).
		Documents(ctx)
	defer iter.Stop()

	var vegetables []*api.Vegetable
	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return nil, fmt.Errorf("failed to iterate available vegetables: %w", err)
		}

		var veg api.Vegetable
		if err := doc.DataTo(&veg); err != nil {
			return nil, fmt.Errorf("failed to decode vegetable %q: %w", doc.Ref.ID, err)
		}

		// Vérification par distance
		dist := distanceKm(lat, lon, veg.Latitude, veg.Longitude)
		if dist > veg.DeliveryRadiusKm {
			continue
		}
		if veg.DeliveryRadiusKm < radiusKm {
			continue
		}

		// Vérification par mot-clé (dans le nom ou la description)
		if keyword != nil {
			kw := *keyword
			if !containsIgnoreCase(veg.Name, kw) && !containsIgnoreCase(veg.Description, kw) {
				continue
			}
		}

		vegetables = append(vegetables, &veg)
	}
	return vegetables, nil
}

func containsIgnoreCase(text, substr string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}

func distanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}
