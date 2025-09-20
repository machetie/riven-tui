package models

// StatsResponse represents the system statistics response
type StatsResponse struct {
	TotalItems         int            `json:"total_items"`
	TotalMovies        int            `json:"total_movies"`
	TotalShows         int            `json:"total_shows"`
	TotalSeasons       int            `json:"total_seasons"`
	TotalEpisodes      int            `json:"total_episodes"`
	TotalSymlinks      int            `json:"total_symlinks"`
	IncompleteItems    int            `json:"incomplete_items"`
	IncompleteRetries  map[string]int `json:"incomplete_retries"` // Media item log string: number of retries
	States             map[States]int `json:"states"`
}

// StateResponse represents the available states response
type StateResponse struct {
	Success bool     `json:"success"`
	States  []string `json:"states"`
}

// ServicesResponse represents the services status response
type ServicesResponse map[string]bool

// RDUser represents Real-Debrid user information
type RDUser struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Points   int      `json:"points"`   // User's RD points
	Locale   string   `json:"locale"`
	Avatar   string   `json:"avatar"`   // URL to the user's avatar
	Type     UserType `json:"type"`
	Premium  int      `json:"premium"`  // Premium subscription left in seconds
}

// LogsResponse represents the logs response
type LogsResponse struct {
	Logs []string `json:"logs"`
}

// EventResponse represents the events response
type EventResponse struct {
	Events map[string][]int `json:"events"`
}

// StreamEventResponse represents streaming event response
type StreamEventResponse struct {
	Data map[string]interface{} `json:"data"`
}

// MountResponse represents the mount files response
type MountResponse struct {
	Files map[string]string `json:"files"`
}

// UploadLogsResponse represents the upload logs response
type UploadLogsResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"` // URL to the uploaded log file. 50M Filesize limit. 180 day retention.
}

// CalendarResponse represents the calendar response
type CalendarResponse struct {
	// This would contain calendar data - structure not fully defined in the API spec
	Data interface{} `json:"data,omitempty"`
}

// TraktOAuthInitiateResponse represents Trakt OAuth initiation response
type TraktOAuthInitiateResponse struct {
	AuthURL string `json:"auth_url"`
}
