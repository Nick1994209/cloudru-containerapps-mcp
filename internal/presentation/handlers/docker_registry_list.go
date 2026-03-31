package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterGetListDockerRegistriesTool registers the get list docker registries tool with the MCP server
func (s *MCPServer) RegisterGetListDockerRegistriesTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get list of Docker Registries from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
	)
	getListDockerRegistriesTool := mcp.NewTool("cloudru_get_list_docker_registries", toolOptions...)

	mcpServer.AddTool(getListDockerRegistriesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the service
		dockerRegistries, err := s.dockerRegistryService.GetListDockerRegistries(projectID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(dockerRegistries, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}
