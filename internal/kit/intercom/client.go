package intercom

import (
	"github.com/AddMile/backend/internal/kit/httpclient"
)

const version = "2.12"

type HTTPClient struct {
	client *httpclient.HTTPClient
}

func NewClient(token string) (*HTTPClient, error) {
	if token == "" {
		return nil, ErrNoToken
	}

	headers := httpclient.DefaultHeaders{
		"Authorization":    "Bearer " + token,
		"Intercom-Version": version,
		"Accept":           "application/json",
		"Content-Type":     "application/json",
	}

	client, err := httpclient.New(httpclient.WithDefaultHeaders(headers))
	if err != nil {
		return nil, err
	}

	return &HTTPClient{client: client}, nil
}
