package application

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// ContainerAppsApplication implements the ContainerAppsService and DockerRegistryService interfaces
type ContainerAppsApplication struct{}

// NewContainerAppsApplication creates a new ContainerAppsApplication
func NewContainerAppsApplication() domain.ContainerAppsService {
	return &ContainerAppsApplication{}
}

// GetListContainerApps gets a list of ContainerApps from Cloud.ru API
func (c *ContainerAppsApplication) GetListContainerApps(projectID string, credentials domain.Credentials) ([]domain.ContainerApp, error) {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to ContainerApps API
	url := fmt.Sprintf("https://containers.api.cloud.ru/v1/containers?projectId=%s", projectID)
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
	log.Printf("GetListContainerApps response - Status: %d, Body length: %d, Body: %s", resp.StatusCode, len(body), string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as a wrapper object containing a slice of ContainerApp
	var response struct {
		Data []domain.ContainerApp `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse containerapps response: %w body length: %d body: %s", err, len(body), string(body))
	}
	containerApps := response.Data

	return containerApps, nil
}

// GetContainerApp gets a specific ContainerApp from Cloud.ru API
func (c *ContainerAppsApplication) GetContainerApp(projectID string, containerAppName string, credentials domain.Credentials) (*domain.ContainerApp, error) {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to ContainerApps API
	url := fmt.Sprintf("https://containers.api.cloud.ru/v1/containers/%s?projectId=%s", containerAppName, projectID)
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
	log.Printf("GetContainerApp response - Status: %d, Body length: %d, Body: %s", resp.StatusCode, len(body), string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body with status %d", resp.StatusCode)
	}

	// Parse response
	var containerApp domain.ContainerApp
	if err := json.Unmarshal(body, &containerApp); err != nil {
		return nil, fmt.Errorf("failed to parse containerapp response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &containerApp, nil
}

// CreateContainerApp creates a new ContainerApp in Cloud.ru
func (c *ContainerAppsApplication) CreateContainerApp(projectID string, containerAppName string, containerAppPort int, containerAppImage string, credentials domain.Credentials) (*domain.ContainerApp, error) {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"name":        containerAppName,
		"projectId":   projectID,
		"description": fmt.Sprintf("Container App %s created via MCP", containerAppName),
		"configuration": map[string]interface{}{
			"ingress": map[string]interface{}{
				"publiclyAccessible": true,
			},
			"autoDeployments": map[string]interface{}{
				"enabled": false,
				"pattern": "",
			},
			"privileged": false,
		},
		"template": map[string]interface{}{
			"timeout":     "60s",
			"idleTimeout": "600s",
			"containers": []map[string]interface{}{
				{
					"name":          containerAppName,
					"image":         containerAppImage,
					"containerPort": containerAppPort,
					"env": []map[string]string{
						{
							"name":  "CONTAINERAPP_NAME",
							"value": containerAppName,
						},
					},
					"resources": []map[string]string{
						{
							"cpu":    "0.1",
							"memory": "256Mi",
						},
					},
				},
			},
		},
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make request to ContainerApps API
	url := "https://containers.api.cloud.ru/v2/containers/"
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
	log.Printf("CreateContainerApp response - Status: %d, Body length: %d, Body: %s", resp.StatusCode, len(body), string(body))

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body with status %d", resp.StatusCode)
	}

	// Parse response
	var containerApp domain.ContainerApp
	if err := json.Unmarshal(body, &containerApp); err != nil {
		return nil, fmt.Errorf("failed to parse containerapp response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &containerApp, nil
}

// DeleteContainerApp deletes a ContainerApp from Cloud.ru
func (c *ContainerAppsApplication) DeleteContainerApp(projectID string, containerAppName string, credentials domain.Credentials) error {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Make DELETE request to ContainerApps API
	// According to the API documentation: DELETE https://containers.api.cloud.ru/v2/containers/<containerapp_name>
	url := fmt.Sprintf("https://containers.api.cloud.ru/v2/containers/%s?projectId=%s", containerAppName, projectID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// According to the API documentation, a successful deletion should return 204 No Content
	// but we'll accept 200 OK as well
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// StartContainerApp starts a ContainerApp in Cloud.ru
func (c *ContainerAppsApplication) StartContainerApp(projectID string, containerAppName string, credentials domain.Credentials) error {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Make POST request to ContainerApps API to start the container app
	// According to the API documentation: POST https://containers.api.cloud.ru/v2/containers/<containerapp_name>:start
	url := fmt.Sprintf("https://containers.api.cloud.ru/v2/containers/%s:start?projectId=%s", containerAppName, projectID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// According to the API documentation, a successful start should return 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// StopContainerApp stops a ContainerApp in Cloud.ru
func (c *ContainerAppsApplication) StopContainerApp(projectID string, containerAppName string, credentials domain.Credentials) error {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Make POST request to ContainerApps API to stop the container app
	// According to the API documentation: POST https://containers.api.cloud.ru/v2/containers/<containerapp_name>:stop
	url := fmt.Sprintf("https://containers.api.cloud.ru/v2/containers/%s:stop?projectId=%s", containerAppName, projectID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// According to the API documentation, a successful stop should return 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// getAccessToken gets an access token using KEY_ID and KEY_SECRET
func (c *ContainerAppsApplication) getAccessToken(keyID, keySecret string) (string, error) {
	url := "https://iam.api.cloud.ru/api/v1/auth/token"

	payload := strings.NewReader(fmt.Sprintf(`{"keyId": "%s","secret": "%s"}`, keyID, keySecret))

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
		return "", fmt.Errorf("failed to read containerapps response body: %w", err)
	}

	// Log the response for debugging
	log.Printf("getAccessToken response - Status: %d", resp.StatusCode)

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

// GetListDockerRegistries gets a list of Docker Registries from Cloud.ru API
func (c *ContainerAppsApplication) GetListDockerRegistries(projectID string, credentials domain.Credentials) ([]domain.DockerRegistry, error) {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to Docker Registries API
	url := fmt.Sprintf("https://ar.api.cloud.ru/v1/projects/%s/registries", projectID)
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
func (c *ContainerAppsApplication) CreateDockerRegistry(projectID string, registryName string, isPublic bool, credentials domain.Credentials) (*domain.DockerRegistry, error) {
	// Get access token using KEY_ID and KEY_SECRET
	token, err := c.getAccessToken(credentials.KeyID, credentials.KeySecret)
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
	url := fmt.Sprintf("https://ar.api.cloud.ru/v1/projects/%s/registries", projectID)
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
