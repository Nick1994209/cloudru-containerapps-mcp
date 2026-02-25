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

// ContainerAppsApplication implements the ContainerAppsService interface
type ContainerAppsApplication struct {
	authService domain.AuthService
	cfg         *config.Config
}

// NewContainerAppsApplication creates a new ContainerAppsApplication
func NewContainerAppsApplication(cfg *config.Config) domain.ContainerAppsService {
	return &ContainerAppsApplication{
		authService: NewAuthApplication(cfg),
		cfg:         cfg,
	}
}

// GetListContainerApps gets a list of ContainerApps from Cloud.ru API
func (c *ContainerAppsApplication) GetListContainerApps(projectID string) ([]domain.ContainerApp, error) {
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to ContainerApps API
	url := fmt.Sprintf("%s/v1/containers?projectId=%s", c.cfg.API.ContainersAPI, projectID)
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
func (c *ContainerAppsApplication) GetContainerApp(projectID string, containerAppName string) (*domain.ContainerApp, error) {
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to ContainerApps API
	url := fmt.Sprintf("%s/v1/containers/%s?projectId=%s", c.cfg.API.ContainersAPI, containerAppName, projectID)
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
func (c *ContainerAppsApplication) CreateContainerApp(request domain.CreateContainerAppRequest) (*domain.ContainerApp, error) {
	projectID := request.ProjectID
	containerAppName := request.ContainerAppName
	containerAppPort := request.ContainerAppPort
	containerAppImage := request.ContainerAppImage
	autoDeploymentsEnabled := request.AutoDeploymentsEnabled
	autoDeploymentsPattern := request.AutoDeploymentsPattern
	privileged := request.Privileged
	idleTimeout := request.IdleTimeout
	timeout := request.Timeout
	cpu := request.CPU
	minInstanceCount := request.MinInstanceCount
	maxInstanceCount := request.MaxInstanceCount
	description := request.Description
	publiclyAccessible := request.PubliclyAccessible
	protocol := request.Protocol
	environmentVariables := request.EnvironmentVariables
	command := request.Command
	args := request.Args

	// Set default values if not provided
	if minInstanceCount == 0 {
		minInstanceCount = 0
	}
	if maxInstanceCount == 0 {
		maxInstanceCount = 1
	}
	if description == "" {
		description = fmt.Sprintf("Container App %s created via MCP", containerAppName)
	}
	if protocol == "" {
		protocol = "http"
	}
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Map CPU to memory
	var memory string
	switch cpu {
	case "0.1":
		memory = "256Mi"
	case "0.2":
		memory = "512Mi"
	case "0.5":
		memory = "1Gi"
	case "1":
		memory = "2Gi"
	default:
		cpu = "0.1"
		memory = "256Mi" // default memory for unknown CPU values
	}

	// Parse environment variables from format <name>='<value>';<next_name>='value2'
	var envVars []map[string]interface{}
	if environmentVariables != "" {
		// Split by semicolon
		variables := strings.Split(environmentVariables, ";")
		for _, variable := range variables {
			// Split by first equals sign
			parts := strings.SplitN(variable, "=", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
					value = value[1 : len(value)-1]
				}
				envVars = append(envVars, map[string]interface{}{
					"name":  name,
					"value": value,
					"type":  "plain",
				})
			}
		}
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"name":        containerAppName,
		"projectId":   projectID,
		"description": description,
		"configuration": map[string]interface{}{
			"ingress": map[string]interface{}{
				"publiclyAccessible": publiclyAccessible,
			},
			"autoDeployments": map[string]interface{}{
				"enabled": autoDeploymentsEnabled,
				"pattern": autoDeploymentsPattern,
			},
			"privileged": privileged,
		},
		"template": map[string]interface{}{
			"timeout":     timeout,
			"idleTimeout": idleTimeout,
			"protocol":    protocol,
			"scaling": map[string]interface{}{
				"minInstanceCount": minInstanceCount,
				"maxInstanceCount": maxInstanceCount,
			},
			"containers": []map[string]interface{}{
				{
					"name":          containerAppName,
					"image":         containerAppImage,
					"containerPort": containerAppPort,
					"resources": map[string]string{
						"cpu":    cpu,
						"memory": memory,
					},
					"env": envVars,
				},
			},
		},
	}

	// Add command and args to the container if provided
	if len(command) > 0 {
		payload["template"].(map[string]interface{})["containers"].([]map[string]interface{})[0]["command"] = command
	}
	if len(args) > 0 {
		payload["template"].(map[string]interface{})["containers"].([]map[string]interface{})[0]["args"] = args
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make request to ContainerApps API
	url := fmt.Sprintf("%s/v2/containers/", c.cfg.API.ContainersAPI)
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
func (c *ContainerAppsApplication) DeleteContainerApp(projectID string, containerAppName string) error {
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Make DELETE request to ContainerApps API
	// According to the API documentation: DELETE https://containers.api.cloud.ru/v2/containers/<containerapp_name>
	url := fmt.Sprintf("%s/v2/containers/%s?projectId=%s", c.cfg.API.ContainersAPI, containerAppName, projectID)
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
func (c *ContainerAppsApplication) StartContainerApp(projectID string, containerAppName string) error {
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Make POST request to ContainerApps API to start the container app
	// According to the API documentation: POST https://containers.api.cloud.ru/v2/containers/<containerapp_name>:start
	url := fmt.Sprintf("%s/v2/containers/%s:start?projectId=%s", c.cfg.API.ContainersAPI, containerAppName, projectID)
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
func (c *ContainerAppsApplication) StopContainerApp(projectID string, containerAppName string) error {
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Make POST request to ContainerApps API to stop the container app
	// According to the API documentation: POST https://containers.api.cloud.ru/v2/containers/<containerapp_name>:stop
	url := fmt.Sprintf("%s/v2/containers/%s:stop?projectId=%s", c.cfg.API.ContainersAPI, containerAppName, projectID)
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

// GetContainerAppLogs gets logs for a specific ContainerApp from Cloud.ru API
func (c *ContainerAppsApplication) GetContainerAppLogs(projectID string, containerAppName string) (*domain.ContainerAppLogs, error) {
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to ContainerApps API for logs
	url := fmt.Sprintf("%s/v2/containers/%s/logs?projectId=%s", c.cfg.API.ContainersAPI, containerAppName, projectID)
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
	log.Printf("GetContainerAppLogs response - Status: %d, Body length: %d, Body: %s", resp.StatusCode, len(body), string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as a wrapper object containing a slice of ContainerAppLogEntry
	var response domain.ContainerAppLogs
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app logs response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &response, nil
}

// GetContainerAppSystemLogs gets system logs for a specific ContainerApp from Cloud.ru API
func (c *ContainerAppsApplication) GetContainerAppSystemLogs(projectID string, containerAppName string) (*domain.ContainerAppSystemLogs, error) {
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to ContainerApps API for system logs
	url := fmt.Sprintf("%s/v2/containers/%s/systemLogs?projectId=%s", c.cfg.API.ContainersAPI, containerAppName, projectID)
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
	log.Printf("GetContainerAppSystemLogs response - Status: %d, Body length: %d, Body: %s", resp.StatusCode, len(body), string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as a wrapper object containing a slice of ContainerAppSystemLogEntry
	var response domain.ContainerAppSystemLogs
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app system logs response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &response, nil
}
