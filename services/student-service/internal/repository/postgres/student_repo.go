package postgres

import (
	"context"
	"fmt"

	"github.com/elearning/student-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type studentRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) domain.StudentRepository {
	return &studentRepository{pool: pool}
}

func (r *studentRepository) CreateStudent(ctx context.Context, s *domain.Student) error {
	query := `
		INSERT INTO students (id, account_id, email, status, license_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query, s.ID, s.AccountID, s.Email, s.Status, s.LicenseID, s.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec insert student: %w", err)
	}
	return nil
}

func (r *studentRepository) GetStudent(ctx context.Context, id string) (*domain.Student, error) {
	query := `SELECT id, account_id, email, status, license_id, created_at FROM students WHERE id = $1`
	var s domain.Student
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.AccountID, &s.Email, &s.Status, &s.LicenseID, &s.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, fmt.Errorf("query row student: %w", err)
	}
	return &s, nil
}

func (r *studentRepository) ListStudents(ctx context.Context, accountID string) ([]*domain.Student, error) {
	query := `SELECT id, account_id, email, status, license_id, created_at FROM students WHERE account_id = $1`
	rows, err := r.pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("query students: %w", err)
	}
	defer rows.Close()

	var students []*domain.Student
	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(&s.ID, &s.AccountID, &s.Email, &s.Status, &s.LicenseID, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan student: %w", err)
		}
		students = append(students, &s)
	}
	return students, nil
}

func (r *studentRepository) UpdateStudentStatus(ctx context.Context, id string, status domain.StudentStatus) error {
	query := `UPDATE students SET status = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("exec update student status: %w", err)
	}
	return nil
}

func (r *studentRepository) UpdateStudentLicense(ctx context.Context, id string, licenseID string) error {
	query := `UPDATE students SET license_id = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, licenseID, id)
	if err != nil {
		return fmt.Errorf("exec update student license: %w", err)
	}
	return nil
}

func (r *studentRepository) CreateGroup(ctx context.Context, g *domain.StudentGroup) error {
	query := `
		INSERT INTO groups (id, account_id, name, parent_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, g.ID, g.AccountID, g.Name, g.ParentID, g.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec insert group: %w", err)
	}
	return nil
}

func (r *studentRepository) UpdateGroup(ctx context.Context, g *domain.StudentGroup) error {
	query := `UPDATE groups SET name = $1, parent_id = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, g.Name, g.ParentID, g.ID)
	if err != nil {
		return fmt.Errorf("exec update group: %w", err)
	}
	return nil
}

func (r *studentRepository) DeleteGroup(ctx context.Context, id string) error {
	query := `DELETE FROM groups WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec delete group: %w", err)
	}
	return nil
}

func (r *studentRepository) GetGroup(ctx context.Context, id string) (*domain.StudentGroup, error) {
	query := `SELECT id, account_id, name, parent_id, created_at FROM groups WHERE id = $1`
	var g domain.StudentGroup
	err := r.pool.QueryRow(ctx, query, id).Scan(&g.ID, &g.AccountID, &g.Name, &g.ParentID, &g.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("query row group: %w", err)
	}
	return &g, nil
}

func (r *studentRepository) ListGroups(ctx context.Context, accountID string) ([]*domain.StudentGroup, error) {
	query := `SELECT id, account_id, name, parent_id, created_at FROM groups WHERE account_id = $1`
	rows, err := r.pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("query groups: %w", err)
	}
	defer rows.Close()

	var groups []*domain.StudentGroup
	for rows.Next() {
		var g domain.StudentGroup
		if err := rows.Scan(&g.ID, &g.AccountID, &g.Name, &g.ParentID, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}
		groups = append(groups, &g)
	}
	return groups, nil
}

func (r *studentRepository) AddStudentToGroup(ctx context.Context, studentID, groupID string) error {
	query := `INSERT INTO student_group_members (student_id, group_id) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, studentID, groupID)
	if err != nil {
		return fmt.Errorf("exec insert student_group_member: %w", err)
	}
	return nil
}
