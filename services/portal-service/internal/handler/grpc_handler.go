package handler

import (
	"context"

	portalpb "github.com/Romasmi/e-learning-arhitecture/gen/go/portal"
	"github.com/elearning/portal-service/internal/domain"
	"github.com/elearning/portal-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	portalpb.UnimplementedPortalServiceServer
	usecase *usecase.PortalUsecase
}

func NewGRPCHandler(u *usecase.PortalUsecase) *GRPCHandler {
	return &GRPCHandler{usecase: u}
}

func (h *GRPCHandler) CreatePortal(ctx context.Context, req *portalpb.CreatePortalRequest) (*portalpb.CreatePortalResponse, error) {
	config := domain.LMSConfig{}
	if req.LmsConfig != nil {
		config.ThemeColor = req.LmsConfig.ThemeColor
		config.LogoURL = req.LmsConfig.LogoUrl
		config.EnableSocialLogin = req.LmsConfig.EnableSocialLogin
	}

	portal, err := h.usecase.CreatePortal(ctx, req.Code, req.Name, config)
	if err != nil {
		if err == domain.ErrDuplicateCode {
			return nil, status.Error(codes.AlreadyExists, "portal code already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create portal: %v", err)
	}

	return &portalpb.CreatePortalResponse{
		Portal: mapDomainToProto(portal),
	}, nil
}

func (h *GRPCHandler) GetPortal(ctx context.Context, req *portalpb.GetPortalRequest) (*portalpb.GetPortalResponse, error) {
	portal, err := h.usecase.GetPortal(ctx, req.Id)
	if err != nil {
		if err == domain.ErrPortalNotFound {
			return nil, status.Error(codes.NotFound, "portal not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get portal: %v", err)
	}

	return &portalpb.GetPortalResponse{
		Portal: mapDomainToProto(portal),
	}, nil
}

func (h *GRPCHandler) UpdatePortalConfig(ctx context.Context, req *portalpb.UpdatePortalConfigRequest) (*portalpb.UpdatePortalConfigResponse, error) {
	config := domain.LMSConfig{}
	if req.LmsConfig != nil {
		config.ThemeColor = req.LmsConfig.ThemeColor
		config.LogoURL = req.LmsConfig.LogoUrl
		config.EnableSocialLogin = req.LmsConfig.EnableSocialLogin
	}

	portal, err := h.usecase.UpdatePortalConfig(ctx, req.Id, req.Name, config)
	if err != nil {
		if err == domain.ErrPortalNotFound {
			return nil, status.Error(codes.NotFound, "portal not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update portal config: %v", err)
	}

	return &portalpb.UpdatePortalConfigResponse{
		Portal: mapDomainToProto(portal),
	}, nil
}

func (h *GRPCHandler) ArchivePortal(ctx context.Context, req *portalpb.ArchivePortalRequest) (*portalpb.ArchivePortalResponse, error) {
	err := h.usecase.ArchivePortal(ctx, req.Id)
	if err != nil {
		if err == domain.ErrPortalNotFound {
			return nil, status.Error(codes.NotFound, "portal not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to archive portal: %v", err)
	}

	return &portalpb.ArchivePortalResponse{Success: true}, nil
}

func (h *GRPCHandler) ListPortals(ctx context.Context, req *portalpb.ListPortalsRequest) (*portalpb.ListPortalsResponse, error) {
	portals, err := h.usecase.ListPortals(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list portals: %v", err)
	}

	pbPortals := make([]*portalpb.Portal, len(portals))
	for i, p := range portals {
		pbPortals[i] = mapDomainToProto(p)
	}

	return &portalpb.ListPortalsResponse{Portals: pbPortals}, nil
}

func mapDomainToProto(p *domain.Portal) *portalpb.Portal {
	return &portalpb.Portal{
		Id:        p.ID,
		Code:      p.Code,
		Name:      p.Name,
		Status:    string(p.Status),
		DomainUrl: p.DomainURL(),
		LmsConfig: &portalpb.LMSConfig{
			ThemeColor:        p.LMSConfig.ThemeColor,
			LogoUrl:           p.LMSConfig.LogoURL,
			EnableSocialLogin: p.LMSConfig.EnableSocialLogin,
		},
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
}
