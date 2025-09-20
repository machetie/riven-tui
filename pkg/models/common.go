package models

import "time"

// Common response structures

// MessageResponse represents a standard API message response
type MessageResponse struct {
	Message string   `json:"message"`
	TMDBIds []string `json:"tmdb_ids,omitempty"`
	TVDBIds []string `json:"tvdb_ids,omitempty"`
}

// RootResponse represents the root API endpoint response
type RootResponse struct {
	Message string   `json:"message"`
	Version string   `json:"version"`
	TMDBIds []string `json:"tmdb_ids,omitempty"`
	TVDBIds []string `json:"tvdb_ids,omitempty"`
}

// HTTPValidationError represents validation errors from the API
type HTTPValidationError struct {
	Detail []ValidationError `json:"detail"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Loc  []interface{} `json:"loc"`
	Msg  string        `json:"msg"`
	Type string        `json:"type"`
}

// States represents the possible states of media items
type States string

const (
	StateUnknown            States = "Unknown"
	StateUnreleased         States = "Unreleased"
	StateOngoing            States = "Ongoing"
	StateRequested          States = "Requested"
	StateIndexed            States = "Indexed"
	StateScraped            States = "Scraped"
	StateDownloaded         States = "Downloaded"
	StateSymlinked          States = "Symlinked"
	StateCompleted          States = "Completed"
	StatePartiallyCompleted States = "PartiallyCompleted"
	StateFailed             States = "Failed"
	StatePaused             States = "Paused"
)

// MediaType represents the type of media
type MediaType string

const (
	MediaTypeMovie MediaType = "movie"
	MediaTypeTV    MediaType = "tv"
	MediaTypeItem  MediaType = "item"
)

// SortOrder represents sorting options for media items
type SortOrder string

const (
	SortDateDesc  SortOrder = "date_desc"
	SortDateAsc   SortOrder = "date_asc"
	SortTitleAsc  SortOrder = "title_asc"
	SortTitleDesc SortOrder = "title_desc"
)

// UpdateType represents update types for new releases
type UpdateType string

const (
	UpdateTypeSeries   UpdateType = "series"
	UpdateTypeSeasons  UpdateType = "seasons"
	UpdateTypeEpisodes UpdateType = "episodes"
)

// DownloadType represents the type of download
type DownloadType string

const (
	DownloadTypeCached   DownloadType = "cached"
	DownloadTypeUncached DownloadType = "uncached"
)

// UserType represents the type of Real-Debrid user
type UserType string

const (
	UserTypeFree    UserType = "free"
	UserTypePremium UserType = "premium"
)

// Common utility functions

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns a pointer to the given bool
func BoolPtr(b bool) *bool {
	return &b
}

// TimePtr returns a pointer to the given time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// Event represents a real-time event from the server
type Event struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Message   string                 `json:"message,omitempty"`
}
