package vegetable

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/7d4b9/utrade/backend/internal/http/api"
	"github.com/7d4b9/utrade/images/vegetable"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config = viper.New()

const (
	gcloudProjectIDConfig                            = "google_cloud_project_id"
	validatedImagesModeratorPubSubSubscriptionConfig = "validated_images_backend_pubsub_subscription"
	createdImagesModeratorPubSubTopic                = "created_images_moderator_pubsub_topic"
)

func init() {
	config.AutomaticEnv()
	config.SetEnvPrefix("vegetable")

	config.BindEnv(gcloudProjectIDConfig, "GCLOUD_PROJECT_ID")

	config.BindEnv(createdImagesModeratorPubSubTopic, "CREATED_IMAGES_MODERATOR_PUBSUB_TOPIC")
}

type Storage interface {
	SetVegetableImageUploaded(ctx context.Context, vegetableID string, imageIndex int, imagePath string) error
}

type VegetableClient struct {
	validatedImagesSubscriptionID string

	storage                    Storage
	exit                       func()
	pubSubClient               *pubsub.Client
	validateImagesPublishTopic *pubsub.Topic
}

func NewVegetableClient(storage Storage) (*VegetableClient, error) {
	gcloudProjectID := config.GetString(gcloudProjectIDConfig)
	pubSubClient, err := pubsub.NewClient(context.Background(), gcloudProjectID)
	if err != nil {
		return nil, fmt.Errorf("cannot create pubsub client: %w", err)
	}

	validateImagePublishTopic := config.GetString(createdImagesModeratorPubSubTopic)
	if validateImagePublishTopic == "" {
		return nil, fmt.Errorf("image moderator input topic not set in config: %s", createdImagesModeratorPubSubTopic)
	}

	validatedImagesTopicSubscription := config.GetString(validatedImagesModeratorPubSubSubscriptionConfig)
	if validatedImagesTopicSubscription == "" {
		return nil, fmt.Errorf("image moderator subscription not set in config: %s", validatedImagesModeratorPubSubSubscriptionConfig)
	}

	validateImagesPublishTopic := pubSubClient.Topic(validateImagePublishTopic)
	v := &VegetableClient{
		pubSubClient:                  pubSubClient,
		storage:                       storage,
		validateImagesPublishTopic:    validateImagesPublishTopic,
		validatedImagesSubscriptionID: validatedImagesTopicSubscription,
	}

	log.Info().Str("image_moderator_input_topic", validateImagePublishTopic).Msg("Vegetable client running")

	return v, nil
}

func (v *VegetableClient) Close() {
	if v.exit != nil {
		v.exit()
	}
	v.exit = nil
	if v.pubSubClient != nil {
		if err := v.pubSubClient.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close pubsub client")
		}
	}
	v.pubSubClient = nil
	log.Info().Msg("Vegetable client closed")
}

// RunValidatedImagesSubscription starts the subscription to receive validated images from the pubsub topic.
// It will block until the context is canceled or an error occurs.
func (v *VegetableClient) RunValidatedImagesSubscription(ctx context.Context) error {
	if err := v.receiveValidatedImages(ctx); err != nil {
		return fmt.Errorf("failed to receive validated images: %w", err)
	}
	return nil
}

func (v *VegetableClient) receiveValidatedImages(ctx context.Context) error {

	validatedImagesSub := v.pubSubClient.Subscription(v.validatedImagesSubscriptionID)

	validatedImagesSub.ReceiveSettings.MaxOutstandingMessages = 500
	validatedImagesSub.ReceiveSettings.NumGoroutines = 100

	err := validatedImagesSub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("Received message: %s", string(msg.Data))
		var payload vegetable.VegetableValidatedImageMessage
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			log.Error().Err(err).Msg("Invalid message payload")
			msg.Nack()
			return
		}
		log.Debug().Str("vegetable_id", payload.VegetableID).Msg("Processing validated image")

		// Acknowledge the message before processing to avoid reprocessing in case of errors
		defer msg.Ack()
		if err := v.storage.SetVegetableImageUploaded(ctx, payload.VegetableID, payload.ImageIndex, payload.ImagePath); err != nil {
			log.Error().Fields(map[string]any{
				"vegetable_id": payload.VegetableID,
				"image_index":  payload.ImageIndex,
				"data":         string(msg.Data),
			}).Err(err).Msg("Failed to update vegetable image URL")
			// msg.Nack()
			return
		}

		// msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("failed to receive pubsub subscription: %w", err)
	}
	return nil
}

func (v *VegetableClient) SetImageValidation(ctx context.Context, vegetableID string, img *api.VegetableImage, imageIndex int) error {
	payload, err := json.Marshal(&vegetable.VegetableCreatedImageMessage{
		VegetableID: vegetableID,
		ImageIndex:  imageIndex,
		ImagePath:   img.Path,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal image validation message: %w", err)
	}
	res := v.validateImagesPublishTopic.Publish(ctx, &pubsub.Message{Data: payload})
	if _, err := res.Get(ctx); err != nil {
		return fmt.Errorf("failed to publish image validation message: %w", err)
	}
	v.validateImagesPublishTopic.Flush()
	log.Debug().
		Fields(map[string]any{
			"vegetable_id": vegetableID,
		}).
		Msg("Published image validation message")
	return nil
}
