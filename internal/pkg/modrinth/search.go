package modrinth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/mworzala/mc/internal/pkg/modrinth/facet"
)

type SearchIndex string

const (
	noIndex   SearchIndex = ""
	Relevance SearchIndex = "relevance"
	Downloads SearchIndex = "downloads"
	Follows   SearchIndex = "follows"
	Newest    SearchIndex = "newest"
	Updated   SearchIndex = "updated"
)

type SearchRequest struct {
	Query  string
	Facets facet.Root
	Index  SearchIndex // Default Relevance
	Offset int         // Default 0
	Limit  int         // Default 10, max 100
}

type SearchResponse struct {
	Hits      []SearchResult `json:"hits"`
	Offset    int            `json:"offset"`
	Limit     int            `json:"limit"`
	TotalHits int            `json:"total_hits"`
}

type SearchResult struct {
	ProjectID    string        `json:"project_id"`
	ProjectType  ProjectType   `json:"project_type"`
	Slug         string        `json:"slug"`
	DateCreated  time.Time     `json:"date_created"`
	DateModified time.Time     `json:"date_modified"`
	Title        string        `json:"title"`
	Author       string        `json:"author"`
	Categories   []string      `json:"categories"`
	ClientSide   SupportStatus `json:"client_side"`
	ServerSide   SupportStatus `json:"server_side"`

	Downloads         int      `json:"downloads"`
	Follows           int      `json:"follows"`
	IconURL           *string  `json:"icon_url"`
	Color             *int     `json:"color"` // The RGB color of the project, automatically generated from the project icon
	Description       string   `json:"description"`
	DisplayCategories []string `json:"display_categories"`
	Gallery           []string `json:"gallery"`          // All gallery images attached to the project
	FeaturedGallery   *string  `json:"featured_gallery"` // The featured gallery image of the project
	License           string   `json:"license"`          // The SPDX license ID of a project

	ThreadID           *string             `json:"thread_id"` // The ID of the moderation thread associated with this project
	MonetizationStatus *MonetizationStatus `json:"monetization_status"`

	Versions      []string `json:"versions"`       // A list of the minecraft versions supported by the project
	LatestVersion string   `json:"latest_version"` // The latest version of minecraft that this project supports
}

func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	if req.Index == noIndex {
		req.Index = Relevance
	}
	if !searchIndexValidationMap[req.Index] {
		return nil, fmt.Errorf("unknown search index: %s", req.Index)
	}
	if req.Offset < 0 {
		return nil, errors.New("offset must be positive")
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit < 0 {
		return nil, errors.New("limit must be positive")
	}
	if req.Limit > 100 {
		return nil, errors.New("limit must be less than or equal to 100")
	}
	facets, err := facet.ToString(req.Facets)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	if req.Query != "" {
		params.Add("query", req.Query)
	}
	if facets != "" {
		params.Add("facets", facets)
	}
	params.Add("index", string(req.Index))
	params.Add("offset", strconv.Itoa(req.Offset))
	params.Add("limit", strconv.Itoa(req.Limit))

	return get[SearchResponse](c, ctx, "/search", params)
}

var searchIndexValidationMap = map[SearchIndex]bool{
	Relevance: true,
	Downloads: true,
	Follows:   true,
	Newest:    true,
	Updated:   true,
}
