package cloudru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/utils"
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

// makeHTTPRequest makes an HTTP request to the Cloud.ru API
func (c *ContainerAppsApplication) makeHTTPRequest(method, path string, body []byte) ([]byte, error) {
	token, err := c.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

	url := c.cfg.API.ContainersAPI + path
	req, err := http.NewRequest(method, url, bodyReader)
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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if status code is 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// GetListContainerApps gets a list of ContainerApps from Cloud.ru API
func (c *ContainerAppsApplication) GetListContainerApps(projectID string) ([]domain.ContainerApp, error) {
	// Make request to ContainerApps API
	path := fmt.Sprintf("/v1/containers?projectId=%s", projectID)
	body, err := c.makeHTTPRequest("GET", path, nil)
	if err != nil {
		return nil, err
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

// getContainerAppRaw gets the raw response body from the ContainerApps API
func (c *ContainerAppsApplication) getContainerAppRaw(projectID string, containerAppName string) ([]byte, error) {
	// Make request to ContainerApps API
	path := fmt.Sprintf("/v1/containers/%s?projectId=%s", containerAppName, projectID)
	body, err := c.makeHTTPRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body")
	}

	return body, nil
}

// GetContainerApp gets a specific ContainerApp from Cloud.ru API
func (c *ContainerAppsApplication) GetContainerApp(projectID string, containerAppName string) (*domain.ContainerApp, error) {
	// Get the raw response body
	rawBody, err := c.getContainerAppRaw(projectID, containerAppName)
	if err != nil {
		return nil, err
	}

	// Parse response
	var containerApp domain.ContainerApp
	if err := json.Unmarshal(rawBody, &containerApp); err != nil {
		return nil, fmt.Errorf("failed to parse containerapp response for '%s': %w body length: %d body: %s", containerAppName, err, len(rawBody), string(rawBody))
	}

	return &containerApp, nil
}

// CreateContainerApp creates a new ContainerApp in Cloud.ru
func (c *ContainerAppsApplication) CreateContainerApp(request domain.CreateContainerAppRequest) (*domain.Operation, error) {
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
		protocol = "http_1"
	}

	// Parse environment variables
	envVars := utils.ParseEnvironmentVariables(environmentVariables)

	// Map CPU to memory
	cpu, memory := utils.ParseCPU(cpu)

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
	path := "/v2/containers/"
	body, err := c.makeHTTPRequest("POST", path, jsonPayload)
	if err != nil {
		return nil, err
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body")
	}

	// Create and return an Operation object by parsing the response body
	var response domain.Operation
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app operation response for '%s': %w body length: %d body: %s", containerAppName, err, len(body), string(body))
	}

	return &response, nil
}

// DeleteContainerApp deletes a ContainerApp from Cloud.ru
func (c *ContainerAppsApplication) DeleteContainerApp(projectID string, containerAppName string) (*domain.Operation, error) {
	// Make DELETE request to ContainerApps API
	// According to the API documentation: DELETE https://containers.api.cloud.ru/v2/containers/<containerapp_name>
	path := fmt.Sprintf("/v2/containers/%s?projectId=%s", containerAppName, projectID)
	body, err := c.makeHTTPRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	// Create and return an Operation object by parsing the response body
	var response domain.Operation
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app operation response for '%s': %w body length: %d body: %s", containerAppName, err, len(body), string(body))
	}

	return &response, nil
}

// StartContainerApp starts a ContainerApp in Cloud.ru
func (c *ContainerAppsApplication) StartContainerApp(projectID string, containerAppName string) (*domain.Operation, error) {
	// Make POST request to ContainerApps API to start the container app
	// According to the API documentation: POST https://containers.api.cloud.ru/v2/containers/<containerapp_name>:start
	path := fmt.Sprintf("/v2/containers/%s:start?projectId=%s", containerAppName, projectID)
	body, err := c.makeHTTPRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	// Create and return an Operation object by parsing the response body
	var response domain.Operation
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app operation response for '%s': %w body length: %d body: %s", containerAppName, err, len(body), string(body))
	}

	return &response, nil
}

// StopContainerApp stops a ContainerApp in Cloud.ru
func (c *ContainerAppsApplication) StopContainerApp(projectID string, containerAppName string) (*domain.Operation, error) {
	// Make POST request to ContainerApps API to stop the container app
	// According to the API documentation: POST https://containers.api.cloud.ru/v2/containers/<containerapp_name>:stop
	path := fmt.Sprintf("/v2/containers/%s:stop?projectId=%s", containerAppName, projectID)
	body, err := c.makeHTTPRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	var response domain.Operation
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app operation response for '%s': %w body length: %d body: %s", containerAppName, err, len(body), string(body))
	}

	return &response, nil
}

// GetContainerAppLogs gets logs for a specific ContainerApp from Cloud.ru API
func (c *ContainerAppsApplication) GetContainerAppLogs(projectID string, containerAppName string) (*domain.ContainerAppLogs, error) {
	// Make request to ContainerApps API for logs
	path := fmt.Sprintf("/v2/containers/%s/logs?projectId=%s", containerAppName, projectID)
	body, err := c.makeHTTPRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response as a wrapper object containing a slice of ContainerAppLogEntry
	var response domain.ContainerAppLogs
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app logs response for '%s': %w body length: %d body: %s", containerAppName, err, len(body), string(body))
	}

	return &response, nil
}

// GetContainerAppSystemLogs gets system logs for a specific ContainerApp from Cloud.ru API
func (c *ContainerAppsApplication) GetContainerAppSystemLogs(projectID string, containerAppName string) (*domain.ContainerAppSystemLogs, error) {
	// Make request to ContainerApps API for system logs
	path := fmt.Sprintf("/v2/containers/%s/systemLogs?projectId=%s", containerAppName, projectID)
	body, err := c.makeHTTPRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response as a wrapper object containing a slice of ContainerAppSystemLogEntry
	var response domain.ContainerAppSystemLogs
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app system logs response for '%s': %w body length: %d body: %s", containerAppName, err, len(body), string(body))
	}

	return &response, nil
}

// PatchContainerApp patches a ContainerApp in Cloud.ru
func (c *ContainerAppsApplication) PatchContainerApp(projectID string, containerAppName string, updateRequest domain.PatchContainerAppRequest) (*domain.Operation, error) {
	// First, get the current container app state
	rawBody, err := c.getContainerAppRaw(projectID, containerAppName)
	if err != nil {
		return nil, fmt.Errorf("failed to get current container app state for '%s': %w", containerAppName, err)
	}

	// Parse the current state
	var currentContainerApp map[string]interface{}
	if err := json.Unmarshal(rawBody, &currentContainerApp); err != nil {
		return nil, fmt.Errorf("failed to parse current container app state for '%s': %w", containerAppName, err)
	}

	// Use environment variables directly (already parsed)
	var envVars []map[string]interface{}
	if updateRequest.EnvironmentVariables != nil {
		envVars = utils.ParseEnvironmentVariables(*updateRequest.EnvironmentVariables)
	}

	// Map CPU to memory if provided
	var cpu, memory *string
	if updateRequest.CPU != nil && *updateRequest.CPU != "" {
		cpuVal, memoryVal := utils.ParseCPU(*updateRequest.CPU)
		cpu = &cpuVal
		memory = &memoryVal
	}

	// Update description if provided
	if updateRequest.Description != nil {
		currentContainerApp["description"] = *updateRequest.Description
	}

	// Update configuration section
	if updateRequest.PubliclyAccessible != nil {
		if config, ok := currentContainerApp["configuration"].(map[string]interface{}); ok {
			if ingress, ok := config["ingress"].(map[string]interface{}); ok {
				ingress["publiclyAccessible"] = *updateRequest.PubliclyAccessible
			} else {
				config["ingress"] = map[string]interface{}{
					"publiclyAccessible": *updateRequest.PubliclyAccessible,
				}
			}
		} else {
			currentContainerApp["configuration"] = map[string]interface{}{
				"ingress": map[string]interface{}{
					"publiclyAccessible": *updateRequest.PubliclyAccessible,
				},
			}
		}
	}

	if updateRequest.AutoDeploymentsEnabled != nil || updateRequest.AutoDeploymentsPattern != nil {
		if config, ok := currentContainerApp["configuration"].(map[string]interface{}); ok {
			if autoDeployments, ok := config["autoDeployments"].(map[string]interface{}); ok {
				if updateRequest.AutoDeploymentsEnabled != nil {
					autoDeployments["enabled"] = *updateRequest.AutoDeploymentsEnabled
				}
				if updateRequest.AutoDeploymentsPattern != nil {
					autoDeployments["pattern"] = *updateRequest.AutoDeploymentsPattern
				}
			} else {
				autoDeployments := map[string]interface{}{}
				if updateRequest.AutoDeploymentsEnabled != nil {
					autoDeployments["enabled"] = *updateRequest.AutoDeploymentsEnabled
				}
				if updateRequest.AutoDeploymentsPattern != nil {
					autoDeployments["pattern"] = *updateRequest.AutoDeploymentsPattern
				}
				config["autoDeployments"] = autoDeployments
			}
		} // else not required, configuration should exists in body
	}

	// Update template section
	if template, ok := currentContainerApp["template"].(map[string]interface{}); ok {
		// Update timeout if provided
		if updateRequest.Timeout != nil {
			template["timeout"] = *updateRequest.Timeout
		}

		// Update idleTimeout if provided
		if updateRequest.IdleTimeout != nil {
			template["idleTimeout"] = *updateRequest.IdleTimeout
		}

		// Update protocol if provided
		if updateRequest.Protocol != nil {
			template["protocol"] = *updateRequest.Protocol
		}

		// Update scaling section
		if updateRequest.MinInstanceCount != nil || updateRequest.MaxInstanceCount != nil {
			if scaling, ok := template["scaling"].(map[string]interface{}); ok {
				if updateRequest.MinInstanceCount != nil {
					scaling["minInstanceCount"] = *updateRequest.MinInstanceCount
				}
				if updateRequest.MaxInstanceCount != nil {
					scaling["maxInstanceCount"] = *updateRequest.MaxInstanceCount
				}
			} // else not required, scaling should exists in template
		}

		// Update container section
		if containers, ok := template["containers"].([]interface{}); ok && len(containers) > 0 {
			if container, ok := containers[0].(map[string]interface{}); ok {
				// Update image if provided
				if updateRequest.ContainerAppImage != nil {
					container["image"] = *updateRequest.ContainerAppImage
				}

				// Update containerPort if provided
				if updateRequest.ContainerAppPort != nil {
					container["containerPort"] = *updateRequest.ContainerAppPort
				}

				// Update resources if provided
				if cpu != nil && memory != nil {
					container["resources"] = map[string]string{
						"cpu":    *cpu,
						"memory": *memory,
					}
				}

				// Update environment variables if provided
				if len(envVars) > 0 {
					container["env"] = envVars
				}

				// Update command if provided
				if len(updateRequest.Command) > 0 {
					container["command"] = updateRequest.Command
				}

				// Update args if provided
				if len(updateRequest.Args) > 0 {
					container["args"] = updateRequest.Args
				}
			}
		}
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(currentContainerApp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make PATCH request to ContainerApps API
	path := fmt.Sprintf("/v2/containers/%s?projectId=%s", containerAppName, projectID)
	body, err := c.makeHTTPRequest("PATCH", path, jsonPayload)
	if err != nil {
		return nil, err
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body")
	}

	// Create and return an Operation object by parsing the response body
	var response domain.Operation
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse container app operation response for '%s': %w body length: %d body: %s", containerAppName, err, len(body), string(body))
	}

	return &response, nil
}
