package presentation

import (
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/presentation/handlers"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer wraps the handlers.MCPServer for backward compatibility
type MCPServer struct {
	*handlers.MCPServer
}

// NewMCPServer creates a new MCP server with the required services
func NewMCPServer(descriptionService domain.DescriptionService, dockerService domain.DockerService, containerAppsService domain.ContainerAppsService, dockerRegistryService domain.ArtifactRegistryService) *MCPServer {
	return &MCPServer{
		MCPServer: handlers.NewMCPServer(descriptionService, dockerService, containerAppsService, dockerRegistryService),
	}
}

// RegisterAllTools registers all tools with the MCP server
func (s *MCPServer) RegisterAllTools(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterAllTools(mcpServer)
}

// RegisterDescriptionTool registers the description tool with the MCP server
func (s *MCPServer) RegisterDescriptionTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterDescriptionTool(mcpServer)
}

// RegisterDockerLoginTool registers the docker login tool with the MCP server
func (s *MCPServer) RegisterDockerLoginTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterDockerLoginTool(mcpServer)
}

// RegisterDockerBuildAndPushTool registers the docker build and push tool with the MCP server
func (s *MCPServer) RegisterDockerBuildAndPushTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterDockerBuildAndPushTool(mcpServer)
}

// RegisterGetListContainerAppsTool registers the get list container apps tool with the MCP server
func (s *MCPServer) RegisterGetListContainerAppsTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterGetListContainerAppsTool(mcpServer)
}

// RegisterGetContainerAppTool registers the get container app tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterGetContainerAppTool(mcpServer)
}

// RegisterPatchContainerAppTool registers the patch container app tool with the MCP server
func (s *MCPServer) RegisterPatchContainerAppTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterPatchContainerAppTool(mcpServer)
}

// RegisterCreateContainerAppTool registers the create container app tool with the MCP server
func (s *MCPServer) RegisterCreateContainerAppTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterCreateContainerAppTool(mcpServer)
}

// RegisterDeleteContainerAppTool registers the delete container app tool with the MCP server
func (s *MCPServer) RegisterDeleteContainerAppTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterDeleteContainerAppTool(mcpServer)
}

// RegisterStartContainerAppTool registers the start container app tool with the MCP server
func (s *MCPServer) RegisterStartContainerAppTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterStartContainerAppTool(mcpServer)
}

// RegisterStopContainerAppTool registers the stop container app tool with the MCP server
func (s *MCPServer) RegisterStopContainerAppTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterStopContainerAppTool(mcpServer)
}

// RegisterGetContainerAppLogsTool registers the get container app logs tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppLogsTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterGetContainerAppLogsTool(mcpServer)
}

// RegisterGetContainerAppSystemLogsTool registers the get container app system logs tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppSystemLogsTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterGetContainerAppSystemLogsTool(mcpServer)
}

// RegisterGetListDockerRegistriesTool registers the get list docker registries tool with the MCP server
func (s *MCPServer) RegisterGetListDockerRegistriesTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterGetListDockerRegistriesTool(mcpServer)
}

// RegisterCreateDockerRegistryTool registers the create docker registry tool with the MCP server
func (s *MCPServer) RegisterCreateDockerRegistryTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterCreateDockerRegistryTool(mcpServer)
}

// RegisterGetRegistryImagesTool registers the get registry images tool with the MCP server
func (s *MCPServer) RegisterGetRegistryImagesTool(mcpServer *server.MCPServer) {
	s.MCPServer.RegisterGetRegistryImagesTool(mcpServer)
}
