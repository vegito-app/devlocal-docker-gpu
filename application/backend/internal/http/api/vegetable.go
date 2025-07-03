package api

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"strconv"
	"sync"
	"time"

	"github.com/7d4b9/utrade/backend/internal/http"
	v1 "github.com/7d4b9/utrade/backend/internal/http/api/internal/v1"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var ErrVegetableNotFound = errors.New("vegetable not found")

type VegetableImage struct {
	Path          string               `json:"path"`
	UploadedAt    time.Time            `json:"uploadedAt"`
	Status        VegetableImageStatus `json:"status"`
	DownloadToken *string              `json:"downloadToken,omitempty"`
}

type VegetableImageStatus string

const (
	VegetableImageStatusPending  = VegetableImageStatus(v1.VegetableImageStatusPending)
	VegetableImageStatusUploaded = VegetableImageStatus(v1.VegetableImageStatusUploaded)
)

type Vegetable struct {
	Active            bool             `json:"active"`
	AvailabilityDate  time.Time        `json:"availabilityDate"`
	AvailabilityType  string           `json:"availabilityType"`
	CreatedAt         time.Time        `json:"createdAt"`
	Description       string           `json:"description"`
	ID                string           `json:"id,omitempty"`
	Images            []VegetableImage `json:"images"`
	Name              string           `json:"name"`
	OwnerID           string           `json:"ownerId"`
	PriceCents        int              `json:"priceCents"`
	QuantityAvailable int              `json:"quantityAvailable"`
	SaleType          string           `json:"saleType"`
	UserCreatedAt     time.Time        `json:"userCreatedAt,omitempty"`

	// 🆕 Champs de géolocalisation
	DeliveryRadiusKm float64 `json:"deliveryRadiusKm"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
}

// Implementations must return ErrVegetableNotFound when a requested vegetable
// does not exist so that HTTP handlers can translate it into a 404 response.
type VegetableStorage interface {
	StoreVegetable(ctx context.Context, userID string, v Vegetable) (err error)

	// GetVegetable returns a vegetable by its ID.
	// Returns ErrVegetableNotFound if the vegetable does not exist.
	GetVegetable(ctx context.Context, userID, id string) (*Vegetable, error)

	// ListVegetables returns all vegetables for a user.
	// If no vegetables are found, it returns an empty slice.
	ListVegetables(ctx context.Context, userID string) ([]*Vegetable, error)

	// ListAvailableVegetables returns the vegetables that can be delivered to a given location.
	ListAvailableVegetables(ctx context.Context, deliveryRadiusKm float64, lat float64, lon float64, keyword *string) ([]*Vegetable, error)

	// DeleteVegetable removes a vegetable by ID.
	// Returns ErrVegetableNotFound if the vegetable does not exist.
	DeleteVegetable(ctx context.Context, userID, id string) error

	// UpdateMainImage updates the image order by placing the selected image at index 0.
	UpdateMainImage(ctx context.Context, userID, vegetableID string, mainImageCurrentIndex int) error
}

// VegetableImageValidator defines the interface for validating vegetable images.
type VegetableImageValidator interface {
	// SetImageValidation should schedules the validation of a vegetable image and immediately
	// returns without blocking. The image will be validated asynchronously.
	// It should return an error if the image cannot be scheduled for validation.
	SetImageValidation(ctx context.Context, userID string, image *VegetableImage, imageIndex int) error
}

// VegetableService defines all routes handled by VegetableService
type VegetableService struct {
	imageValidator               VegetableImageValidator
	isValidatedImagesServedByCDN bool
	storage                      VegetableStorage
}

// NewVegetableService initializes the VegetableService with the provided storage and image validator.
// It also sets up the HTTP handlers for the vegetable-related endpoints.
func NewVegetableService(mux *nethttp.ServeMux, storage VegetableStorage, imageValidator VegetableImageValidator) (*VegetableService, error) {
	isValidatedImagesCDN := config.GetBool(isValidatedImagesCDNEnabled)
	service := &VegetableService{
		storage:                      storage,
		imageValidator:               imageValidator,
		isValidatedImagesServedByCDN: isValidatedImagesCDN,
	}
	mux.Handle("POST /api/vegetables", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.CreateVegetable),
		FirebaseAuthMiddleware))
	mux.Handle("GET /api/vegetables", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.ListVegetables),
		FirebaseAuthMiddleware))
	mux.Handle("GET /api/vegetables/available", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.ListAvailableVegetables),
		FirebaseAuthMiddleware))
	mux.Handle("GET /api/vegetable", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.GetVegetable),
		FirebaseAuthMiddleware))
	mux.Handle("DELETE /api/vegetables/{id}", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.DeleteVegetable),
		FirebaseAuthMiddleware))
	mux.Handle("PUT /api/vegetables/{id}", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.UpdateVegetable),
		FirebaseAuthMiddleware))
	mux.Handle("PUT /api/vegetables/{id}/main-image", http.ApplyMiddleware(
		nethttp.HandlerFunc(service.UpdateMainImage),
		FirebaseAuthMiddleware))
	return service, nil
}

// CreateVegetable handles the creation of a new vegetable.
// It decodes the request body into a VegetableRequest, validates it, and stores it in the database.
// If the vegetable is successfully created, it returns a 201 Created response with the vegetable details.
func (s *VegetableService) CreateVegetable(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	requestBody, ok := http.RequestBodyFromContext(r)
	if !ok {
		log.Error().Msg("create vegetable failed to get request body from context")
		nethttp.Error(w, "failed to read request body", nethttp.StatusInternalServerError)
		return
	}
	if len(requestBody) == 0 {
		log.Error().Msg("create vegetable failed: empty request body")
		nethttp.Error(w, "empty request body", nethttp.StatusBadRequest)
		return
	}
	var vegetableRequest v1.VegetableRequest
	// Decode the request body into VegetableRequest
	if err := json.Unmarshal(requestBody, &vegetableRequest); err != nil {
		log.Error().Err(err).Msg("create vegetable failed to decode vegetable payload")
		nethttp.Error(w, "invalid payload", nethttp.StatusBadRequest)
		return
	}
	if vegetableRequest.ID != "" {
		nethttp.Error(w, "ID must not be provided on create", nethttp.StatusBadRequest)
		return
	}
	vegetable := Vegetable{
		Active:            true, // always active on creation
		AvailabilityDate:  vegetableRequest.AvailabilityDate,
		AvailabilityType:  string(vegetableRequest.AvailabilityType),
		CreatedAt:         time.Now().UTC(),
		DeliveryRadiusKm:  vegetableRequest.DeliveryRadiusKm,
		Description:       vegetableRequest.Description,
		ID:                uuid.NewString(),
		Latitude:          vegetableRequest.Latitude,
		Longitude:         vegetableRequest.Longitude,
		Name:              vegetableRequest.Name,
		OwnerID:           vegetableRequest.OwnerID,
		PriceCents:        vegetableRequest.PriceCents,
		QuantityAvailable: vegetableRequest.QuantityAvailable,
		SaleType:          vegetableRequest.SaleType,
		UserCreatedAt:     time.Now().UTC(),
	}
	for _, img := range vegetableRequest.Images {
		vegetableImage := VegetableImage{
			Path:          img.Path,
			UploadedAt:    time.Now().UTC(),
			Status:        VegetableImageStatusPending,
			DownloadToken: img.DownloadToken,
		}
		vegetable.Images = append(vegetable.Images, vegetableImage)
	}
	var wg sync.WaitGroup
	defer wg.Wait()
	// Store the vegetable with images that need validation blanked out.
	// Do this first to not blank out images that are already validated.
	err := s.storage.StoreVegetable(ctx, userID, vegetable)
	if err != nil {
		nethttp.Error(w, "store failed", nethttp.StatusInternalServerError)
		return
	}
	for i, img := range vegetableRequest.Images {
		wg.Add(1)
		go func(index int, img *v1.VegetableImageRequest) {
			defer wg.Done()
			vegetableImage := &VegetableImage{
				Path:       img.Path,
				UploadedAt: time.Now().UTC(),
				Status:     VegetableImageStatusPending,
			}
			if err := s.imageValidator.SetImageValidation(ctx, vegetable.ID, vegetableImage, index); err != nil {
				log.Error().Err(err).Msg("failed to set image validation")
			}
		}(i, &img)
	}
	createdVegetableImagesResponse := responseImages(&vegetable, userID, s.isValidatedImagesServedByCDN)
	response := &v1.VegetableResponse{
		// WeightGrams:       vegetable.WeightGrams,
		Active:            vegetable.Active,
		AvailabilityDate:  &vegetable.AvailabilityDate,
		AvailabilityType:  v1.VegetableAvailabilityType(vegetable.AvailabilityType),
		CreatedAt:         vegetable.CreatedAt,
		DeliveryRadiusKm:  vegetable.DeliveryRadiusKm,
		Description:       vegetable.Description,
		ID:                vegetable.ID,
		Images:            createdVegetableImagesResponse,
		Latitude:          vegetable.Latitude,
		Longitude:         vegetable.Longitude,
		Name:              vegetable.Name,
		OwnerID:           vegetable.OwnerID,
		PriceCents:        vegetable.PriceCents,
		QuantityAvailable: vegetable.QuantityAvailable,
		SaleType:          vegetable.SaleType,
		UserCreatedAt:     vegetable.UserCreatedAt,
	}
	log.Debug().Str("id", vegetable.ID).Msg("Created vegetable")
	// Set the response status to Created (201)
	w.WriteHeader(nethttp.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("failed to encode vegetable response")
		nethttp.Error(w, "failed to encode response", nethttp.StatusInternalServerError)
		return
	}
}

// ListVegetables retrieves all vegetables for the authenticated user.
// It returns a list of vegetables with their images enriched based on the user's permissions.
// If no vegetables are found, it returns an empty slice.
func (s *VegetableService) ListVegetables(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	veggies, err := s.storage.ListVegetables(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to list vegetables")
		nethttp.Error(w, "list failed", nethttp.StatusInternalServerError)
		return
	}
	if veggies == nil {
		veggies = []*Vegetable{}
	}
	// Prepare response with enriched images
	resp := []*v1.VegetableResponse{}
	for _, veggie := range veggies {
		imagesResp := responseImages(veggie, userID, s.isValidatedImagesServedByCDN)
		resp = append(resp, &v1.VegetableResponse{
			// WeightGrams:       veggie.WeightGrams,
			Active:            veggie.Active,
			AvailabilityDate:  &veggie.AvailabilityDate,
			AvailabilityType:  v1.VegetableAvailabilityType(veggie.AvailabilityType),
			CreatedAt:         veggie.CreatedAt,
			DeliveryRadiusKm:  veggie.DeliveryRadiusKm,
			Description:       veggie.Description,
			ID:                veggie.ID,
			Images:            imagesResp,
			Latitude:          veggie.Latitude,
			Longitude:         veggie.Longitude,
			Name:              veggie.Name,
			OwnerID:           veggie.OwnerID,
			PriceCents:        veggie.PriceCents,
			QuantityAvailable: veggie.QuantityAvailable,
			SaleType:          veggie.SaleType,
			UserCreatedAt:     veggie.UserCreatedAt,
		})
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msg("failed to encode vegetables response")
		nethttp.Error(w, "failed to encode response", nethttp.StatusInternalServerError)
		return
	}
	log.Debug().Int("count", len(veggies)).Msg("Listed vegetables")
}

// responseImages prepares the vegetable images response, handling visibility based on user ID and CDN status.
// If the image is pending and the user is not the owner, it will not return the path as the image is not yet validated.
func responseImages(veggie *Vegetable, userID string, isValidatedImagesServedByCDN bool) []v1.VegetableImageResponse {
	var imagesResp []v1.VegetableImageResponse
	for i := range veggie.Images {
		var path string
		if veggie.Images[i].Status == VegetableImageStatusPending && veggie.OwnerID != userID && !isValidatedImagesServedByCDN {
			path = ""
		} else {
			path = veggie.Images[i].Path
		}
		servedByCdn := (veggie.Images[i].Status == VegetableImageStatusUploaded) && isValidatedImagesServedByCDN
		imagesResp = append(imagesResp, v1.VegetableImageResponse{
			Path:          path,
			UploadedAt:    veggie.Images[i].UploadedAt,
			Status:        v1.VegetableImageStatus(veggie.Images[i].Status),
			ServedByCdn:   servedByCdn,
			DownloadToken: veggie.Images[i].DownloadToken,
		})
	}
	return imagesResp
}

// GetVegetable retrieves a specific vegetable by its ID for the authenticated user.
// It returns a 404 error if the vegetable does not exist or is not found.
func (s *VegetableService) GetVegetable(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	id := r.URL.Query().Get("id")
	if id == "" {
		nethttp.Error(w, "missing id", nethttp.StatusBadRequest)
		return
	}
	veggie, err := s.storage.GetVegetable(ctx, userID, id)
	if err != nil {
		if errors.Is(err, ErrVegetableNotFound) {
			nethttp.Error(w, "vegetable not found", nethttp.StatusNotFound)
			return
		}
		nethttp.Error(w, "get failed", nethttp.StatusInternalServerError)
		return
	}
	// Prepare enriched image response
	imagesResp := responseImages(veggie, userID, s.isValidatedImagesServedByCDN)
	resp := &v1.VegetableResponse{
		Active:            veggie.Active,
		AvailabilityDate:  &veggie.AvailabilityDate,
		AvailabilityType:  v1.VegetableAvailabilityType(veggie.AvailabilityType),
		CreatedAt:         veggie.CreatedAt,
		DeliveryRadiusKm:  veggie.DeliveryRadiusKm,
		Description:       veggie.Description,
		ID:                veggie.ID,
		Images:            imagesResp,
		Latitude:          veggie.Latitude,
		Longitude:         veggie.Longitude,
		Name:              veggie.Name,
		OwnerID:           veggie.OwnerID,
		PriceCents:        veggie.PriceCents,
		QuantityAvailable: veggie.QuantityAvailable,
		SaleType:          veggie.SaleType,
		UserCreatedAt:     veggie.UserCreatedAt,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msg("failed to encode vegetable response")
		nethttp.Error(w, "failed to encode response", nethttp.StatusInternalServerError)
		return
	}
	log.Debug().Str("id", id).Msg("Retrieved vegetable")
}

// DeleteVegetable removes a vegetable by its ID for the authenticated user.
// It returns a 404 error if the vegetable does not exist or is not found.
func (s *VegetableService) DeleteVegetable(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	id := r.PathValue("id")
	if id == "" {
		nethttp.Error(w, "missing id", nethttp.StatusBadRequest)
		return
	}
	if err := s.storage.DeleteVegetable(ctx, userID, id); err != nil {
		if errors.Is(err, ErrVegetableNotFound) {
			nethttp.Error(w, "vegetable not found", nethttp.StatusNotFound)
			return
		}
		nethttp.Error(w, "delete failed", nethttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

// UpdateVegetable handles updating an existing vegetable.
func (s *VegetableService) UpdateVegetable(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	id := r.PathValue("id")
	if id == "" {
		nethttp.Error(w, "missing id", nethttp.StatusBadRequest)
		return
	}
	requestBody, ok := http.RequestBodyFromContext(r)
	if !ok {
		log.Error().Msg("create vegetable failed to get request body from context")
		nethttp.Error(w, "failed to read request body", nethttp.StatusInternalServerError)
		return
	}
	if len(requestBody) == 0 {
		log.Error().Msg("create vegetable failed: empty request body")
		nethttp.Error(w, "empty request body", nethttp.StatusBadRequest)
		return
	}
	var vegetableRequest v1.VegetableRequest
	if err := json.Unmarshal(requestBody, &vegetableRequest); err != nil {
		log.Error().Err(err).Msg("update vegetable failed to decode vegetable payload")
		nethttp.Error(w, "invalid payload", nethttp.StatusBadRequest)
		return
	}
	veggie, err := s.storage.GetVegetable(ctx, userID, id)
	if err != nil {
		if errors.Is(err, ErrVegetableNotFound) {
			nethttp.Error(w, "vegetable not found", nethttp.StatusNotFound)
			return
		}
		nethttp.Error(w, "get failed", nethttp.StatusInternalServerError)
		return
	}

	// Update fields
	veggie.Name = vegetableRequest.Name
	veggie.Description = vegetableRequest.Description
	veggie.SaleType = vegetableRequest.SaleType
	// veggie.WeightGrams = vegetableRequest.WeightGrams
	veggie.PriceCents = vegetableRequest.PriceCents
	veggie.Active = vegetableRequest.Active
	veggie.AvailabilityType = string(vegetableRequest.AvailabilityType)
	veggie.AvailabilityDate = vegetableRequest.AvailabilityDate
	veggie.QuantityAvailable = int(vegetableRequest.QuantityAvailable)
	veggie.Images = nil
	for _, img := range vegetableRequest.Images {
		vegetableImage := VegetableImage{
			Path:          img.Path,
			UploadedAt:    time.Now().UTC(),
			Status:        VegetableImageStatusPending,
			DownloadToken: img.DownloadToken,
		}
		veggie.Images = append(veggie.Images, vegetableImage)
	}

	err = s.storage.StoreVegetable(ctx, userID, *veggie)
	if err != nil {
		nethttp.Error(w, "store failed", nethttp.StatusInternalServerError)
		return
	}

	imagesResp := responseImages(veggie, userID, s.isValidatedImagesServedByCDN)
	resp := &v1.VegetableResponse{
		ID:          veggie.ID,
		Name:        veggie.Name,
		OwnerID:     veggie.OwnerID,
		Description: veggie.Description,
		SaleType:    veggie.SaleType,
		// WeightGrams:       veggie.WeightGrams,
		PriceCents:        veggie.PriceCents,
		Images:            imagesResp,
		CreatedAt:         veggie.CreatedAt,
		UserCreatedAt:     veggie.UserCreatedAt,
		Active:            veggie.Active,
		AvailabilityType:  v1.VegetableAvailabilityType(veggie.AvailabilityType),
		AvailabilityDate:  &veggie.AvailabilityDate,
		QuantityAvailable: veggie.QuantityAvailable,
	}
	log.Debug().Str("id", id).Msg("Updated vegetable")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msg("failed to encode vegetable response")
		nethttp.Error(w, "failed to encode response", nethttp.StatusInternalServerError)
		return
	}
}

// UpdateMainImage updates the main image of a vegetable by moving the selected image to index 0.
// It expects the request body to contain the index of the current main image.
func (s *VegetableService) UpdateMainImage(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)
	id := r.PathValue("id")
	if id == "" {
		nethttp.Error(w, "missing id", nethttp.StatusBadRequest)
		return
	}
	requestBody, ok := http.RequestBodyFromContext(r)
	if !ok {
		nethttp.Error(w, "failed to get request body", nethttp.StatusInternalServerError)
		return
	}
	var payload v1.UpdateMainImageRequest
	if err := json.Unmarshal(requestBody, &payload); err != nil {
		log.Error().Err(err).Msg("failed to decode main image index")
		nethttp.Error(w, "invalid payload", nethttp.StatusBadRequest)
		return
	}
	if err := s.storage.UpdateMainImage(ctx, userID, id, payload.MainImageIndex); err != nil {
		if errors.Is(err, ErrVegetableNotFound) {
			nethttp.Error(w, "vegetable not found", nethttp.StatusNotFound)
			return
		}
		nethttp.Error(w, "update main image failed", nethttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

// ListAvailableVegetables returns vegetables available for a given location and optional keyword.
func (s *VegetableService) ListAvailableVegetables(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := r.Context()
	userID := authenticatedUserID(r)

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	deliveryRadiusKMstr := r.URL.Query().Get("deliveryRadiusKm")
	keyword := r.URL.Query().Get("keyword")

	if latStr == "" || lonStr == "" || deliveryRadiusKMstr == "" {
		nethttp.Error(w, "missing required parameters: lat, lon, deliveryRadiusKm", nethttp.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		nethttp.Error(w, "invalid lat", nethttp.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		nethttp.Error(w, "invalid lon", nethttp.StatusBadRequest)
		return
	}
	radiusKm, err := strconv.ParseFloat(deliveryRadiusKMstr, 64)
	if err != nil {
		nethttp.Error(w, "invalid radiusKm", nethttp.StatusBadRequest)
		return
	}

	var kw *string
	if keyword != "" {
		kw = &keyword
	}

	veggies, err := s.storage.ListAvailableVegetables(ctx, radiusKm, lat, lon, kw)
	if err != nil {
		log.Error().Err(err).Msg("failed to list available vegetables")
		nethttp.Error(w, "internal error", nethttp.StatusInternalServerError)
		return
	}
	if veggies == nil {
		veggies = []*Vegetable{}
	}

	resp := []*v1.VegetableResponse{}
	for _, veggie := range veggies {
		imagesResp := responseImages(veggie, userID, s.isValidatedImagesServedByCDN)
		resp = append(resp, &v1.VegetableResponse{
			ID:                veggie.ID,
			Name:              veggie.Name,
			OwnerID:           veggie.OwnerID,
			Description:       veggie.Description,
			SaleType:          veggie.SaleType,
			PriceCents:        veggie.PriceCents,
			Images:            imagesResp,
			CreatedAt:         veggie.CreatedAt,
			UserCreatedAt:     veggie.UserCreatedAt,
			Active:            veggie.Active,
			AvailabilityType:  v1.VegetableAvailabilityType(veggie.AvailabilityType),
			AvailabilityDate:  &veggie.AvailabilityDate,
			QuantityAvailable: veggie.QuantityAvailable,
		})
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
		nethttp.Error(w, "failed to encode response", nethttp.StatusInternalServerError)
		return
	}
}
