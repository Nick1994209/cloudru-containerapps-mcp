package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetContainerAppTool registers the get container app tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"containerapp_name",
	)
	getContainerAppTool := mcp.NewTool("cloudru_get_containerapp", toolOptions...)

	mcpServer.AddTool(getContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		containerApp, err := s.containerAppsService.GetContainerApp(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(containerApp, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}
