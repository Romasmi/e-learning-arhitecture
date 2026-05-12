package domain

import (
	"context"
	"time"
)

type StudentStatus string

const (
	StudentStatusActive   StudentStatus = "ACTIVE"
	StudentStatusInactive StudentStatus = "INACTIVE"
)

type Student struct {
	ID        string        `json:"id"`
	AccountID string        `json:"account_id"`
	Email     string        `json:"email"`
	Status    StudentStatus `json:"status"`
	LicenseID *string       `json:"license_id,omitzero"`
	CreatedAt time.Time     `json:"created_at"`
}

type StudentGroup struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Name      string    `json:"name"`
	ParentID  *string   `json:"parent_id,omitzero"`
	CreatedAt time.Time `json:"created_at"`
}

type StudentRepository interface {
	CreateStudent(ctx context.Context, s *Student) error
	GetStudent(ctx context.Context, id string) (*Student, error)
	ListStudents(ctx context.Context, accountID string) ([]*Student, error)
	UpdateStudentStatus(ctx context.Context, id string, status StudentStatus) error
	UpdateStudentLicense(ctx context.Context, id string, licenseID string) error

	CreateGroup(ctx context.Context, g *StudentGroup) error
	UpdateGroup(ctx context.Context, g *StudentGroup) error
	DeleteGroup(ctx context.Context, id string) error
	GetGroup(ctx context.Context, id string) (*StudentGroup, error)
	ListGroups(ctx context.Context, accountID string) ([]*StudentGroup, error)
	AddStudentToGroup(ctx context.Context, studentID, groupID string) error
}
