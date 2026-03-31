package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterDescriptionTool registers the description tool with the MCP server
func (s *MCPServer) RegisterDescriptionTool(mcpServer *server.MCPServer) {
	descriptionTool := mcp.NewTool("cloudru_containerapps_description",
		mcp.WithDescription("Returns usage instructions for Cloud.ru Container Apps MCP"),
	)

	mcpServer.AddTool(descriptionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(s.descriptionService.GetDescription()), nil
	})
}
