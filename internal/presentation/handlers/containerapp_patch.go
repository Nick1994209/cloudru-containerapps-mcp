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
		containerAppPortStr, _ := s.getMCPFieldValue("containerapp_port", request)

		// Convert port to integer pointer
		var containerAppPort *int
		if containerAppPortStr != "" {
			var port int
			n, err := fmt.Sscanf(containerAppPortStr, "%d", &port)
			if err != nil || n != 1 {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse containerapp_port: %s", containerAppPortStr)), nil
			}
			containerAppPort = &port
		}

		// Get container app image
		containerAppImage, _ := s.getMCPFieldValue("containerapp_image", request)

		// Get auto deployments enabled
		autoDeploymentsEnabledStr, _ := s.getMCPFieldValue("containerapp_auto_deployments_enabled", request)
		var autoDeploymentsEnabled *bool
		if autoDeploymentsEnabledStr != "" {
			var enabled bool
			n, err := fmt.Sscanf(autoDeploymentsEnabledStr, "%t", &enabled)
			if err != nil || n != 1 {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse containerapp_auto_deployments_enabled: %s", autoDeploymentsEnabledStr)), nil
			}
			autoDeploymentsEnabled = &enabled
		}

		// Get auto deployments pattern
		autoDeploymentsPattern, _ := s.getMCPFieldValue("containerapp_auto_deployments_pattern", request)

		// Get idle timeout
		idleTimeout, _ := s.getMCPFieldValue("containerapp_idle_timeout", request)

		// Get timeout
		timeout, _ := s.getMCPFieldValue("containerapp_timeout", request)

		// Get CPU
		cpu, _ := s.getMCPFieldValue("containerapp_cpu", request)

		// Get min instance count
		minInstanceCountStr, _ := s.getMCPFieldValue("containerapp_min_instance_count", request)
		var minInstanceCount *int
		if minInstanceCountStr != "" {
			var count int
			n, err := fmt.Sscanf(minInstanceCountStr, "%d", &count)
			if err != nil || n != 1 {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse containerapp_min_instance_count: %s", minInstanceCountStr)), nil
			}
			minInstanceCount = &count
		}

		// Get max instance count
		maxInstanceCountStr, _ := s.getMCPFieldValue("containerapp_max_instance_count", request)
		var maxInstanceCount *int
		if maxInstanceCountStr != "" {
			var count int
			n, err := fmt.Sscanf(maxInstanceCountStr, "%d", &count)
			if err != nil || n != 1 {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse containerapp_max_instance_count: %s", maxInstanceCountStr)), nil
			}
			maxInstanceCount = &count
		}

		// Get description
		description, _ := s.getMCPFieldValue("containerapp_description", request)

		// Get publicly accessible
		publiclyAccessibleStr, _ := s.getMCPFieldValue("containerapp_publicly_accessible", request)
		var publiclyAccessible *bool
		if publiclyAccessibleStr != "" {
			var accessible bool
			n, err := fmt.Sscanf(publiclyAccessibleStr, "%t", &accessible)
			if err != nil || n != 1 {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse containerapp_publicly_accessible: %s", publiclyAccessibleStr)), nil
			}
			publiclyAccessible = &accessible
		}

		// Get protocol
		protocol, _ := s.getMCPFieldValue("containerapp_protocol", request)

		// Get environment variables
		environmentVariables, _ := s.getMCPFieldValue("containerapp_environment_variables", request)

		// Get command
		commandStr, _ := s.getMCPFieldValue("containerapp_command", request)
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
		argsStr, _ := s.getMCPFieldValue("containerapp_args", request)
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
					return containerAppPort
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
					return autoDeploymentsEnabled
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
					return minInstanceCount
				}
				return nil
			}(),
			MaxInstanceCount: func() *int {
				if checkRequestHasKey(request, "containerapp_max_instance_count") {
					return maxInstanceCount
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
					return publiclyAccessible
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
