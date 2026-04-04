package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetListJobsTool registers the get list jobs tool with the MCP server
func (s *MCPServer) RegisterGetListJobsTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get paginated list of jobs from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"page_size",
	)

	getListJobsTool := mcp.NewTool("cloudru_jobs_list", toolOptions...)

	mcpServer.AddTool(getListJobsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get optional parameters
		pageSize, _ := s.getMCPFieldValue("page_size", request)

		// Log the parameters for debugging
		fmt.Printf("GetListJobs called with projectID: %s, pageSize: %s\n",
			projectID, pageSize)

		// Call the service
		jobs, err := s.jobsService.GetListJobs(projectID, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(jobs, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}
