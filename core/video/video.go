package video

import "time"

// Video models videos.
// A course can contain many videos.
// A video can be contained by a course only.
// URL is not marhsalled to JSON to avoid security issues.
type Video struct {
	ID          string    `json:"id" db:"video_id"`
	CourseID    string    `json:"courseId" db:"course_id"`
	Index       int       `json:"index" db:"index"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Free        bool      `json:"free" db:"free"`
	URL         string    `json:"-" db:"url"`
	ImageURL    string    `json:"imageUrl" db:"image_url"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	Version     int       `json:"-" db:"version"`
}

// VideoNew contains all the information needed to insert a new video.
type VideoNew struct {
	CourseID    string `json:"courseId" validate:"required"`
	Index       int    `json:"index" validate:"required,gte=0"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Free        bool   `json:"free" validate:"required"`
	URL         string `json:"url" validate:"omitempty,url"`
	ImageURL    string `json:"imageUrl" validate:"required"`
}

// VideoUp specifies the data of videos that can be updated.
type VideoUp struct {
	CourseID    *string `json:"courseId"`
	Index       *int    `json:"index" validate:"omitempty,gte=0"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Free        *bool   `json:"free"`
	URL         *string `json:"url" validate:"omitempty,url"`
	ImageURL    *string `json:"imageUrl"`
}

// Progress models users' progress on videos.
type Progress struct {
	VideoID   string    `json:"videoId" db:"video_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Progress  int       `json:"progress" db:"progress"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// ProgressUp contains the data of a progress which can be updated.
type ProgressUp struct {
	Progress int `json:"progress" validate:"gte=0,lte=100"`
}
