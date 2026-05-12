package handler

import (
	"context"

	studentpb "github.com/Romasmi/e-learning-arhitecture/gen/go/student"
	"github.com/elearning/student-service/internal/domain"
	"github.com/elearning/student-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	studentpb.UnimplementedStudentServiceServer
	usecase *usecase.StudentUsecase
}

func NewGRPCHandler(u *usecase.StudentUsecase) *GRPCHandler {
	return &GRPCHandler{usecase: u}
}

func (h *GRPCHandler) CreateStudent(ctx context.Context, req *studentpb.CreateStudentRequest) (*studentpb.CreateStudentResponse, error) {
	student, err := h.usecase.CreateStudent(ctx, req.AccountId, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create student: %v", err)
	}

	return &studentpb.CreateStudentResponse{
		Student: mapStudentToProto(student),
	}, nil
}

func (h *GRPCHandler) DeactivateStudent(ctx context.Context, req *studentpb.DeactivateStudentRequest) (*studentpb.DeactivateStudentResponse, error) {
	success, err := h.usecase.DeactivateStudent(ctx, req.StudentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to deactivate student: %v", err)
	}
	return &studentpb.DeactivateStudentResponse{Deactivated: success}, nil
}

func (h *GRPCHandler) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.GetStudentResponse, error) {
	student, err := h.usecase.GetStudent(ctx, req.StudentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "student not found: %v", err)
	}
	return &studentpb.GetStudentResponse{
		Student: mapStudentToProto(student),
	}, nil
}

func (h *GRPCHandler) ListStudents(ctx context.Context, req *studentpb.ListStudentsRequest) (*studentpb.ListStudentsResponse, error) {
	students, err := h.usecase.ListStudents(ctx, req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list students: %v", err)
	}

	pbStudents := make([]*studentpb.Student, len(students))
	for i, s := range students {
		pbStudents[i] = mapStudentToProto(s)
	}

	return &studentpb.ListStudentsResponse{Students: pbStudents}, nil
}

func (h *GRPCHandler) AssignLicense(ctx context.Context, req *studentpb.AssignLicenseRequest) (*studentpb.AssignLicenseResponse, error) {
	success, err := h.usecase.AssignLicense(ctx, req.StudentId, req.LicenseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to assign license: %v", err)
	}
	return &studentpb.AssignLicenseResponse{Assigned: success}, nil
}

func (h *GRPCHandler) AddGroup(ctx context.Context, req *studentpb.AddGroupRequest) (*studentpb.AddGroupResponse, error) {
	group, err := h.usecase.AddGroup(ctx, req.AccountId, req.Name, req.ParentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add group: %v", err)
	}
	return &studentpb.AddGroupResponse{
		Group: mapGroupToProto(group),
	}, nil
}

func (h *GRPCHandler) UpdateGroup(ctx context.Context, req *studentpb.UpdateGroupRequest) (*studentpb.UpdateGroupResponse, error) {
	group, err := h.usecase.UpdateGroup(ctx, req.GroupId, req.Name, req.ParentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update group: %v", err)
	}
	return &studentpb.UpdateGroupResponse{
		Group: mapGroupToProto(group),
	}, nil
}

func (h *GRPCHandler) DeleteGroup(ctx context.Context, req *studentpb.DeleteGroupRequest) (*studentpb.DeleteGroupResponse, error) {
	success, err := h.usecase.DeleteGroup(ctx, req.GroupId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete group: %v", err)
	}
	return &studentpb.DeleteGroupResponse{Deleted: success}, nil
}

func (h *GRPCHandler) AddStudentToGroup(ctx context.Context, req *studentpb.AddStudentToGroupRequest) (*studentpb.AddStudentToGroupResponse, error) {
	success, err := h.usecase.AddStudentToGroup(ctx, req.StudentId, req.GroupId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add student to group: %v", err)
	}
	return &studentpb.AddStudentToGroupResponse{Added: success}, nil
}

func (h *GRPCHandler) ListGroups(ctx context.Context, req *studentpb.ListGroupsRequest) (*studentpb.ListGroupsResponse, error) {
	groups, err := h.usecase.ListGroups(ctx, req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list groups: %v", err)
	}

	pbGroups := make([]*studentpb.StudentGroup, len(groups))
	for i, g := range groups {
		pbGroups[i] = mapGroupToProto(g)
	}

	return &studentpb.ListGroupsResponse{Groups: pbGroups}, nil
}

func mapStudentToProto(s *domain.Student) *studentpb.Student {
	licenseID := ""
	if s.LicenseID != nil {
		licenseID = *s.LicenseID
	}
	return &studentpb.Student{
		Id:        s.ID,
		AccountId: s.AccountID,
		Email:     s.Email,
		Status:    studentpb.StudentStatus(studentpb.StudentStatus_value[string(s.Status)]),
		LicenseId: licenseID,
		CreatedAt: timestamppb.New(s.CreatedAt),
	}
}

func mapGroupToProto(g *domain.StudentGroup) *studentpb.StudentGroup {
	parentID := ""
	if g.ParentID != nil {
		parentID = *g.ParentID
	}
	return &studentpb.StudentGroup{
		Id:        g.ID,
		AccountId: g.AccountID,
		Name:      g.Name,
		ParentId:  parentID,
		CreatedAt: timestamppb.New(g.CreatedAt),
	}
}
