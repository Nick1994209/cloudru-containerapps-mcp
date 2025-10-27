package main

import (
	"fmt"
	"log"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/presentation"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create infrastructure layer
	dockerInfrastructure := application.NewDockerApplication(cfg)
	containerAppsService := cloudru.NewContainerAppsApplication(cfg)
	dockerRegistryService := cloudru.NewArtifactRegistryApplication(cfg)

	// Create application layer
	descriptionService := application.NewDescriptionApplication()

	// Log the application description
	log.Println("Application Description:")
	log.Println(descriptionService.GetDescription())

	// Create presentation layer
	mcpServer := presentation.NewMCPServer(descriptionService, dockerInfrastructure, containerAppsService, dockerRegistryService)

	// Create a new MCP server
	s := server.NewMCPServer(
		"Cloud.ru Container Apps MCP",
		"0.0.1",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// Register tools with the MCP server
	mcpServer.RegisterDescriptionTool(s)
	mcpServer.RegisterDockerLoginTool(s)
	mcpServer.RegisterDockerBuildAndPushTool(s)
	mcpServer.RegisterGetListContainerAppsTool(s)
	mcpServer.RegisterGetContainerAppTool(s)
	mcpServer.RegisterCreateContainerAppTool(s)
	mcpServer.RegisterDeleteContainerAppTool(s)
	mcpServer.RegisterStartContainerAppTool(s)
	mcpServer.RegisterStopContainerAppTool(s)
	mcpServer.RegisterGetContainerAppLogsTool(s)
	// mcpServer.RegisterGetContainerAppSystemLogsTool(s)
	mcpServer.RegisterGetListDockerRegistriesTool(s)
	mcpServer.RegisterCreateDockerRegistryTool(s)
	// mcpServer.RegisterGetRegistryImagesTool(s)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
