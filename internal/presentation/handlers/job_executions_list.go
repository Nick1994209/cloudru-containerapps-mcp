package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetListExecutionsTool registers the get list executions tool with the MCP server
func (s *MCPServer) RegisterGetListExecutionsTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get paginated list of job executions from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"job_name",
		"page_size",
	)

	getListExecutionsTool := mcp.NewTool("cloudru_job_executions_list", toolOptions...)

	mcpServer.AddTool(getListExecutionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Get optional parameters
		pageSize, _ := s.getMCPFieldValue("page_size", request)

		// Log the parameters for debugging
		fmt.Printf("GetListExecutions called with projectID: %s, jobName: %s, pageSize: %s\n",
			projectID, jobName, pageSize)

		// Call the service
		executions, err := s.jobsService.GetListExecutions(projectID, jobName, pageSize)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(executions, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}
