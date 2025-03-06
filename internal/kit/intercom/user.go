package intercom

import (
	"context"
	"fmt"
)

type UpdateUserRequest struct {
	ExternalID string `json:"external_id,omitempty"`
	Email      string `json:"email,omitempty"`
	Platform   string `json:"platform,omitempty"`
	Provider   string `json:"provider,omitempty"`
	ChatTier   string `json:"chat_tier,omitempty"`
}

func (c *HTTPClient) UpdateUser(ctx context.Context, userID string, request UpdateUserRequest) error {
	url := fmt.Sprintf("https://api.intercom.io/contacts/%s", userID)

	return c.client.PUT(ctx, url, request, nil)
}

type UserIDByEmailResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
	TotalCount int `json:"total_count"`
}

func (c *HTTPClient) UserIDByEmail(ctx context.Context, email string) (string, error) {
	url := "https://api.intercom.io/contacts/search"

	requestBody := map[string]any{
		"query": map[string]any{
			"operator": "AND",
			"value": []map[string]any{
				{
					"field":    "email",
					"operator": "=",
					"value":    email,
				},
			},
		},
	}

	var responseBody UserIDByEmailResponse
	err := c.client.POST(ctx, url, requestBody, &responseBody)
	if err != nil {
		return "", fmt.Errorf("cannot find user by email: %w", err)
	}

	if responseBody.TotalCount == 0 {
		return "", ErrUserNotFound
	}

	if responseBody.TotalCount > 1 {
		return "", ErrTooManyUsers
	}

	return responseBody.Data[0].ID, nil
}
