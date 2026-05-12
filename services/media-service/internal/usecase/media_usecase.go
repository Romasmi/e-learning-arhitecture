package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elearning/media-service/internal/domain"
	"github.com/elearning/media-service/pkg/kafka"
	"github.com/google/uuid"
)

type MediaUsecase struct {
	repo          domain.AssetRepository
	producer      *kafka.Producer
	transcoderURL string
	s3Bucket      string
}

func NewMediaUsecase(repo domain.AssetRepository, producer *kafka.Producer, transcoderURL, s3Bucket string) *MediaUsecase {
	return &MediaUsecase{
		repo:          repo,
		producer:      producer,
		transcoderURL: transcoderURL,
		s3Bucket:      s3Bucket,
	}
}

func (u *MediaUsecase) UploadAsset(ctx context.Context, lessonID uuid.UUID, content []byte, filename string, assetType domain.AssetType) (*domain.Asset, error) {
	assetID := uuid.New()

	// Mock S3 upload - for simplicity, just generate a URL
	rawURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s", u.s3Bucket, assetID, filename)

	asset := &domain.Asset{
		ID:        assetID,
		LessonID:  lessonID,
		Type:      assetType,
		Status:    domain.AssetStatusPending,
		RawURL:    rawURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	if assetType == domain.AssetTypeVideo {
		// Submit transcode job
		go u.submitTranscodeJob(context.Background(), asset)

		// Emit VideoUploaded
		payload, _ := json.Marshal(asset)
		u.producer.PublishAsync(kafka.Event{
			EventType:  "VideoUploaded",
			AssetID:    asset.ID.String(),
			LessonID:   asset.LessonID.String(),
			Payload:    payload,
			OccurredAt: time.Now(),
		})
	}

	return asset, nil
}

func (u *MediaUsecase) submitTranscodeJob(ctx context.Context, asset *domain.Asset) {
	// Simple HTTP POST to transcoder
	jobReq := map[string]string{
		"asset_id": asset.ID.String(),
		"raw_url":  asset.RawURL,
		"callback": "/internal/transcode-callback",
	}
	body, _ := json.Marshal(jobReq)

	resp, err := http.Post(u.transcoderURL+"/jobs", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("failed to submit transcode job: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		var res map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
			jobID := res["job_id"]
			if jobID != "" {
				_ = u.repo.UpdateAssetJob(ctx, asset.ID, jobID)
			}
		}
	}
}

func (u *MediaUsecase) HandleTranscodeCallback(ctx context.Context, jobID, status string, cdnURLs []string) error {
	asset, err := u.repo.GetAssetByJobID(ctx, jobID)
	if err != nil {
		return err
	}

	newStatus := domain.AssetStatusReady
	eventType := "VideoProcessed"
	if status == "failed" {
		newStatus = domain.AssetStatusFailed
		eventType = "VideoProcessingFailed"
	}

	if err := u.repo.UpdateAssetStatus(ctx, asset.ID, newStatus, cdnURLs); err != nil {
		return err
	}

	asset.Status = newStatus
	asset.CDNURLs = cdnURLs

	payload, _ := json.Marshal(asset)
	u.producer.PublishAsync(kafka.Event{
		EventType:  eventType,
		AssetID:    asset.ID.String(),
		LessonID:   asset.LessonID.String(),
		Payload:    payload,
		OccurredAt: time.Now(),
	})

	return nil
}

func (u *MediaUsecase) GetAsset(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	return u.repo.GetAsset(ctx, id)
}
