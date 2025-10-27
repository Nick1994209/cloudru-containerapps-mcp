package cloudru

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// AuthApplication implements the AuthService interface
type AuthApplication struct {
	creds domain.Credentials
	cfg   *config.Config
}

// NewAuthApplication creates a new AuthApplication
func NewAuthApplication(cfg *config.Config) domain.AuthService {
	return &AuthApplication{
		creds: domain.Credentials{
			KeyID:     cfg.KeyID,
			KeySecret: cfg.KeySecret,
		},
		cfg: cfg,
	}
}

// GetAccessToken gets an access token using KEY_ID and KEY_SECRET
func (a *AuthApplication) GetAccessToken() (string, error) {
	url := fmt.Sprintf("%s/api/v1/auth/token", a.cfg.API.IAMAPI)

	payload := strings.NewReader(fmt.Sprintf(`{"keyId": "%s","secret": "%s"}`, a.creds.KeyID, a.creds.KeySecret))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read authentication response body: %w", err)
	}

	// Log the response for debugging
	log.Printf("GetAccessToken response - Status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		return "", fmt.Errorf("authentication API returned empty response body with status %d", resp.StatusCode)
	}

	// Parse response to get token
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return result.AccessToken, nil
}
