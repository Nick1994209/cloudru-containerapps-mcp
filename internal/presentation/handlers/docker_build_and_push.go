package handlers

import (
	"context"
	"fmt"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterDockerBuildAndPushTool registers the docker build and push tool with the MCP server
func (s *MCPServer) RegisterDockerBuildAndPushTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Build and push Docker image to Cloud.ru Artifact Registry (Docker registry)",
		"registry_name",
		"repository_name",
		"image_version",
		"dockerfile_path",
		"dockerfile_target",
		"dockerfile_folder",
		"show_commands",
	)
	dockerPushTool := mcp.NewTool("cloudru_docker_build_and_push", toolOptions...)

	mcpServer.AddTool(dockerPushTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		registryName, err := s.getMCPFieldValue("registry_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		repositoryName, err := s.getMCPFieldValue("repository_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		imageVersion, _ := s.getMCPFieldValue("image_version", request)
		dockerfilePath, _ := request.RequireString("dockerfile_path")
		dockerfileTarget, _ := request.RequireString("dockerfile_target")
		dockerfileFolder, _ := request.RequireString("dockerfile_folder")

		image := domain.DockerImage{
			RegistryName:     registryName,
			RepositoryName:   repositoryName,
			ImageVersion:     imageVersion,
			DockerfilePath:   dockerfilePath,
			DockerfileTarget: dockerfileTarget,
			DockerfileFolder: dockerfileFolder,
		}

		// Determine whether to execute build/push or just return commands
		showCommands, err := s.getMCPBooleanFieldValue("show_commands", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if showCommands {
			buildCmd, pushCmd, err := s.dockerService.ShowBuildAndPushCommands(image)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			combined := fmt.Sprintf("Run Docker build command:\n'%s'\n and then run docker push command:\n'%s'. IMPORTANT! Use platform for building image.", buildCmd, pushCmd)
			return mcp.NewToolResultText(combined), nil
		}

		result, err := s.dockerService.BuildAndPush(image)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully built and pushed Docker image to Cloud.ru Artifact Registry: %s", result)), nil
	})
}
