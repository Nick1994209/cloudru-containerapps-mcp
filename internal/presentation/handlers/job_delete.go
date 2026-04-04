package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterDeleteJobTool registers the delete job tool with the MCP server
func (s *MCPServer) RegisterDeleteJobTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Delete a Job from Cloud.ru. WARNING: This action cannot be undone!",
		"project_id",
		"job_name",
	)
	deleteJobTool := mcp.NewTool("cloudru_delete_job", toolOptions...)

	mcpServer.AddTool(deleteJobTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Confirmation prompt - in MCP context, we'll add a warning in the description
		// but the actual confirmation would typically happen in the client UI

		// Call the service
		operation, err := s.jobsService.DeleteJob(projectID, jobName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(operation, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted Job: %s\n%s", jobName, string(result))), nil
	})
}
