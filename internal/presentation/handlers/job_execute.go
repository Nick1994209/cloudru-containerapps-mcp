package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterExecuteJobTool registers the execute job tool with the MCP server
func (s *MCPServer) RegisterExecuteJobTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Execute a Job in Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"job_name",
		"params",
	)
	executeJobTool := mcp.NewTool("cloudru_execute_job", toolOptions...)

	mcpServer.AddTool(executeJobTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Get optional params
		paramsStr, _ := s.getMCPFieldValue("params", request)

		// Parse params if provided
		params := make(map[string]interface{})
		if paramsStr != "" {
			if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse params: %v", err)), nil
			}
		}

		// Call the service
		jobExecution, err := s.jobsService.ExecuteJob(projectID, jobName, params)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(jobExecution, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}
