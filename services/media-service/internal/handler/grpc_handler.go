package handler

import (
	"context"

	mediapb "github.com/Romasmi/e-learning-arhitecture/gen/go/media"
	"github.com/elearning/media-service/internal/domain"
	"github.com/elearning/media-service/internal/usecase"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	mediapb.UnimplementedMediaServiceServer
	uc *usecase.MediaUsecase
}

func NewGRPCHandler(uc *usecase.MediaUsecase) *GRPCHandler {
	return &GRPCHandler{uc: uc}
}

func (h *GRPCHandler) UploadVideo(ctx context.Context, req *mediapb.UploadVideoRequest) (*mediapb.UploadVideoResponse, error) {
	lessonID, err := uuid.Parse(req.LessonId)
	if err != nil {
		return nil, err
	}
	asset, err := h.uc.UploadAsset(ctx, lessonID, req.Content, req.Filename, domain.AssetTypeVideo)
	if err != nil {
		return nil, err
	}
	return &mediapb.UploadVideoResponse{Asset: toProtoAsset(asset)}, nil
}

func (h *GRPCHandler) UploadPDF(ctx context.Context, req *mediapb.UploadPDFRequest) (*mediapb.UploadPDFResponse, error) {
	lessonID, err := uuid.Parse(req.LessonId)
	if err != nil {
		return nil, err
	}
	asset, err := h.uc.UploadAsset(ctx, lessonID, req.Content, req.Filename, domain.AssetTypePDF)
	if err != nil {
		return nil, err
	}
	return &mediapb.UploadPDFResponse{Asset: toProtoAsset(asset)}, nil
}

func (h *GRPCHandler) UploadImage(ctx context.Context, req *mediapb.UploadImageRequest) (*mediapb.UploadImageResponse, error) {
	lessonID, err := uuid.Parse(req.LessonId)
	if err != nil {
		return nil, err
	}
	asset, err := h.uc.UploadAsset(ctx, lessonID, req.Content, req.Filename, domain.AssetTypeImage)
	if err != nil {
		return nil, err
	}
	return &mediapb.UploadImageResponse{Asset: toProtoAsset(asset)}, nil
}

func (h *GRPCHandler) GetAsset(ctx context.Context, req *mediapb.GetAssetRequest) (*mediapb.GetAssetResponse, error) {
	assetID, err := uuid.Parse(req.AssetId)
	if err != nil {
		return nil, err
	}
	asset, err := h.uc.GetAsset(ctx, assetID)
	if err != nil {
		return nil, err
	}
	return &mediapb.GetAssetResponse{Asset: toProtoAsset(asset)}, nil
}

func (h *GRPCHandler) TranscodeCallback(ctx context.Context, req *mediapb.TranscodeCallbackRequest) (*mediapb.TranscodeCallbackResponse, error) {
	err := h.uc.HandleTranscodeCallback(ctx, req.JobId, req.Status, req.CdnUrls)
	if err != nil {
		return nil, err
	}
	return &mediapb.TranscodeCallbackResponse{Ok: true}, nil
}

func toProtoAsset(a *domain.Asset) *mediapb.Asset {
	return &mediapb.Asset{
		Id:        a.ID.String(),
		LessonId:  a.LessonID.String(),
		Type:      mapType(a.Type),
		Status:    mapStatus(a.Status),
		RawUrl:    a.RawURL,
		CdnUrls:   a.CDNURLs,
		CreatedAt: timestamppb.New(a.CreatedAt),
	}
}

func mapType(t domain.AssetType) mediapb.AssetType {
	switch t {
	case domain.AssetTypeVideo:
		return mediapb.AssetType_VIDEO
	case domain.AssetTypePDF:
		return mediapb.AssetType_PDF
	case domain.AssetTypeImage:
		return mediapb.AssetType_IMAGE
	default:
		return mediapb.AssetType_ASSET_TYPE_UNSPECIFIED
	}
}

func mapStatus(s domain.AssetStatus) mediapb.AssetStatus {
	switch s {
	case domain.AssetStatusPending:
		return mediapb.AssetStatus_PENDING
	case domain.AssetStatusReady:
		return mediapb.AssetStatus_READY
	case domain.AssetStatusFailed:
		return mediapb.AssetStatus_FAILED
	default:
		return mediapb.AssetStatus_ASSET_STATUS_UNSPECIFIED
	}
}
