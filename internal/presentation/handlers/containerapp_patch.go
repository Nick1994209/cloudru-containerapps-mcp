package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterPatchContainerAppTool registers the patch container app tool with the MCP server
func (s *MCPServer) RegisterPatchContainerAppTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Patch a Container App in Cloud.ru. This will get the current state, merge with new values, and update the container app.",
		"project_id",
		"containerapp_name",
		"containerapp_port",
		"containerapp_image",
		"containerapp_auto_deployments_enabled",
		"containerapp_auto_deployments_pattern",
		"containerapp_idle_timeout",
		"containerapp_timeout",
		"containerapp_cpu",
		"containerapp_min_instance_count",
		"containerapp_max_instance_count",
		"containerapp_description",
		"containerapp_publicly_accessible",
		"containerapp_protocol",
		"containerapp_environment_variables",
		"containerapp_command",
		"containerapp_args",
	)
	patchContainerAppTool := mcp.NewTool("cloudru_patch_containerapp", toolOptions...)

	mcpServer.AddTool(patchContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Get container app port
		containerAppPortStr, err := s.getMCPFieldValue("containerapp_port", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert port to integer
		var containerAppPort int
		fmt.Sscanf(containerAppPortStr, "%d", &containerAppPort)

		// Get container app image
		containerAppImage, err := s.getMCPFieldValue("containerapp_image", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get auto deployments enabled
		autoDeploymentsEnabled, err := s.getMCPBooleanFieldValue("containerapp_auto_deployments_enabled", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get auto deployments pattern
		autoDeploymentsPattern, err := s.getMCPFieldValue("containerapp_auto_deployments_pattern", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get idle timeout
		idleTimeout, err := s.getMCPFieldValue("containerapp_idle_timeout", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get timeout
		timeout, err := s.getMCPFieldValue("containerapp_timeout", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get CPU
		cpu, err := s.getMCPFieldValue("containerapp_cpu", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get min instance count
		minInstanceCountStr, err := s.getMCPFieldValue("containerapp_min_instance_count", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var minInstanceCount int
		if minInstanceCountStr != "" {
			fmt.Sscanf(minInstanceCountStr, "%d", &minInstanceCount)
		}

		// Get max instance count
		maxInstanceCountStr, err := s.getMCPFieldValue("containerapp_max_instance_count", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var maxInstanceCount int
		if maxInstanceCountStr != "" {
			fmt.Sscanf(maxInstanceCountStr, "%d", &maxInstanceCount)
		}

		// Get description
		description, err := s.getMCPFieldValue("containerapp_description", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get publicly accessible
		publiclyAccessible, err := s.getMCPBooleanFieldValue("containerapp_publicly_accessible", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get protocol
		protocol, err := s.getMCPFieldValue("containerapp_protocol", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get environment variables
		environmentVariables, err := s.getMCPFieldValue("containerapp_environment_variables", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get command
		commandStr, err := s.getMCPFieldValue("containerapp_command", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var command []string
		if commandStr != "" {
			// Split by comma
			command = strings.Split(commandStr, ",")
			// Trim spaces from each command
			for i, cmd := range command {
				command[i] = strings.TrimSpace(cmd)
			}
		}

		// Get args
		argsStr, err := s.getMCPFieldValue("containerapp_args", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var args []string
		if argsStr != "" {
			// Split by comma
			args = strings.Split(argsStr, ",")
			// Trim spaces from each arg
			for i, arg := range args {
				args[i] = strings.TrimSpace(arg)
			}
		}

		// Create the patch request
		patchRequest := domain.PatchContainerAppRequest{
			ProjectID:        projectID,
			ContainerAppName: containerAppName,
			ContainerAppPort: func() *int {
				if checkRequestHasKey(request, "containerapp_port") {
					return &containerAppPort
				}
				return nil
			}(),
			ContainerAppImage: func() *string {
				if checkRequestHasKey(request, "containerapp_image") {
					return &containerAppImage
				}
				return nil
			}(),
			AutoDeploymentsEnabled: func() *bool {
				if checkRequestHasKey(request, "containerapp_auto_deployments_enabled") {
					return &autoDeploymentsEnabled
				}
				return nil
			}(),
			AutoDeploymentsPattern: func() *string {
				if checkRequestHasKey(request, "containerapp_auto_deployments_pattern") {
					return &autoDeploymentsPattern
				}
				return nil
			}(),
			IdleTimeout: func() *string {
				if checkRequestHasKey(request, "containerapp_idle_timeout") {
					return &idleTimeout
				}
				return nil
			}(),
			Timeout: func() *string {
				if checkRequestHasKey(request, "containerapp_timeout") {
					return &timeout
				}
				return nil
			}(),
			CPU: func() *string {
				if checkRequestHasKey(request, "containerapp_cpu") {
					return &cpu
				}
				return nil
			}(),
			MinInstanceCount: func() *int {
				if checkRequestHasKey(request, "containerapp_min_instance_count") {
					return &minInstanceCount
				}
				return nil
			}(),
			MaxInstanceCount: func() *int {
				if checkRequestHasKey(request, "containerapp_max_instance_count") {
					return &maxInstanceCount
				}
				return nil
			}(),
			Description: func() *string {
				if checkRequestHasKey(request, "containerapp_description") {
					return &description
				}
				return nil
			}(),
			PubliclyAccessible: func() *bool {
				if checkRequestHasKey(request, "containerapp_publicly_accessible") {
					return &publiclyAccessible
				}
				return nil
			}(),
			Protocol: func() *string {
				if checkRequestHasKey(request, "containerapp_protocol") {
					return &protocol
				}
				return nil
			}(),
			EnvironmentVariables: func() *string {
				if checkRequestHasKey(request, "containerapp_environment_variables") {
					return &environmentVariables
				}
				return nil
			}(),
			Command: func() []string {
				if checkRequestHasKey(request, "containerapp_command") {
					return command
				}
				return nil
			}(),
			Args: func() []string {
				if checkRequestHasKey(request, "containerapp_args") {
					return args
				}
				return nil
			}(),
		}

		// Call the service
		operation, err := s.containerAppsService.PatchContainerApp(projectID, containerAppName, patchRequest)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(operation, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully patched Container App: %s\n%s", containerAppName, string(result))), nil
	})
}
