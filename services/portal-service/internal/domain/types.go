package domain

import (
	"errors"
	"time"
)

type PortalStatus string

const (
	PortalStatusActive   PortalStatus = "ACTIVE"
	PortalStatusArchived PortalStatus = "ARCHIVED"
)

type LMSConfig struct {
	ThemeColor        string `json:"theme_color"`
	LogoURL           string `json:"logo_url"`
	EnableSocialLogin bool   `json:"enable_social_login"`
}

type Portal struct {
	ID        string       `json:"id"`
	Code      string       `json:"code"`
	Name      string       `json:"name"`
	Status    PortalStatus `json:"status"`
	LMSConfig LMSConfig    `json:"lms_config"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (p *Portal) DomainURL() string {
	if p.Code == "" {
		return ""
	}
	return p.Code + ".e-learning.com"
}

var (
	ErrPortalNotFound = errors.New("portal not found")
	ErrDuplicateCode  = errors.New("portal code already exists")
)
