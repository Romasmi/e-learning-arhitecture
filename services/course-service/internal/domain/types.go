package domain

import (
	"time"
)

type CourseStatus string

const (
	CourseStatusDraft    CourseStatus = "DRAFT"
	CourseStatusArchived CourseStatus = "ARCHIVED"
)

type Course struct {
	ID          string       `json:"id"`
	PortalID    string       `json:"portal_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      CourseStatus `json:"status"`
	Chapters    []Chapter    `json:"chapters,omitempty"`
	CertTest    *CertTest    `json:"cert_test,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Chapter struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Title     string    `json:"title"`
	Position  int32     `json:"position"`
	Lessons   []Lesson  `json:"lessons,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Lesson struct {
	ID          string    `json:"id"`
	ChapterID   string    `json:"chapter_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	AssetID     string    `json:"asset_id"`
	AssetStatus string    `json:"asset_status"`
	Position    int32     `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
}

type CertTest struct {
	ID        string     `json:"id"`
	CourseID  string     `json:"course_id"`
	PassScore int32      `json:"pass_score"`
	Questions []Question `json:"questions,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Question struct {
	ID           string   `json:"id"`
	CertTestID   string   `json:"cert_test_id"`
	Text         string   `json:"text"`
	Options      []string `json:"options"`
	CorrectIndex int32    `json:"correct_index"`
}

type CourseEvent struct {
	EventType  string    `json:"event_type"`
	CourseID   string    `json:"course_id"`
	PortalID   string    `json:"portal_id"`
	Payload    any       `json:"payload"`
	OccurredAt time.Time `json:"occurred_at"`
}
