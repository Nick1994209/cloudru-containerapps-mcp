package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterDockerLoginTool registers the docker login tool with the MCP server
func (s *MCPServer) RegisterDockerLoginTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions("Login to Cloud.ru Artifact registry (Docker registry)", "registry_name")
	dockerLoginTool := mcp.NewTool("cloudru_docker_login", toolOptions...)

	mcpServer.AddTool(dockerLoginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Using helper functions for type-safe argument access
		registryName, err := request.RequireString("registry_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result, err := s.dockerService.Login(registryName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Successfully login to Cloud.ru Artifact Registry: %s", result)), nil
	})
}
