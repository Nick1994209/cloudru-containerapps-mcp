package cloudru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/utils"
)

// JobsApplication implements the JobsService interface
type JobsApplication struct {
	authService domain.AuthService
	cfg         *config.Config
}

// NewJobsApplication creates a new JobsApplication
func NewJobsApplication(cfg *config.Config) domain.JobsService {
	return &JobsApplication{
		authService: NewAuthApplication(cfg),
		cfg:         cfg,
	}
}

// GetListJobs gets a list of Jobs from Cloud.ru API
func (j *JobsApplication) GetListJobs(projectID string, pageSize string) ([]domain.Job, error) {
	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Set default pageSize to 100 if not provided
	if pageSize == "" {
		pageSize = "100"
	}

	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs?projectId=%s&pageSize=%s", j.cfg.API.ContainersAPI, projectID, pageSize)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as a wrapper object containing a slice of Job
	var response struct {
		Data []domain.Job `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse jobs response: %w body length: %d body: %s", err, len(body), string(body))
	}
	jobs := response.Data

	return jobs, nil
}

// GetJob gets a specific Job from Cloud.ru API by name
func (j *JobsApplication) GetJob(projectID string, jobName string) (*domain.Job, error) {
	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as a Job
	var job domain.Job
	if err := json.Unmarshal(body, &job); err != nil {
		return nil, fmt.Errorf("failed to parse job response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &job, nil
}

// CreateJob creates a new Job in Cloud.ru
func (j *JobsApplication) CreateJob(request domain.CreateJobRequest) (*domain.Operation, error) {
	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Map CPU to memory
	cpu, memory := j.parseCPU(request.JobCPU)
	envVars := j.parseEnvironmentVariables(request.JobEnvironmentVariables)

	// Prepare request body according to swagger spec
	requestBody := map[string]interface{}{
		"projectId":      request.ProjectID,
		"name":           request.JobName,
		"description":    request.JobDescription,
		"runImmediately": request.JobRunImmediately,
		"configuration": map[string]interface{}{
			"privileged": request.JobPrivileged,
		},
		"template": map[string]interface{}{
			"maxRetries": request.JobRetryCount,
			"timeout":    request.JobExecutionTimeout,
			"containers": []map[string]interface{}{
				{
					"name":  request.JobName,
					"image": request.JobImage,
					"resources": map[string]interface{}{
						"cpu":    cpu,
						"memory": memory,
					},
					"env": envVars,
				},
			},
		},
	}

	if containers, ok := requestBody["template"].(map[string]interface{})["containers"].([]map[string]interface{}); ok && len(containers) > 0 {
		if len(request.JobCommand) > 0 {
			containers[0]["command"] = request.JobCommand
		}
		if len(request.JobArgs) > 0 {
			containers[0]["args"] = request.JobArgs
		}
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs", j.cfg.API.ContainersAPI)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as an Operation
	var operation domain.Operation
	if err := json.Unmarshal(body, &operation); err != nil {
		return nil, fmt.Errorf("failed to parse operation response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &operation, nil
}

// DeleteJob deletes a specific Job from Cloud.ru
func (j *JobsApplication) DeleteJob(projectID string, jobName string) (*domain.Operation, error) {
	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID)
	req, err := http.NewRequest("DELETE", url, nil)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as an Operation
	var operation domain.Operation
	if err := json.Unmarshal(body, &operation); err != nil {
		return nil, fmt.Errorf("failed to parse operation response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &operation, nil
}

// ExecuteJob executes a specific Job in Cloud.ru
func (j *JobsApplication) ExecuteJob(projectID string, jobName string, params map[string]interface{}) (*domain.JobExecution, error) {
	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Prepare request body
	requestBody := map[string]interface{}{
		"projectId": projectID,
		"jobName":   jobName,
		"params":    params,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s:execute", j.cfg.API.ContainersAPI, jobName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as a JobExecution
	var jobExecution domain.JobExecution
	if err := json.Unmarshal(body, &jobExecution); err != nil {
		return nil, fmt.Errorf("failed to parse job execution response: %w body length: %d body: %s", err, len(body), string(body))
	}

	return &jobExecution, nil
}

// GetListExecutions gets a list of Job Executions from Cloud.ru API
func (j *JobsApplication) GetListExecutions(projectID string, jobName string, pageSize string) ([]domain.JobExecution, error) {
	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Set default pageSize to 100 if not provided
	if pageSize == "" {
		pageSize = "100"
	}

	// Build URL with query parameters
	url := fmt.Sprintf("%s/v2/jobs/%s/executions?projectId=%s&pageSize=%s", j.cfg.API.ContainersAPI, jobName, projectID, pageSize)

	// Make request to Jobs API
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response as a wrapper object containing a slice of JobExecution
	var response struct {
		Data []domain.JobExecution `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse job executions response: %w body length: %d body: %s", err, len(body), string(body))
	}
	executions := response.Data

	return executions, nil
}

// getJobRaw gets the raw response body from the Jobs API
func (j *JobsApplication) getJobRaw(projectID string, jobName string) ([]byte, error) {
	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body with status %d for project %s and job %s", resp.StatusCode, projectID, jobName)
	}

	return body, nil
}

// parseEnvironmentVariables parses environment variables from format <name>='<value>';<next_name>='value2'
func (j *JobsApplication) parseEnvironmentVariables(environmentVariables string) []map[string]interface{} {
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
				})
			}
		}
	}
	return envVars
}

// parseCPU maps CPU allocation to memory allocation
func (j *JobsApplication) parseCPU(cpu string) (string, string) {
	var memory string
	switch cpu {
	case "0.1":
		memory = "256Mi"
	case "0.2":
		memory = "512Mi"
	case "0.3":
		memory = "768Mi"
	case "0.5":
		memory = "1024Mi"
	case "1":
		memory = "4096Mi"
	default:
		// Default to 0.1 CPU and 256Mi memory for unknown values
		cpu = "0.1"
		memory = "256Mi"
	}

	return cpu, memory
}

// PatchJob patches a Job in Cloud.ru
func (j *JobsApplication) PatchJob(projectID string, jobName string, updateRequest domain.PatchJobRequest) (*domain.Operation, error) {
	// First, get the current job state
	rawBody, err := j.getJobRaw(projectID, jobName)
	if err != nil {
		return nil, fmt.Errorf("failed to get current job state for project %s and job %s: %w", projectID, jobName, err)
	}

	// Parse the current state into a map
	var currentJobMap map[string]interface{}
	if err := json.Unmarshal(rawBody, &currentJobMap); err != nil {
		return nil, fmt.Errorf("failed to parse current job state for project %s and job %s: %w", projectID, jobName, err)
	}

	// Use environment variables directly (already parsed)
	var envVars []map[string]interface{}
	if updateRequest.JobEnvironmentVariables != nil {
		envVars = j.parseEnvironmentVariables(*updateRequest.JobEnvironmentVariables)
	}

	// Map CPU to memory if provided
	var cpu, memory *string
	if updateRequest.JobCPU != nil && *updateRequest.JobCPU != "" {
		cpuVal, memoryVal := j.parseCPU(*updateRequest.JobCPU)
		cpu = &cpuVal
		memory = &memoryVal
	}

	// Prepare the new payload - only include fields that are being updated
	newPayload := map[string]interface{}{}

	// Add description if provided
	if updateRequest.JobDescription != nil {
		newPayload["description"] = *updateRequest.JobDescription
	}

	// Add runImmediately if provided
	if updateRequest.JobRunImmediately != nil {
		newPayload["runImmediately"] = *updateRequest.JobRunImmediately
	}

	// Build configuration section
	configSection := map[string]interface{}{}
	if updateRequest.JobPrivileged != nil {
		configSection["privileged"] = *updateRequest.JobPrivileged
	}
	if len(configSection) > 0 {
		newPayload["configuration"] = configSection
	}

	// Build template section
	templateSection := map[string]interface{}{}

	// Build template section with timeout and retry count
	if updateRequest.JobExecutionTimeout != nil {
		templateSection["timeout"] = *updateRequest.JobExecutionTimeout
	}
	if updateRequest.JobRetryCount != nil {
		templateSection["maxRetries"] = *updateRequest.JobRetryCount
	}

	// Build container section - start with existing container data
	containerSection := map[string]interface{}{}

	// Extract entire container from current job to preserve all fields
	if template, ok := currentJobMap["template"].(map[string]interface{}); ok {
		if containers, ok := template["containers"].([]interface{}); ok && len(containers) > 0 {
			if container, ok := containers[0].(map[string]interface{}); ok {
				// Copy all existing container fields
				for k, v := range container {
					containerSection[k] = v
				}
			}
		}
	}

	// Update image if explicitly provided in the update request
	if updateRequest.JobImage != nil {
		containerSection["image"] = *updateRequest.JobImage
	}

	// Update resources if provided
	if cpu != nil && memory != nil {
		containerSection["resources"] = map[string]string{
			"cpu":    *cpu,
			"memory": *memory,
		}
	}
	// Update environment variables if provided
	if len(envVars) > 0 {
		containerSection["env"] = envVars
	}
	// Update command if provided
	if len(updateRequest.JobCommand) > 0 {
		containerSection["command"] = updateRequest.JobCommand
	}
	// Update args if provided
	if len(updateRequest.JobArgs) > 0 {
		containerSection["args"] = updateRequest.JobArgs
	}
	if len(containerSection) > 0 {
		templateSection["containers"] = []map[string]interface{}{containerSection}
	}

	if len(templateSection) > 0 {
		newPayload["template"] = templateSection
	}

	// Convert payload to JSON
	jsonPayloanewPayload1, err := json.Marshal(newPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload for project %s and job %s: %w", projectID, jobName, err)
	}
	// Convert payload to JSON
	jsonPayloanewcurrentJobMap, err := json.Marshal(currentJobMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload for project %s and job %s: %w", projectID, jobName, err)
	}

	fmt.Println(string(jsonPayloanewPayload1))
	fmt.Println(string(jsonPayloanewcurrentJobMap))

	// Merge new payload with old data (deep merge)
	mergedPayload := utils.DeepMerge(newPayload, currentJobMap)

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(mergedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload for project %s and job %s: %w", projectID, jobName, err)
	}

	fmt.Println("PATCH: %s body=%s", fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID), string(jsonPayload))

	// Make PATCH request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body with status %d for project %s and job %s", resp.StatusCode, projectID, jobName)
	}

	// Create and return an Operation object by parsing the response body
	var response domain.Operation
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse job operation response for project %s and job %s: %w body length: %d body: %s", projectID, jobName, err, len(body), string(body))
	}

	return &response, nil
}
