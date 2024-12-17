package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	url string

	client *http.Client
}

func (c *Client) Shorten(ctx context.Context, targetURL string) (PostLinksResponse, error) {
	req := PostLinksRequest{
		Link: CreateLink{
			TargetURL: targetURL,
		},
	}
	res := PostLinksResponse{}

	return res, c.doRequest(ctx, http.MethodPost, "v1/links", nil, req, &res)
}

func (c *Client) GetLink(ctx context.Context, linkID string) (GetLinkResponse, error) {
	res := GetLinkResponse{}

	return res, c.doRequest(ctx, http.MethodGet, fmt.Sprintf("v1/links/%s", linkID), nil, nil, &res)
}

func (c *Client) doRequest(ctx context.Context, method, path string, headers map[string]string, payload, response any) error {
	var reqBody io.Reader = nil
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		reqBody = strings.NewReader(string(jsonBytes))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.url+path, reqBody)
	if err != nil {
		return err
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d with body %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(&response)
}

func New(url string, opts ...Option) *Client {
	config := DefaultConfig
	for _, opt := range opts {
		opt(&config)
	}

	return &Client{
		url:    url,
		client: config.client,
	}
}
