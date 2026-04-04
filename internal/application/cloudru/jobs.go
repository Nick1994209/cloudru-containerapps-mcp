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

// makeHTTPRequest makes an HTTP request to the Cloud.ru API
func (j *JobsApplication) makeHTTPRequest(method, url string, body []byte) ([]byte, error) {
	token, err := j.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// GetListJobs gets a list of Jobs from Cloud.ru API
func (j *JobsApplication) GetListJobs(projectID string, pageSize string) ([]domain.Job, error) {
	// Set default pageSize to 100 if not provided
	if pageSize == "" {
		pageSize = "100"
	}

	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs?projectId=%s&pageSize=%s", j.cfg.API.ContainersAPI, projectID, pageSize)
	body, err := j.makeHTTPRequest("GET", url, nil)
	if err != nil {
		return nil, err
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
	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID)
	body, err := j.makeHTTPRequest("GET", url, nil)
	if err != nil {
		return nil, err
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
	// Map CPU to memory
	cpu, memory := utils.ParseCPU(request.JobCPU)
	envVars := utils.ParseEnvironmentVariables(request.JobEnvironmentVariables)

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
	body, err := j.makeHTTPRequest("POST", url, jsonBody)
	if err != nil {
		return nil, err
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
	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID)
	body, err := j.makeHTTPRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
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
	body, err := j.makeHTTPRequest("POST", url, jsonBody)
	if err != nil {
		return nil, err
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
	// Set default pageSize to 100 if not provided
	if pageSize == "" {
		pageSize = "100"
	}

	// Build URL with query parameters
	url := fmt.Sprintf("%s/v2/jobs/%s/executions?projectId=%s&pageSize=%s", j.cfg.API.ContainersAPI, jobName, projectID, pageSize)

	// Make request to Jobs API
	body, err := j.makeHTTPRequest("GET", url, nil)
	if err != nil {
		return nil, err
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
	// Make request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID)
	body, err := j.makeHTTPRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body for project %s and job %s", projectID, jobName)
	}

	return body, nil
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

	// Update description if provided
	if updateRequest.JobDescription != nil {
		currentJobMap["description"] = *updateRequest.JobDescription
	}

	// Update runImmediately if provided
	if updateRequest.JobRunImmediately != nil {
		currentJobMap["runImmediately"] = *updateRequest.JobRunImmediately
	}

	// Update configuration section
	if updateRequest.JobPrivileged != nil {
		if config, ok := currentJobMap["configuration"].(map[string]interface{}); ok {
			config["privileged"] = *updateRequest.JobPrivileged
		} else {
			currentJobMap["configuration"] = map[string]interface{}{
				"privileged": *updateRequest.JobPrivileged,
			}
		}
	}

	// Use environment variables directly (already parsed)
	var envVars []map[string]interface{}
	if updateRequest.JobEnvironmentVariables != nil {
		envVars = utils.ParseEnvironmentVariables(*updateRequest.JobEnvironmentVariables)
	}

	// Map CPU to memory if provided
	var cpu, memory *string
	if updateRequest.JobCPU != nil && *updateRequest.JobCPU != "" {
		cpuVal, memoryVal := utils.ParseCPU(*updateRequest.JobCPU)
		cpu = &cpuVal
		memory = &memoryVal
	}

	// Update template section
	if template, ok := currentJobMap["template"].(map[string]interface{}); ok {
		// Update timeout if provided
		if updateRequest.JobExecutionTimeout != nil {
			template["timeout"] = *updateRequest.JobExecutionTimeout
		}

		// Update maxRetries if provided
		if updateRequest.JobRetryCount != nil {
			template["maxRetries"] = *updateRequest.JobRetryCount
		}

		// Update container section
		if containers, ok := template["containers"].([]interface{}); ok && len(containers) > 0 {
			if container, ok := containers[0].(map[string]interface{}); ok {
				// Update image if provided
				if updateRequest.JobImage != nil {
					container["image"] = *updateRequest.JobImage
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
				if len(updateRequest.JobCommand) > 0 {
					container["command"] = updateRequest.JobCommand
				}

				// Update args if provided
				if len(updateRequest.JobArgs) > 0 {
					container["args"] = updateRequest.JobArgs
				}
			}
		}
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(currentJobMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload for project %s and job %s: %w", projectID, jobName, err)
	}

	// Make PATCH request to Jobs API
	url := fmt.Sprintf("%s/v2/jobs/%s?projectId=%s", j.cfg.API.ContainersAPI, jobName, projectID)
	body, err := j.makeHTTPRequest("PATCH", url, jsonPayload)
	if err != nil {
		return nil, err
	}

	// Check if body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("API returned empty response body for project %s and job %s", projectID, jobName)
	}

	// Create and return an Operation object by parsing the response body
	var response domain.Operation
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse job operation response for project %s and job %s: %w body length: %d body: %s", projectID, jobName, err, len(body), string(body))
	}

	return &response, nil
}
