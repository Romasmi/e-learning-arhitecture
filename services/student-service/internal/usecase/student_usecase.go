package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"log/slog"

	authapi "github.com/Romasmi/e-learning-arhitecture/gen/go/auth"
	"github.com/elearning/student-service/internal/domain"
	"github.com/elearning/student-service/pkg/kafka"
	"github.com/google/uuid"
)

type StudentUsecase struct {
	repo       domain.StudentRepository
	producer   *kafka.Producer
	authClient authapi.AuthServiceClient
}

func NewStudentUsecase(repo domain.StudentRepository, producer *kafka.Producer, authClient authapi.AuthServiceClient) *StudentUsecase {
	return &StudentUsecase{
		repo:       repo,
		producer:   producer,
		authClient: authClient,
	}
}

func (u *StudentUsecase) CreateStudent(ctx context.Context, accountID, email, password string) (*domain.Student, error) {
	// 1. Call auth-service to create credentials
	slog.Info("calling auth-service register", "email", email, "portal_id", accountID)
	resp, err := u.authClient.Register(ctx, &authapi.RegisterRequest{
		Email:    email,
		Password: password,
		PortalId: accountID, // Using accountID as portalID for students
		Role:     "student",
	})
	if err != nil {
		slog.Error("auth register failed", "error", err)
		return nil, fmt.Errorf("auth register: %w", err)
	}

	slog.Info("auth register success", "user_id", resp.UserId)

	// 2. Create student in DB
	student := &domain.Student{
		ID:        resp.UserId,
		AccountID: accountID,
		Email:     email,
		Status:    domain.StudentStatusActive,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateStudent(ctx, student); err != nil {
		slog.Error("repo create student failed", "error", err)
		return nil, fmt.Errorf("repo create student: %w", err)
	}

	// 3. Emit StudentCreated event
	u.emitEvent(ctx, "StudentCreated", accountID, student)

	return student, nil
}

func (u *StudentUsecase) DeactivateStudent(ctx context.Context, id string) (bool, error) {
	student, err := u.repo.GetStudent(ctx, id)
	if err != nil {
		return false, err
	}

	if err := u.repo.UpdateStudentStatus(ctx, id, domain.StudentStatusInactive); err != nil {
		return false, err
	}

	u.emitEvent(ctx, "StudentDeactivated", student.AccountID, map[string]string{"student_id": id})

	return true, nil
}

func (u *StudentUsecase) GetStudent(ctx context.Context, id string) (*domain.Student, error) {
	return u.repo.GetStudent(ctx, id)
}

func (u *StudentUsecase) ListStudents(ctx context.Context, accountID string) ([]*domain.Student, error) {
	return u.repo.ListStudents(ctx, accountID)
}

func (u *StudentUsecase) AssignLicense(ctx context.Context, studentID, licenseID string) (bool, error) {
	student, err := u.repo.GetStudent(ctx, studentID)
	if err != nil {
		return false, err
	}

	if err := u.repo.UpdateStudentLicense(ctx, studentID, licenseID); err != nil {
		return false, err
	}

	u.emitEvent(ctx, "LicenseAssigned", student.AccountID, map[string]string{
		"student_id": studentID,
		"license_id": licenseID,
	})

	return true, nil
}

func (u *StudentUsecase) AddGroup(ctx context.Context, accountID, name, parentID string) (*domain.StudentGroup, error) {
	group := &domain.StudentGroup{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	if parentID != "" {
		group.ParentID = &parentID
	}

	if err := u.repo.CreateGroup(ctx, group); err != nil {
		return nil, err
	}

	u.emitEvent(ctx, "GroupCreated", accountID, group)

	return group, nil
}

func (u *StudentUsecase) UpdateGroup(ctx context.Context, groupID, name, parentID string) (*domain.StudentGroup, error) {
	group, err := u.repo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	group.Name = name
	if parentID != "" {
		group.ParentID = &parentID
	} else {
		group.ParentID = nil
	}

	if err := u.repo.UpdateGroup(ctx, group); err != nil {
		return nil, err
	}

	u.emitEvent(ctx, "GroupUpdated", group.AccountID, group)

	return group, nil
}

func (u *StudentUsecase) DeleteGroup(ctx context.Context, groupID string) (bool, error) {
	group, err := u.repo.GetGroup(ctx, groupID)
	if err != nil {
		return false, err
	}

	if err := u.repo.DeleteGroup(ctx, groupID); err != nil {
		return false, err
	}

	u.emitEvent(ctx, "GroupDeleted", group.AccountID, map[string]string{"group_id": groupID})

	return true, nil
}

func (u *StudentUsecase) AddStudentToGroup(ctx context.Context, studentID, groupID string) (bool, error) {
	student, err := u.repo.GetStudent(ctx, studentID)
	if err != nil {
		return false, err
	}

	if err := u.repo.AddStudentToGroup(ctx, studentID, groupID); err != nil {
		return false, err
	}

	u.emitEvent(ctx, "StudentAddedToGroup", student.AccountID, map[string]string{
		"student_id": studentID,
		"group_id":   groupID,
	})

	return true, nil
}

func (u *StudentUsecase) ListGroups(ctx context.Context, accountID string) ([]*domain.StudentGroup, error) {
	return u.repo.ListGroups(ctx, accountID)
}

func (u *StudentUsecase) emitEvent(ctx context.Context, eventType, accountID string, payload any) {
	payloadJSON, _ := json.Marshal(payload)
	u.producer.PublishAsync(kafka.Event{
		EventType:  eventType,
		AccountID:  accountID,
		Payload:    payloadJSON,
		OccurredAt: time.Now(),
	})
}
