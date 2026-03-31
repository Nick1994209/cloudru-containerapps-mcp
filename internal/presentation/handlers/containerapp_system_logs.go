package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetContainerAppSystemLogsTool registers the get container app system logs tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppSystemLogsTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get system logs for a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"containerapp_name",
	)
	getContainerAppSystemLogsTool := mcp.NewTool("cloudru_get_containerapp_system_logs", toolOptions...)

	mcpServer.AddTool(getContainerAppSystemLogsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get container app name
		containerAppName, err := s.getMCPFieldValue("containerapp_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the service
		containerAppSystemLogs, err := s.containerAppsService.GetContainerAppSystemLogs(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output with limit of 200 records
		limitedData, err := marshalJSONWithLimit(containerAppSystemLogs.Data, 200)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		// Create the response structure with limited data
		response := map[string]interface{}{
			"data": limitedData,
		}

		// Marshal the final response
		result, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}
