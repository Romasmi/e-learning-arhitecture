package domain

import (
	"context"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *Course) error
	UpdateCourse(ctx context.Context, course *Course) error
	GetCourse(ctx context.Context, id string) (*Course, error)
	ListCourses(ctx context.Context, portalID string) ([]*Course, error)

	AddChapter(ctx context.Context, chapter *Chapter) error
	AddLesson(ctx context.Context, lesson *Lesson) error

	UpsertCertTest(ctx context.Context, certTest *CertTest) error

	UpdateLessonAssetStatus(ctx context.Context, assetID string, status string) error
}
