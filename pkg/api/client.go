package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"riven-tui/pkg/config"
	"riven-tui/pkg/models"
)

// Client represents the Riven API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Riven API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(cfg.API.Endpoint, "/"),
		token:   cfg.API.Token,
		httpClient: &http.Client{
			Timeout: cfg.API.Timeout,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// parseResponse parses the HTTP response into the target struct
func (c *Client) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if target == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// Health checks the API health
func (c *Client) Health(ctx context.Context) (*models.MessageResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/health", nil)
	if err != nil {
		return nil, err
	}

	var result models.MessageResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetRoot gets the root API information
func (c *Client) GetRoot(ctx context.Context) (*models.RootResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/", nil)
	if err != nil {
		return nil, err
	}

	var result models.RootResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetStats gets system statistics
func (c *Client) GetStats(ctx context.Context) (*models.StatsResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/stats", nil)
	if err != nil {
		return nil, err
	}

	var result models.StatsResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetServices gets service status
func (c *Client) GetServices(ctx context.Context) (models.ServicesResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/services", nil)
	if err != nil {
		return nil, err
	}

	var result models.ServicesResponse
	err = c.parseResponse(resp, &result)
	return result, err
}

// GetStates gets available item states
func (c *Client) GetStates(ctx context.Context) (*models.StateResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/items/states", nil)
	if err != nil {
		return nil, err
	}

	var result models.StateResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetLogs gets system logs
func (c *Client) GetLogs(ctx context.Context) (*models.LogsResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/logs", nil)
	if err != nil {
		return nil, err
	}

	var result models.LogsResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetEvents gets system events
func (c *Client) GetEvents(ctx context.Context) (*models.EventResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/events", nil)
	if err != nil {
		return nil, err
	}

	var result models.EventResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetRDUser gets Real-Debrid user information
func (c *Client) GetRDUser(ctx context.Context) (*models.RDUser, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/rd", nil)
	if err != nil {
		return nil, err
	}

	var result models.RDUser
	err = c.parseResponse(resp, &result)
	return &result, err
}

// ItemsParams represents parameters for getting items
type ItemsParams struct {
	Limit    *int
	Page     *int
	Type     *string
	States   *string
	Sort     *models.SortOrder
	Search   *string
	Extended *bool
	IsAnime  *bool
}

// GetItems gets media items with optional filters
func (c *Client) GetItems(ctx context.Context, params *ItemsParams) (*models.ItemsResponse, error) {
	path := "/api/v1/items"

	if params != nil {
		query := url.Values{}
		if params.Limit != nil {
			query.Set("limit", strconv.Itoa(*params.Limit))
		}
		if params.Page != nil {
			query.Set("page", strconv.Itoa(*params.Page))
		}
		if params.Type != nil {
			query.Set("type", *params.Type)
		}
		if params.States != nil {
			query.Set("states", *params.States)
		}
		if params.Sort != nil {
			query.Set("sort", string(*params.Sort))
		}
		if params.Search != nil {
			query.Set("search", *params.Search)
		}
		if params.Extended != nil {
			query.Set("extended", strconv.FormatBool(*params.Extended))
		}
		if params.IsAnime != nil {
			query.Set("is_anime", strconv.FormatBool(*params.IsAnime))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.ItemsResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetItem gets a single media item by ID
func (c *Client) GetItem(ctx context.Context, id string, mediaType *models.MediaType, withStreams *bool) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/items/%s", id)

	query := url.Values{}
	if mediaType != nil {
		query.Set("media_type", string(*mediaType))
	}
	if withStreams != nil {
		query.Set("with_streams", strconv.FormatBool(*withStreams))
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// AddItems adds media items by TMDB or TVDB IDs
func (c *Client) AddItems(ctx context.Context, tmdbIds, tvdbIds *string, mediaType *models.MediaType) (*models.MessageResponse, error) {
	path := "/api/v1/items/add"

	query := url.Values{}
	if tmdbIds != nil {
		query.Set("tmdb_ids", *tmdbIds)
	}
	if tvdbIds != nil {
		query.Set("tvdb_ids", *tvdbIds)
	}
	if mediaType != nil {
		query.Set("media_type", string(*mediaType))
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.MessageResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// RemoveItems removes media items by IDs
func (c *Client) RemoveItems(ctx context.Context, ids string) (*models.RemoveResponse, error) {
	path := fmt.Sprintf("/api/v1/items/remove?ids=%s", url.QueryEscape(ids))

	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.RemoveResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// RetryItems retries media items by IDs
func (c *Client) RetryItems(ctx context.Context, ids string) (*models.RetryResponse, error) {
	path := fmt.Sprintf("/api/v1/items/retry?ids=%s", url.QueryEscape(ids))

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.RetryResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// ResetItems resets media items by IDs
func (c *Client) ResetItems(ctx context.Context, ids string) (*models.ResetResponse, error) {
	path := fmt.Sprintf("/api/v1/items/reset?ids=%s", url.QueryEscape(ids))

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.ResetResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// PauseItems pauses media items by IDs
func (c *Client) PauseItems(ctx context.Context, ids string) (*models.PauseResponse, error) {
	path := fmt.Sprintf("/api/v1/items/pause?ids=%s", url.QueryEscape(ids))

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.PauseResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// UnpauseItems unpauses media items by IDs
func (c *Client) UnpauseItems(ctx context.Context, ids string) (*models.PauseResponse, error) {
	path := fmt.Sprintf("/api/v1/items/unpause?ids=%s", url.QueryEscape(ids))

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.PauseResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetItemsByIMDBIds gets media items by IMDB IDs
func (c *Client) GetItemsByIMDBIds(ctx context.Context, imdbIds string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/items/imdb/%s", url.PathEscape(imdbIds))

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// RetryLibraryItems retries items in the library that failed to download
func (c *Client) RetryLibraryItems(ctx context.Context) (*models.RetryResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/items/retry_library", nil)
	if err != nil {
		return nil, err
	}

	var result models.RetryResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// UpdateOngoingItems updates state for ongoing and unreleased items
func (c *Client) UpdateOngoingItems(ctx context.Context) (*models.UpdateOngoingResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/items/update_ongoing", nil)
	if err != nil {
		return nil, err
	}

	var result models.UpdateOngoingResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// UpdateNewReleasesParams represents parameters for updating new releases
type UpdateNewReleasesParams struct {
	UpdateType *models.UpdateType
	Hours      *int
}

// UpdateNewReleases updates state for new releases
func (c *Client) UpdateNewReleases(ctx context.Context, params *UpdateNewReleasesParams) (*models.UpdateNewReleasesResponse, error) {
	path := "/api/v1/items/update_new_releases"

	if params != nil {
		query := url.Values{}
		if params.UpdateType != nil {
			query.Set("update_type", string(*params.UpdateType))
		}
		if params.Hours != nil {
			query.Set("hours", strconv.Itoa(*params.Hours))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.UpdateNewReleasesResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetItemStreams gets streams for a specific item
func (c *Client) GetItemStreams(ctx context.Context, itemID int) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/items/%d/streams", itemID)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// BlacklistStream blacklists a stream for an item
func (c *Client) BlacklistStream(ctx context.Context, itemID, streamID int) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/items/%d/streams/%d/blacklist", itemID, streamID)

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// UnblacklistStream removes a stream from blacklist for an item
func (c *Client) UnblacklistStream(ctx context.Context, itemID, streamID int) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/items/%d/streams/%d/unblacklist", itemID, streamID)

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// ResetItemStreams resets all streams for a media item
func (c *Client) ResetItemStreams(ctx context.Context, itemID int) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/items/%d/streams/reset", itemID)

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// ReindexParams represents parameters for reindexing items
type ReindexParams struct {
	ItemID *int
	TVDBId *string
	TMDBId *string
	IMDBId *string
}

// ReindexItem reindexes an item with Composite Indexer
func (c *Client) ReindexItem(ctx context.Context, params *ReindexParams) (*models.ReindexResponse, error) {
	path := "/api/v1/items/reindex"

	if params != nil {
		query := url.Values{}
		if params.ItemID != nil {
			query.Set("item_id", strconv.Itoa(*params.ItemID))
		}
		if params.TVDBId != nil {
			query.Set("tvdb_id", *params.TVDBId)
		}
		if params.TMDBId != nil {
			query.Set("tmdb_id", *params.TMDBId)
		}
		if params.IMDBId != nil {
			query.Set("imdb_id", *params.IMDBId)
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.ReindexResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// FfprobeMediaFiles parses a media file
func (c *Client) FfprobeMediaFiles(ctx context.Context, id int) (*models.FfprobeResponse, error) {
	path := fmt.Sprintf("/api/v1/items/ffprobe?id=%d", id)

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.FfprobeResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GenerateAPIKey generates a new API key
func (c *Client) GenerateAPIKey(ctx context.Context) (*models.MessageResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/generateapikey", nil)
	if err != nil {
		return nil, err
	}

	var result models.MessageResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetMount gets all files in the Riven VFS mount
func (c *Client) GetMount(ctx context.Context) (*models.MountResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/mount", nil)
	if err != nil {
		return nil, err
	}

	var result models.MountResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// UploadLogs uploads the latest log file to paste.c-net.org
func (c *Client) UploadLogs(ctx context.Context) (*models.UploadLogsResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/upload_logs", nil)
	if err != nil {
		return nil, err
	}

	var result models.UploadLogsResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetCalendar fetches the calendar of all items in the library
func (c *Client) GetCalendar(ctx context.Context) (*models.CalendarResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/calendar", nil)
	if err != nil {
		return nil, err
	}

	var result models.CalendarResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// InitiateTraktOAuth initiates Trakt OAuth flow
func (c *Client) InitiateTraktOAuth(ctx context.Context) (*models.TraktOAuthInitiateResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/trakt/oauth/initiate", nil)
	if err != nil {
		return nil, err
	}

	var result models.TraktOAuthInitiateResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// TraktOAuthCallback handles Trakt OAuth callback
func (c *Client) TraktOAuthCallback(ctx context.Context, code string) (*models.MessageResponse, error) {
	path := fmt.Sprintf("/api/v1/trakt/oauth/callback?code=%s", url.QueryEscape(code))

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.MessageResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// Settings API endpoints

// GetSettingsSchema gets the JSON schema for the settings
func (c *Client) GetSettingsSchema(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/settings/schema", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// LoadSettings loads settings from file
func (c *Client) LoadSettings(ctx context.Context) (*models.MessageResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/settings/load", nil)
	if err != nil {
		return nil, err
	}

	var result models.MessageResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// SaveSettings saves settings to file
func (c *Client) SaveSettings(ctx context.Context) (*models.MessageResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/settings/save", nil)
	if err != nil {
		return nil, err
	}

	var result models.MessageResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// GetAllSettings gets all settings
func (c *Client) GetAllSettings(ctx context.Context) (interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/settings/get/all", nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// GetSettings gets specific settings by paths
func (c *Client) GetSettings(ctx context.Context, paths string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/settings/get/%s", url.PathEscape(paths))

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// SetAllSettings sets all settings
func (c *Client) SetAllSettings(ctx context.Context, settings map[string]interface{}) (*models.MessageResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/settings/set/all", settings)
	if err != nil {
		return nil, err
	}

	var result models.MessageResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// SetSettingsRequest represents a settings update request
type SetSettingsRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// SetSettings sets specific settings
func (c *Client) SetSettings(ctx context.Context, settings []SetSettingsRequest) (*models.MessageResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/settings/set", settings)
	if err != nil {
		return nil, err
	}

	var result models.MessageResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// Scraping API endpoints

// ScrapeItemParams represents parameters for scraping an item
type ScrapeItemParams struct {
	ItemID    *string
	TMDBId    *string
	TVDBId    *string
	IMDBId    *string
	MediaType *models.MediaType
}

// ScrapeItem gets streams for an item by any supported ID
func (c *Client) ScrapeItem(ctx context.Context, params *ScrapeItemParams) (*models.ScrapeItemResponse, error) {
	path := "/api/v1/scrape/scrape"

	if params != nil {
		query := url.Values{}
		if params.ItemID != nil {
			query.Set("item_id", *params.ItemID)
		}
		if params.TMDBId != nil {
			query.Set("tmdb_id", *params.TMDBId)
		}
		if params.TVDBId != nil {
			query.Set("tvdb_id", *params.TVDBId)
		}
		if params.IMDBId != nil {
			query.Set("imdb_id", *params.IMDBId)
		}
		if params.MediaType != nil {
			query.Set("media_type", string(*params.MediaType))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.ScrapeItemResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// StartManualSessionParams represents parameters for starting a manual scraping session
type StartManualSessionParams struct {
	ItemID    *string
	TMDBId    *string
	TVDBId    *string
	IMDBId    *string
	MediaType *models.MediaType
	Magnet    *string
}

// StartManualSession starts a manual scraping session
func (c *Client) StartManualSession(ctx context.Context, params *StartManualSessionParams) (*models.StartSessionResponse, error) {
	path := "/api/v1/scrape/scrape/start_session"

	if params != nil {
		query := url.Values{}
		if params.ItemID != nil {
			query.Set("item_id", *params.ItemID)
		}
		if params.TMDBId != nil {
			query.Set("tmdb_id", *params.TMDBId)
		}
		if params.TVDBId != nil {
			query.Set("tvdb_id", *params.TVDBId)
		}
		if params.IMDBId != nil {
			query.Set("imdb_id", *params.IMDBId)
		}
		if params.MediaType != nil {
			query.Set("media_type", string(*params.MediaType))
		}
		if params.Magnet != nil {
			query.Set("magnet", *params.Magnet)
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.StartSessionResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// SelectFiles selects files for torrent id
func (c *Client) SelectFiles(ctx context.Context, sessionID string, container interface{}) (*models.SelectFilesResponse, error) {
	path := fmt.Sprintf("/api/v1/scrape/scrape/select_files/%s", sessionID)

	resp, err := c.doRequest(ctx, "POST", path, container)
	if err != nil {
		return nil, err
	}

	var result models.SelectFilesResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// UpdateAttributes matches container files to item
func (c *Client) UpdateAttributes(ctx context.Context, sessionID string, data interface{}) (*models.UpdateAttributesResponse, error) {
	path := fmt.Sprintf("/api/v1/scrape/scrape/update_attributes/%s", sessionID)

	resp, err := c.doRequest(ctx, "POST", path, data)
	if err != nil {
		return nil, err
	}

	var result models.UpdateAttributesResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// AbortManualSession aborts a manual scraping session
func (c *Client) AbortManualSession(ctx context.Context, sessionID string) (*models.SessionResponse, error) {
	path := fmt.Sprintf("/api/v1/scrape/scrape/abort_session/%s", sessionID)

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.SessionResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// CompleteManualSession completes a manual scraping session
func (c *Client) CompleteManualSession(ctx context.Context, sessionID string) (*models.SessionResponse, error) {
	path := fmt.Sprintf("/api/v1/scrape/scrape/complete_session/%s", sessionID)

	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.SessionResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// ParseTorrentTitles parses an array of torrent titles
func (c *Client) ParseTorrentTitles(ctx context.Context, titles []string) (*models.ParseTorrentTitleResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/scrape/parse", titles)
	if err != nil {
		return nil, err
	}

	var result models.ParseTorrentTitleResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// Streaming API endpoints

// GetEventTypes gets available event types for streaming
func (c *Client) GetEventTypes(ctx context.Context) (interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/stream/event_types", nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}

// StreamEvents streams events of a specific type
func (c *Client) StreamEvents(ctx context.Context, eventType string) (*models.StreamEventResponse, error) {
	path := fmt.Sprintf("/api/v1/stream/%s", eventType)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result models.StreamEventResponse
	err = c.parseResponse(resp, &result)
	return &result, err
}

// StreamEventsSSE streams real-time events using Server-Sent Events
func (c *Client) StreamEventsSSE(ctx context.Context, eventChan chan<- models.Event) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v1/events/stream", nil)
	if err != nil {
		return err
	}

	// Add authentication header
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("event stream failed with status: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse SSE format
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			var event models.Event
			if err := json.Unmarshal([]byte(data), &event); err == nil {
				select {
				case eventChan <- event:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}

	return scanner.Err()
}

// Webhook endpoints

// OverseerrWebhook handles Overseerr webhook
func (c *Client) OverseerrWebhook(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/webhook/overseerr", data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = c.parseResponse(resp, &result)
	return result, err
}
