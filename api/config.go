package api

import "net/http"

type Config struct {
	client *http.Client
}

type Option func(*Config)

func WithClient(client *http.Client) func(*Config) {
	return func(c *Config) {
		c.client = client
	}
}

func WithDefaultClient() func(*Config) {
	return func(c *Config) {
		c.client = http.DefaultClient
	}
}

var DefaultConfig = Config{
	client: http.DefaultClient,
}
