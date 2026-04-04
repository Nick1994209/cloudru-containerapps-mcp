package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterCreateJobTool registers the create job tool with the MCP server
func (s *MCPServer) RegisterCreateJobTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Create a new Job in Cloud.ru",
		"project_id",
		"job_name",
		"job_image",
		"job_privileged",
		"job_cpu",
		"job_description",
		"job_environment_variables",
		"job_command",
		"job_args",
		"job_retry_count",
		"job_execution_timeout",
		"job_run_immediately",
	)
	createJobTool := mcp.NewTool("cloudru_create_job", toolOptions...)

	mcpServer.AddTool(createJobTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get job name
		jobName, err := s.getMCPFieldValue("job_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get job image
		jobImage, err := s.getMCPFieldValue("job_image", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get privileged
		privileged, err := s.getMCPBooleanFieldValue("job_privileged", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get CPU
		cpu, _ := s.getMCPFieldValue("job_cpu", request)
		if cpu == "" {
			cpu = "0.1"
		}

		// Get description
		description, _ := s.getMCPFieldValue("job_description", request)

		// Get environment variables
		environmentVariables, _ := s.getMCPFieldValue("job_environment_variables", request)

		// Get command
		commandStr, _ := s.getMCPFieldValue("job_command", request)
		var command []string
		if commandStr != "" {
			// Split by comma
			command = strings.Split(commandStr, ",")
			// Trim spaces from each command
			for i, cmd := range command {
				command[i] = strings.TrimSpace(cmd)
			}
		}

		// Get args
		argsStr, _ := s.getMCPFieldValue("job_args", request)
		var args []string
		if argsStr != "" {
			// Split by comma
			args = strings.Split(argsStr, ",")
			// Trim spaces from each arg
			for i, arg := range args {
				args[i] = strings.TrimSpace(arg)
			}
		}

		// Get retry count
		retryCountStr, _ := s.getMCPFieldValue("job_retry_count", request)
		var retryCount uint32
		if retryCountStr != "" {
			if val, err := strconv.ParseUint(retryCountStr, 10, 32); err == nil {
				retryCount = uint32(val)
			}
		}

		// Get execution timeout
		executionTimeoutStr, _ := s.getMCPFieldValue("job_execution_timeout", request)
		var executionTimeout uint32
		if executionTimeoutStr != "" {
			if val, err := strconv.ParseUint(executionTimeoutStr, 10, 32); err == nil {
				executionTimeout = uint32(val)
			}
		}

		// Get run immediately
		runImmediately, err := s.getMCPBooleanFieldValue("job_run_immediately", request)
		if err != nil {
			runImmediately = false // default value
		}

		// Create the request struct
		createRequest := domain.CreateJobRequest{
			ProjectID:               projectID,
			JobName:                 jobName,
			JobImage:                jobImage,
			JobPrivileged:           privileged,
			JobCPU:                  cpu,
			JobDescription:          description,
			JobEnvironmentVariables: environmentVariables,
			JobCommand:              command,
			JobArgs:                 args,
			JobRetryCount:           retryCount,
			JobExecutionTimeout:     executionTimeout,
			JobRunImmediately:       runImmediately,
		}

		// Call the service
		operation, err := s.jobsService.CreateJob(createRequest)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(operation, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created Job: %s\n%s", jobName, string(result))), nil
	})
}
