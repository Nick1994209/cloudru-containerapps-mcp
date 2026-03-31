package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterDeleteContainerAppTool registers the delete container app tool with the MCP server
func (s *MCPServer) RegisterDeleteContainerAppTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Delete a Container App from Cloud.ru. WARNING: This action cannot be undone!",
		"project_id",
		"containerapp_name",
	)
	deleteContainerAppTool := mcp.NewTool("cloudru_delete_containerapp", toolOptions...)

	mcpServer.AddTool(deleteContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Confirmation prompt - in MCP context, we'll add a warning in the description
		// but the actual confirmation would typically happen in the client UI

		// Call the service
		operation, err := s.containerAppsService.DeleteContainerApp(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(operation, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted Container App: %s\n%s", containerAppName, string(result))), nil
	})
}
