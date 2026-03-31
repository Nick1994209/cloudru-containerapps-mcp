package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterCreateDockerRegistryTool registers the create docker registry tool with the MCP server
func (s *MCPServer) RegisterCreateDockerRegistryTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Create a new Docker Registry in Cloud.ru",
		"project_id",
		"registry_name",
		"registry_is_public",
	)
	createDockerRegistryTool := mcp.NewTool("cloudru_create_docker_registry", toolOptions...)

	mcpServer.AddTool(createDockerRegistryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get registry name
		registryName, err := s.getMCPFieldValue("registry_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get registry_is_public flag
		isPublic, err := s.getMCPBooleanFieldValue("registry_is_public", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the service
		dockerRegistry, err := s.dockerRegistryService.CreateDockerRegistry(projectID, registryName, isPublic)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(dockerRegistry, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created Docker Registry: %s\n%s", registryName, string(result))), nil
	})
}
