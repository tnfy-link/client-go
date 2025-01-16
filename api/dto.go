package api

import "time"

type CreateLink struct {
	TargetURL string `json:"targetUrl" validate:"required,http_url"` // Target link to shorten
}

type Link struct {
	ID        string `json:"id"`                     // ID
	TargetURL string `json:"targetUrl" format:"uri"` // Target link
	URL       string `json:"url" format:"uri"`       // Short link

	CreatedAt  time.Time `json:"createdAt" format:"date-time"`  // Created at
	ValidUntil time.Time `json:"validUntil" format:"date-time"` // Valid until
}

type GetLinkResponse struct {
	Link Link `json:"link"`
}

type PostLinksRequest struct {
	Link CreateLink `json:"link"`
}

type PostLinksResponse struct {
	Link Link `json:"link"`
}

type Stats struct {
	Labels map[string]map[string]int `json:"labels"` // Redirects by labels' values
	Total  int                       `json:"total"`  // Total redirects
}

type GetStatsResponse struct {
	Stats Stats `json:"stats"`
}
