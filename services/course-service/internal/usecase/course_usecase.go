package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/elearning/course-service/internal/domain"
	"github.com/elearning/course-service/pkg/kafka"
)

type CourseUsecase struct {
	repo     domain.CourseRepository
	producer *kafka.Producer
}

func NewCourseUsecase(repo domain.CourseRepository, producer *kafka.Producer) *CourseUsecase {
	return &CourseUsecase{
		repo:     repo,
		producer: producer,
	}
}

func (u *CourseUsecase) CreateCourse(ctx context.Context, portalID, title, description string) (*domain.Course, error) {
	course := &domain.Course{
		PortalID:    portalID,
		Title:       title,
		Description: description,
		Status:      domain.CourseStatusDraft,
	}

	if err := u.repo.CreateCourse(ctx, course); err != nil {
		return nil, err
	}

	u.emitEvent(ctx, "CourseCreated", course)

	return course, nil
}

func (u *CourseUsecase) UpdateCourse(ctx context.Context, courseID, title, description string) (*domain.Course, error) {
	course, err := u.repo.GetCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}

	course.Title = title
	course.Description = description

	if err := u.repo.UpdateCourse(ctx, course); err != nil {
		return nil, err
	}

	u.emitEvent(ctx, "CourseUpdated", course)

	return course, nil
}

func (u *CourseUsecase) ArchiveCourse(ctx context.Context, courseID string) (bool, error) {
	course, err := u.repo.GetCourse(ctx, courseID)
	if err != nil {
		return false, err
	}

	course.Status = domain.CourseStatusArchived

	if err := u.repo.UpdateCourse(ctx, course); err != nil {
		return false, err
	}

	u.emitEvent(ctx, "CourseArchived", course)

	return true, nil
}

func (u *CourseUsecase) AddChapter(ctx context.Context, courseID, title string, position int32) (*domain.Chapter, error) {
	chapter := &domain.Chapter{
		CourseID: courseID,
		Title:    title,
		Position: position,
	}

	if err := u.repo.AddChapter(ctx, chapter); err != nil {
		return nil, err
	}

	return chapter, nil
}

func (u *CourseUsecase) AddLesson(ctx context.Context, chapterID, title, lType, assetID string, position int32) (*domain.Lesson, error) {
	lesson := &domain.Lesson{
		ChapterID: chapterID,
		Title:     title,
		Type:      lType,
		AssetID:   assetID,
		Position:  position,
	}

	if err := u.repo.AddLesson(ctx, lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

func (u *CourseUsecase) AttachCertTest(ctx context.Context, courseID string, passScore int32, questions []domain.Question) (*domain.CertTest, error) {
	certTest := &domain.CertTest{
		CourseID:  courseID,
		PassScore: passScore,
		Questions: questions,
	}

	if err := u.repo.UpsertCertTest(ctx, certTest); err != nil {
		return nil, err
	}

	return certTest, nil
}

func (u *CourseUsecase) GetCourse(ctx context.Context, courseID string) (*domain.Course, error) {
	return u.repo.GetCourse(ctx, courseID)
}

func (u *CourseUsecase) ListCourses(ctx context.Context, portalID string) ([]*domain.Course, error) {
	return u.repo.ListCourses(ctx, portalID)
}

func (u *CourseUsecase) HandleVideoProcessed(ctx context.Context, assetID, status string) error {
	return u.repo.UpdateLessonAssetStatus(ctx, assetID, status)
}

func (u *CourseUsecase) emitEvent(ctx context.Context, eventType string, course *domain.Course) {
	payload, _ := json.Marshal(course)
	u.producer.PublishAsync(kafka.Event{
		EventType:  eventType,
		CourseID:   course.ID,
		PortalID:   course.PortalID,
		Payload:    payload,
		OccurredAt: time.Now(),
	})
}
