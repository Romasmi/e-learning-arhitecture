package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/e-learning-arhitecture/user-service/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate")
)

type UserRepository struct {
	db DBQuerier
}

const usersTable = "users"

func CreateUserRepository(db DBQuerier) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, userModel *user.User) (*user.User, error) {
	userId, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %w", err)
	}

	const query = `
		INSERT INTO %s (id, username, firstname, lastname, email, phone)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, firstname, lastname, email, phone, created_at
	`
	sql := fmt.Sprintf(query, usersTable)

	newUser := &user.User{}
	err = r.db.QueryRow(ctx, sql, userId, userModel.Username, userModel.FirstName, userModel.LastName, userModel.Email, userModel.Phone).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Email,
		&newUser.Phone,
		&newUser.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrDuplicate
			}
			return nil, err
		}
		return nil, fmt.Errorf("failed to create userModel: %w", err)
	}
	return newUser, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*user.User, error) {
	const query = `
		SELECT id, username, firstname, lastname, email, phone, created_at
        FROM %s
		WHERE id = $1
	`
	sql := fmt.Sprintf(query, usersTable)

	userModel := &user.User{}
	err := r.db.QueryRow(ctx, sql, userId).Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Email,
		&userModel.Phone,
		&userModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get userModel: %w", err)
	}
	return userModel, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userModel *user.User) (*user.User, error) {
	current, err := r.GetUserById(ctx, userModel.ID)
	if err != nil {
		return nil, err
	}

	if userModel.Username != "" {
		current.Username = userModel.Username
	}
	if userModel.FirstName != "" {
		current.FirstName = userModel.FirstName
	}
	if userModel.LastName != "" {
		current.LastName = userModel.LastName
	}
	if userModel.Email != "" {
		current.Email = userModel.Email
	}
	if userModel.Phone != "" {
		current.Phone = userModel.Phone
	}

	const query = `
		UPDATE %s
		SET username = $2, firstname = $3, lastname = $4, email = $5, phone = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, firstname, lastname, email, phone, created_at
	`
	sql := fmt.Sprintf(query, usersTable)

	updatedUser := &user.User{}
	err = r.db.QueryRow(ctx, sql, current.ID, current.Username, current.FirstName, current.LastName, current.Email, current.Phone).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Email,
		&updatedUser.Phone,
		&updatedUser.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrDuplicate
			}
			return nil, err
		}
		return nil, fmt.Errorf("failed to update userModel: %w", err)
	}
	return updatedUser, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	const query = `
		DELETE FROM %s
		WHERE id = $1
	`
	sql := fmt.Sprintf(query, usersTable)

	result, err := r.db.Exec(ctx, sql, userId)
	if err != nil {
		return fmt.Errorf("failed to delete userModel: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
