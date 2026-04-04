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

// RegisterPatchJobTool registers the patch job tool with the MCP server
func (s *MCPServer) RegisterPatchJobTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Patch a Job in Cloud.ru. This will get the current state, merge with new values, and update the job.",
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
	patchJobTool := mcp.NewTool("cloudru_patch_job", toolOptions...)

	mcpServer.AddTool(patchJobTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		jobImage, _ := s.getMCPFieldValue("job_image", request)

		// Get privileged
		privilegedStr, _ := s.getMCPFieldValue("job_privileged", request)
		var privileged *bool
		if privilegedStr != "" {
			var priv bool
			n, err := fmt.Sscanf(privilegedStr, "%t", &priv)
			if err != nil || n != 1 {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse job_privileged: %s", privilegedStr)), nil
			}
			privileged = &priv
		}

		// Get CPU
		cpu, _ := s.getMCPFieldValue("job_cpu", request)

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
		var retryCount *uint32
		if retryCountStr != "" {
			if val, err := strconv.ParseUint(retryCountStr, 10, 32); err == nil {
				count := uint32(val)
				retryCount = &count
			}
		}

		// Get execution timeout
		executionTimeoutStr, _ := s.getMCPFieldValue("job_execution_timeout", request)
		var executionTimeout *uint32
		if executionTimeoutStr != "" {
			if val, err := strconv.ParseUint(executionTimeoutStr, 10, 32); err == nil {
				timeout := uint32(val)
				executionTimeout = &timeout
			}
		}

		// Get run immediately
		runImmediatelyStr, _ := s.getMCPFieldValue("job_run_immediately", request)
		var runImmediately *bool
		if runImmediatelyStr != "" {
			var run bool
			n, err := fmt.Sscanf(runImmediatelyStr, "%t", &run)
			if err != nil || n != 1 {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse job_run_immediately: %s", runImmediatelyStr)), nil
			}
			runImmediately = &run
		}

		// Create the patch request
		patchRequest := domain.PatchJobRequest{
			ProjectID: projectID,
			JobName:   jobName,
			JobImage: func() *string {
				if checkRequestHasKey(request, "job_image") {
					return &jobImage
				}
				return nil
			}(),
			JobPrivileged: func() *bool {
				if checkRequestHasKey(request, "job_privileged") {
					return privileged
				}
				return nil
			}(),
			JobCPU: func() *string {
				if checkRequestHasKey(request, "job_cpu") {
					return &cpu
				}
				return nil
			}(),
			JobDescription: func() *string {
				if checkRequestHasKey(request, "job_description") {
					return &description
				}
				return nil
			}(),
			JobEnvironmentVariables: func() *string {
				if checkRequestHasKey(request, "job_environment_variables") {
					return &environmentVariables
				}
				return nil
			}(),
			JobCommand: func() []string {
				if checkRequestHasKey(request, "job_command") {
					return command
				}
				return nil
			}(),
			JobArgs: func() []string {
				if checkRequestHasKey(request, "job_args") {
					return args
				}
				return nil
			}(),
			JobRetryCount: func() *uint32 {
				if checkRequestHasKey(request, "job_retry_count") {
					return retryCount
				}
				return nil
			}(),
			JobExecutionTimeout: func() *uint32 {
				if checkRequestHasKey(request, "job_execution_timeout") {
					return executionTimeout
				}
				return nil
			}(),
			JobRunImmediately: func() *bool {
				if checkRequestHasKey(request, "job_run_immediately") {
					return runImmediately
				}
				return nil
			}(),
		}

		// Call the service
		operation, err := s.jobsService.PatchJob(projectID, jobName, patchRequest)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(operation, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully patched Job: %s\n%s", jobName, string(result))), nil
	})
}
