package models

import "time"

// ScrapeItemResponse represents the response from scraping an item
type ScrapeItemResponse struct {
	Message string            `json:"message"`
	Streams map[string]Stream `json:"streams"`
}

// StartSessionResponse represents the response from starting a manual scraping session
type StartSessionResponse struct {
	Message     string            `json:"message"`
	SessionID   string            `json:"session_id"`
	TorrentID   string            `json:"torrent_id"`
	TorrentInfo TorrentInfo       `json:"torrent_info"`
	Containers  *TorrentContainer `json:"containers"`
	ExpiresAt   string            `json:"expires_at"`
}

// SelectFilesResponse represents the response from selecting files
type SelectFilesResponse struct {
	Message      string       `json:"message"`
	DownloadType DownloadType `json:"download_type"`
}

// UpdateAttributesResponse represents the response from updating attributes
type UpdateAttributesResponse struct {
	Message string `json:"message"`
}

// SessionResponse represents a generic session response
type SessionResponse struct {
	Message string `json:"message"`
}

// ParseTorrentTitleResponse represents the response from parsing torrent titles
type ParseTorrentTitleResponse struct {
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

// TorrentInfo represents torrent information from a debrid service
type TorrentInfo struct {
	ID                  interface{}            `json:"id"` // Can be int or string
	Name                string                 `json:"name"`
	Status              *string                `json:"status"`
	InfoHash            *string                `json:"infohash"`
	Progress            *float64               `json:"progress"`
	Bytes               *int64                 `json:"bytes"`
	CreatedAt           *time.Time             `json:"created_at"`
	ExpiresAt           *time.Time             `json:"expires_at"`
	CompletedAt         *time.Time             `json:"completed_at"`
	AlternativeFilename *string                `json:"alternative_filename"`
	Files               map[string]interface{} `json:"files"`
	Links               []string               `json:"links"`
}

// TorrentContainer represents a collection of files from an infohash from a debrid service
type TorrentContainer struct {
	InfoHash string       `json:"infohash"`
	Files    []DebridFile `json:"files,omitempty"`
}

// DebridFile represents a file from a debrid service
type DebridFile struct {
	Filename string `json:"filename"`
	Filesize int64  `json:"filesize,omitempty"`
	Link     string `json:"link,omitempty"`
	ID       string `json:"id,omitempty"`
}

// ShowFileData represents root model for show file data that maps seasons to episodes to file data
// Example:
//
//	{
//	    1: {  # Season 1
//	        1: {"filename": "path/to/s01e01.mkv"},  # Episode 1
//	        2: {"filename": "path/to/s01e02.mkv"}   # Episode 2
//	    },
//	    2: {  # Season 2
//	        1: {"filename": "path/to/s02e01.mkv"}   # Episode 1
//	    }
//	}
type ShowFileData map[int]map[int]DebridFile

// Container represents a container for manual scraping
type Container struct {
	InfoHash string       `json:"infohash"`
	Files    []DebridFile `json:"files,omitempty"`
}
