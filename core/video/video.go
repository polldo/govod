package video

import "time"

type Video struct {
	ID          string    `json:"id" db:"video_id"`
	CourseID    string    `json:"course_id" db:"course_id"`
	Index       int       `json:"index" db:"index"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Free        bool      `json:"free" db:"free"`
	URL         string    `json:"-" db:"url"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Version     int       `json:"-" db:"version"`
}

type URL struct {
	VideoID string `json:"video_id" db:"video_id"`
	URL     string `json:"url" db:"url"`
}

type VideoNew struct {
	CourseID    string `json:"course_id" validate:"required"`
	Index       int    `json:"index" validate:"required,gte=0"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Free        bool   `json:"free" validate:"required"`
	URL         string `json:"url" validate:"omitempty,url"`
	ImageURL    string `json:"image_url" validate:"required"`
}

type VideoUp struct {
	CourseID    *string `json:"course_id"`
	Index       *int    `json:"index" validate:"omitempty,gte=0"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Free        *bool   `json:"free"`
	URL         *string `json:"url" validate:"omitempty,url"`
	ImageURL    *string `json:"image_url"`
}

type Progress struct {
	VideoID   string    `json:"video_id" db:"video_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Progress  int       `json:"progress" db:"progress"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ProgressUp struct {
	Progress int `json:"progress" validate:"gte=0,lte=100"`
}
