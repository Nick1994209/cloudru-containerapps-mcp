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

// ArtifactRegistryApplication implements the ArtifactRegistryService interface
type ArtifactRegistryApplication struct {
	authService domain.AuthService
	cfg         *config.Config
}

// NewArtifactRegistryApplication creates a new ArtifactRegistryApplication
func NewArtifactRegistryApplication(cfg *config.Config) domain.ArtifactRegistryService {
	return &ArtifactRegistryApplication{
		authService: NewAuthApplication(cfg),
		cfg:         cfg,
	}
}

// GetListDockerRegistries gets a list of Docker Registries from Cloud.ru API
func (d *ArtifactRegistryApplication) GetListDockerRegistries(projectID string) ([]domain.DockerRegistry, error) {
	token, err := d.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to Docker Registries API
	url := fmt.Sprintf("%s/v1/projects/%s/registries", d.cfg.API.ArtifactAPI, projectID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the response for debugging
	log.Printf("GetListDockerRegistries response - Status: %d, Body length: %d, Body: %s", resp.StatusCode, len(body), string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		// Return empty slice if no registries found
		return []domain.DockerRegistry{}, nil
	}

	// Parse response
	var response struct {
		Registries []domain.DockerRegistry `json:"registries"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse registries response: %w body length: %d body: %s", err, len(body), string(body))
	}

	// Filter only DOCKER registries
	dockerRegistries := []domain.DockerRegistry{}
	for _, registry := range response.Registries {
		if registry.RegistryType == "DOCKER" {
			dockerRegistries = append(dockerRegistries, registry)
		}
	}

	return dockerRegistries, nil
}

// CreateDockerRegistry creates a new Docker Registry in Cloud.ru
func (d *ArtifactRegistryApplication) CreateDockerRegistry(projectID string, registryName string, isPublic bool) (*domain.DockerRegistry, error) {
	token, err := d.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"name":         registryName,
		"isPublic":     isPublic,
		"registryType": "DOCKER",
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make request to Docker Registries API
	url := fmt.Sprintf("%s/v1/projects/%s/registries", d.cfg.API.ArtifactAPI, projectID)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the response for debugging
	log.Printf("CreateDockerRegistry response - Status: %d, Body length: %d, Body: %s", resp.StatusCode, len(body), string(body))

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body with status %d", resp.StatusCode)
	}

	// Parse response
	var registry domain.DockerRegistry
	if err := json.Unmarshal(body, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &registry, nil
}
