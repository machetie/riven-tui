package models

import "time"

// ItemsResponse represents the paginated items response
type ItemsResponse struct {
	Success    bool                     `json:"success"`
	Items      []map[string]interface{} `json:"items"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalItems int                      `json:"total_items"`
	TotalPages int                      `json:"total_pages"`
}

// MediaItem represents a media item (simplified structure)
type MediaItem struct {
	ID           int                    `json:"id,omitempty"`
	Title        string                 `json:"title,omitempty"`
	Type         string                 `json:"type,omitempty"`
	State        States                 `json:"state,omitempty"`
	TMDBId       string                 `json:"tmdb_id,omitempty"`
	TVDBId       string                 `json:"tvdb_id,omitempty"`
	IMDBId       string                 `json:"imdb_id,omitempty"`
	Year         *int                   `json:"year,omitempty"`
	CreatedAt    *time.Time             `json:"created_at,omitempty"`
	UpdatedAt    *time.Time             `json:"updated_at,omitempty"`
	RequestedAt  *time.Time             `json:"requested_at,omitempty"`
	LastState    States                 `json:"last_state,omitempty"`
	Symlinked    bool                   `json:"symlinked,omitempty"`
	Blacklisted  bool                   `json:"blacklisted,omitempty"`
	Genres       []string               `json:"genres,omitempty"`
	Network      string                 `json:"network,omitempty"`
	Country      string                 `json:"country,omitempty"`
	Language     string                 `json:"language,omitempty"`
	Overview     string                 `json:"overview,omitempty"`
	Runtime      *int                   `json:"runtime,omitempty"`
	Seasons      []Season               `json:"seasons,omitempty"`
	Episodes     []Episode              `json:"episodes,omitempty"`
	Streams      []Stream               `json:"streams,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// Season represents a TV show season
type Season struct {
	ID          int       `json:"id,omitempty"`
	Number      int       `json:"number"`
	Title       string    `json:"title,omitempty"`
	State       States    `json:"state,omitempty"`
	Episodes    []Episode `json:"episodes,omitempty"`
	EpisodeCount int      `json:"episode_count,omitempty"`
	AirDate     *string   `json:"air_date,omitempty"`
}

// Episode represents a TV show episode
type Episode struct {
	ID          int     `json:"id,omitempty"`
	Number      int     `json:"number"`
	Title       string  `json:"title,omitempty"`
	State       States  `json:"state,omitempty"`
	AirDate     *string `json:"air_date,omitempty"`
	Runtime     *int    `json:"runtime,omitempty"`
	Overview    string  `json:"overview,omitempty"`
	SeasonNumber int    `json:"season_number,omitempty"`
	Streams     []Stream `json:"streams,omitempty"`
}

// Stream represents a torrent stream
type Stream struct {
	InfoHash     string     `json:"infohash"`
	RawTitle     string     `json:"raw_title"`
	ParsedTitle  string     `json:"parsed_title"`
	ParsedData   ParsedData `json:"parsed_data"`
	Rank         int        `json:"rank"`
	LevRatio     float64    `json:"lev_ratio"`
	IsCached     bool       `json:"is_cached"`
	Blacklisted  bool       `json:"blacklisted,omitempty"`
}

// ParsedData represents parsed torrent data
type ParsedData struct {
	RawTitle        string   `json:"raw_title"`
	ParsedTitle     string   `json:"parsed_title"`
	NormalizedTitle string   `json:"normalized_title"`
	Trash           bool     `json:"trash"`
	Adult           bool     `json:"adult"`
	Year            *int     `json:"year"`
	Resolution      string   `json:"resolution"`
	Seasons         []int    `json:"seasons"`
	Episodes        []int    `json:"episodes"`
	Complete        bool     `json:"complete"`
	Volumes         []int    `json:"volumes"`
	Languages       []string `json:"languages"`
	Quality         *string  `json:"quality"`
	HDR             []string `json:"hdr"`
	Codec           *string  `json:"codec"`
	Audio           []string `json:"audio"`
	Channels        []string `json:"channels"`
	Dubbed          bool     `json:"dubbed"`
	Subbed          bool     `json:"subbed"`
	Date            *string  `json:"date"`
	Group           *string  `json:"group"`
	Edition         *string  `json:"edition"`
	BitDepth        *string  `json:"bit_depth"`
	Bitrate         *string  `json:"bitrate"`
	Network         *string  `json:"network"`
	Extended        bool     `json:"extended"`
	Converted       bool     `json:"converted"`
	Hardcoded       bool     `json:"hardcoded"`
	Region          *string  `json:"region"`
	PPV             bool     `json:"ppv"`
	Site            *string  `json:"site"`
	Size            *string  `json:"size"`
	Proper          bool     `json:"proper"`
	Repack          bool     `json:"repack"`
	Retail          bool     `json:"retail"`
	Upscaled        bool     `json:"upscaled"`
	Remastered      bool     `json:"remastered"`
	Unrated         bool     `json:"unrated"`
	Uncensored      bool     `json:"uncensored"`
	Documentary     bool     `json:"documentary"`
	Commentary      bool     `json:"commentary"`
	EpisodeCode     *string  `json:"episode_code"`
	Country         *string  `json:"country"`
	Container       *string  `json:"container"`
	Extension       *string  `json:"extension"`
	Extras          []string `json:"extras"`
	Torrent         bool     `json:"torrent"`
	Scene           bool     `json:"scene"`
}

// Action response types

// ResetResponse represents reset action response
type ResetResponse struct {
	Message string   `json:"message"`
	IDs     []string `json:"ids"`
}

// RetryResponse represents retry action response
type RetryResponse struct {
	Message string   `json:"message"`
	IDs     []string `json:"ids"`
}

// RemoveResponse represents remove action response
type RemoveResponse struct {
	Message string   `json:"message"`
	IDs     []string `json:"ids"`
}

// PauseResponse represents pause/unpause action response
type PauseResponse struct {
	Message string   `json:"message"`
	IDs     []string `json:"ids"`
}

// UpdateOngoingResponse represents update ongoing items response
type UpdateOngoingResponse struct {
	Message      string                   `json:"message"`
	UpdatedItems []map[string]interface{} `json:"updated_items"`
}

// UpdateNewReleasesResponse represents update new releases response
type UpdateNewReleasesResponse struct {
	Message      string                   `json:"message"`
	UpdatedItems []map[string]interface{} `json:"updated_items"`
}

// ReindexResponse represents reindex response
type ReindexResponse struct {
	Message string `json:"message"`
}

// FfprobeResponse represents ffprobe response
type FfprobeResponse struct {
	// Structure not fully defined in API spec
	Data interface{} `json:"data,omitempty"`
}
