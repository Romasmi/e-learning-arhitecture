package handler

import (
	"context"

	coursepb "github.com/Romasmi/e-learning-arhitecture/gen/go/course"
	"github.com/elearning/course-service/internal/domain"
	"github.com/elearning/course-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	coursepb.UnimplementedCourseServiceServer
	usecase *usecase.CourseUsecase
}

func NewGRPCHandler(u *usecase.CourseUsecase) *GRPCHandler {
	return &GRPCHandler{usecase: u}
}

func (h *GRPCHandler) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CreateCourseResponse, error) {
	course, err := h.usecase.CreateCourse(ctx, req.PortalId, req.Title, req.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create course: %v", err)
	}

	return &coursepb.CreateCourseResponse{
		Course: mapCourseToProto(course),
	}, nil
}

func (h *GRPCHandler) UpdateCourse(ctx context.Context, req *coursepb.UpdateCourseRequest) (*coursepb.UpdateCourseResponse, error) {
	course, err := h.usecase.UpdateCourse(ctx, req.CourseId, req.Title, req.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update course: %v", err)
	}

	return &coursepb.UpdateCourseResponse{
		Course: mapCourseToProto(course),
	}, nil
}

func (h *GRPCHandler) ArchiveCourse(ctx context.Context, req *coursepb.ArchiveCourseRequest) (*coursepb.ArchiveCourseResponse, error) {
	archived, err := h.usecase.ArchiveCourse(ctx, req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to archive course: %v", err)
	}

	return &coursepb.ArchiveCourseResponse{
		Archived: archived,
	}, nil
}

func (h *GRPCHandler) AddChapter(ctx context.Context, req *coursepb.AddChapterRequest) (*coursepb.AddChapterResponse, error) {
	chapter, err := h.usecase.AddChapter(ctx, req.CourseId, req.Title, req.Position)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add chapter: %v", err)
	}

	return &coursepb.AddChapterResponse{
		Chapter: mapChapterToProto(chapter),
	}, nil
}

func (h *GRPCHandler) AddLesson(ctx context.Context, req *coursepb.AddLessonRequest) (*coursepb.AddLessonResponse, error) {
	lesson, err := h.usecase.AddLesson(ctx, req.ChapterId, req.Title, req.Type.String(), req.AssetId, req.Position)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add lesson: %v", err)
	}

	return &coursepb.AddLessonResponse{
		Lesson: mapLessonToProto(lesson),
	}, nil
}

func (h *GRPCHandler) AttachCertTest(ctx context.Context, req *coursepb.AttachCertTestRequest) (*coursepb.AttachCertTestResponse, error) {
	questions := make([]domain.Question, len(req.Questions))
	for i, q := range req.Questions {
		questions[i] = domain.Question{
			Text:         q.Text,
			Options:      q.Options,
			CorrectIndex: q.CorrectIndex,
		}
	}

	certTest, err := h.usecase.AttachCertTest(ctx, req.CourseId, req.PassScore, questions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to attach cert test: %v", err)
	}

	return &coursepb.AttachCertTestResponse{
		CertTest: mapCertTestToProto(certTest),
	}, nil
}

func (h *GRPCHandler) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.GetCourseResponse, error) {
	course, err := h.usecase.GetCourse(ctx, req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get course: %v", err)
	}

	return &coursepb.GetCourseResponse{
		Course: mapCourseToProto(course),
	}, nil
}

func (h *GRPCHandler) ListCourses(ctx context.Context, req *coursepb.ListCoursesRequest) (*coursepb.ListCoursesResponse, error) {
	courses, err := h.usecase.ListCourses(ctx, req.PortalId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list courses: %v", err)
	}

	pbCourses := make([]*coursepb.Course, len(courses))
	for i, c := range courses {
		pbCourses[i] = mapCourseToProto(c)
	}

	return &coursepb.ListCoursesResponse{
		Courses: pbCourses,
	}, nil
}

func mapCourseToProto(c *domain.Course) *coursepb.Course {
	pbChapters := make([]*coursepb.Chapter, len(c.Chapters))
	for i, ch := range c.Chapters {
		pbChapters[i] = mapChapterToProto(&ch)
	}

	var pbCertTest *coursepb.CertTest
	if c.CertTest != nil {
		pbCertTest = mapCertTestToProto(c.CertTest)
	}

	return &coursepb.Course{
		Id:          c.ID,
		PortalId:    c.PortalID,
		Title:       c.Title,
		Description: c.Description,
		Status:      coursepb.CourseStatus(coursepb.CourseStatus_value[string(c.Status)]),
		Chapters:    pbChapters,
		CertTest:    pbCertTest,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}

func mapChapterToProto(ch *domain.Chapter) *coursepb.Chapter {
	pbLessons := make([]*coursepb.Lesson, len(ch.Lessons))
	for i, l := range ch.Lessons {
		pbLessons[i] = mapLessonToProto(&l)
	}

	return &coursepb.Chapter{
		Id:       ch.ID,
		CourseId: ch.CourseID,
		Title:    ch.Title,
		Position: ch.Position,
		Lessons:  pbLessons,
	}
}

func mapLessonToProto(l *domain.Lesson) *coursepb.Lesson {
	return &coursepb.Lesson{
		Id:          l.ID,
		ChapterId:   l.ChapterID,
		Title:       l.Title,
		Type:        coursepb.LessonType(coursepb.LessonType_value[l.Type]),
		AssetId:     l.AssetID,
		Position:    l.Position,
		AssetStatus: l.AssetStatus,
	}
}

func mapCertTestToProto(ct *domain.CertTest) *coursepb.CertTest {
	pbQuestions := make([]*coursepb.Question, len(ct.Questions))
	for i, q := range ct.Questions {
		pbQuestions[i] = &coursepb.Question{
			Id:           q.ID,
			Text:         q.Text,
			Options:      q.Options,
			CorrectIndex: q.CorrectIndex,
		}
	}

	return &coursepb.CertTest{
		Id:        ct.ID,
		CourseId:  ct.CourseID,
		PassScore: ct.PassScore,
		Questions: pbQuestions,
	}
}
