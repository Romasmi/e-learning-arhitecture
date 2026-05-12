package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elearning/course-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type courseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) domain.CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) CreateCourse(ctx context.Context, course *domain.Course) error {
	query := `INSERT INTO courses (portal_id, title, description, status) 
			  VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, course.PortalID, course.Title, course.Description, course.Status).
		Scan(&course.ID, &course.CreatedAt, &course.UpdatedAt)
}

func (r *courseRepository) UpdateCourse(ctx context.Context, course *domain.Course) error {
	query := `UPDATE courses SET title = $1, description = $2, status = $3, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $4 RETURNING updated_at`
	return r.db.QueryRow(ctx, query, course.Title, course.Description, course.Status, course.ID).
		Scan(&course.UpdatedAt)
}

func (r *courseRepository) GetCourse(ctx context.Context, id string) (*domain.Course, error) {
	course := &domain.Course{}
	query := `SELECT id, portal_id, title, description, status, created_at, updated_at FROM courses WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&course.ID, &course.PortalID, &course.Title, &course.Description, &course.Status, &course.CreatedAt, &course.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("course not found: %w", err)
		}
		return nil, err
	}

	// Fetch Chapters
	chapters, err := r.getChapters(ctx, id)
	if err != nil {
		return nil, err
	}
	course.Chapters = chapters

	// Fetch CertTest
	certTest, err := r.getCertTest(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	course.CertTest = certTest

	return course, nil
}

func (r *courseRepository) getChapters(ctx context.Context, courseID string) ([]domain.Chapter, error) {
	query := `SELECT id, course_id, title, position, created_at FROM chapters WHERE course_id = $1 ORDER BY position`
	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []domain.Chapter
	for rows.Next() {
		var ch domain.Chapter
		if err := rows.Scan(&ch.ID, &ch.CourseID, &ch.Title, &ch.Position, &ch.CreatedAt); err != nil {
			return nil, err
		}

		lessons, err := r.getLessons(ctx, ch.ID)
		if err != nil {
			return nil, err
		}
		ch.Lessons = lessons
		chapters = append(chapters, ch)
	}
	return chapters, nil
}

func (r *courseRepository) getLessons(ctx context.Context, chapterID string) ([]domain.Lesson, error) {
	query := `SELECT id, chapter_id, title, type, asset_id, asset_status, position, created_at FROM lessons WHERE chapter_id = $1 ORDER BY position`
	rows, err := r.db.Query(ctx, query, chapterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		var l domain.Lesson
		var assetID *string
		var assetStatus *string
		if err := rows.Scan(&l.ID, &l.ChapterID, &l.Title, &l.Type, &assetID, &assetStatus, &l.Position, &l.CreatedAt); err != nil {
			return nil, err
		}
		if assetID != nil {
			l.AssetID = *assetID
		}
		if assetStatus != nil {
			l.AssetStatus = *assetStatus
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *courseRepository) getCertTest(ctx context.Context, courseID string) (*domain.CertTest, error) {
	cert := &domain.CertTest{}
	query := `SELECT id, course_id, pass_score, created_at FROM cert_tests WHERE course_id = $1`
	err := r.db.QueryRow(ctx, query, courseID).Scan(&cert.ID, &cert.CourseID, &cert.PassScore, &cert.CreatedAt)
	if err != nil {
		return nil, err
	}

	questions, err := r.getQuestions(ctx, cert.ID)
	if err != nil {
		return nil, err
	}
	cert.Questions = questions
	return cert, nil
}

func (r *courseRepository) getQuestions(ctx context.Context, certTestID string) ([]domain.Question, error) {
	query := `SELECT id, cert_test_id, text, options, correct_index FROM questions WHERE cert_test_id = $1`
	rows, err := r.db.Query(ctx, query, certTestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []domain.Question
	for rows.Next() {
		var q domain.Question
		var optionsJSON []byte
		if err := rows.Scan(&q.ID, &q.CertTestID, &q.Text, &optionsJSON, &q.CorrectIndex); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(optionsJSON, &q.Options); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *courseRepository) ListCourses(ctx context.Context, portalID string) ([]*domain.Course, error) {
	query := `SELECT id, portal_id, title, description, status, created_at, updated_at FROM courses WHERE portal_id = $1`
	rows, err := r.db.Query(ctx, query, portalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*domain.Course
	for rows.Next() {
		c := &domain.Course{}
		if err := rows.Scan(&c.ID, &c.PortalID, &c.Title, &c.Description, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *courseRepository) AddChapter(ctx context.Context, chapter *domain.Chapter) error {
	query := `INSERT INTO chapters (course_id, title, position) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, chapter.CourseID, chapter.Title, chapter.Position).
		Scan(&chapter.ID, &chapter.CreatedAt)
}

func (r *courseRepository) AddLesson(ctx context.Context, lesson *domain.Lesson) error {
	query := `INSERT INTO lessons (chapter_id, title, type, asset_id, asset_status, position) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	var assetID *string
	if lesson.AssetID != "" {
		assetID = &lesson.AssetID
	}
	var assetStatus *string
	if lesson.AssetStatus != "" {
		assetStatus = &lesson.AssetStatus
	}
	return r.db.QueryRow(ctx, query, lesson.ChapterID, lesson.Title, lesson.Type, assetID, assetStatus, lesson.Position).
		Scan(&lesson.ID, &lesson.CreatedAt)
}

func (r *courseRepository) UpsertCertTest(ctx context.Context, certTest *domain.CertTest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete existing cert test and questions due to CASCADE
	_, err = tx.Exec(ctx, "DELETE FROM cert_tests WHERE course_id = $1", certTest.CourseID)
	if err != nil {
		return err
	}

	query := `INSERT INTO cert_tests (course_id, pass_score) VALUES ($1, $2) RETURNING id, created_at`
	err = tx.QueryRow(ctx, query, certTest.CourseID, certTest.PassScore).Scan(&certTest.ID, &certTest.CreatedAt)
	if err != nil {
		return err
	}

	for i := range certTest.Questions {
		q := &certTest.Questions[i]
		optionsJSON, err := json.Marshal(q.Options)
		if err != nil {
			return err
		}
		qQuery := `INSERT INTO questions (cert_test_id, text, options, correct_index) VALUES ($1, $2, $3, $4) RETURNING id`
		err = tx.QueryRow(ctx, qQuery, certTest.ID, q.Text, optionsJSON, q.CorrectIndex).Scan(&q.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *courseRepository) UpdateLessonAssetStatus(ctx context.Context, assetID string, status string) error {
	query := `UPDATE lessons SET asset_status = $1 WHERE asset_id = $2`
	_, err := r.db.Exec(ctx, query, status, assetID)
	return err
}
