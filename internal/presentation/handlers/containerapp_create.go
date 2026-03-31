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

// RegisterCreateContainerAppTool registers the create container app tool with the MCP server
func (s *MCPServer) RegisterCreateContainerAppTool(mcpServer *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Create a new Container App in Cloud.ru",
		"project_id",
		"containerapp_name",
		"containerapp_port",
		"containerapp_image",
		"containerapp_auto_deployments_enabled",
		"containerapp_auto_deployments_pattern",
		"containerapp_privileged",
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
	createContainerAppTool := mcp.NewTool("cloudru_create_containerapp", toolOptions...)

	mcpServer.AddTool(createContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		autoDeploymentsPattern, _ := s.getMCPFieldValue("containerapp_auto_deployments_pattern", request)

		// Get privileged
		privileged, err := s.getMCPBooleanFieldValue("containerapp_privileged", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get idle timeout
		idleTimeout, _ := s.getMCPFieldValue("containerapp_idle_timeout", request)
		if idleTimeout == "" {
			idleTimeout = "600s"
		}

		// Get timeout
		timeout, _ := s.getMCPFieldValue("containerapp_timeout", request)
		if timeout == "" {
			timeout = "60s"
		}

		// Get CPU
		cpu, _ := s.getMCPFieldValue("containerapp_cpu", request)
		if cpu == "" {
			cpu = "0.1"
		}

		// Get min instance count
		minInstanceCountStr, _ := s.getMCPFieldValue("containerapp_min_instance_count", request)
		var minInstanceCount int
		if minInstanceCountStr != "" {
			fmt.Sscanf(minInstanceCountStr, "%d", &minInstanceCount)
		}

		// Get max instance count
		maxInstanceCountStr, _ := s.getMCPFieldValue("containerapp_max_instance_count", request)
		var maxInstanceCount int
		if maxInstanceCountStr != "" {
			fmt.Sscanf(maxInstanceCountStr, "%d", &maxInstanceCount)
		}

		// Get description
		description, _ := s.getMCPFieldValue("containerapp_description", request)

		// Get publicly accessible
		publiclyAccessible, err := s.getMCPBooleanFieldValue("containerapp_publicly_accessible", request)
		if err != nil {
			publiclyAccessible = true // default value
		}

		// Get protocol
		protocol, _ := s.getMCPFieldValue("containerapp_protocol", request)
		if protocol == "" {
			protocol = "http_1"
		}

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

		// Create the request struct
		createRequest := domain.CreateContainerAppRequest{
			ProjectID:              projectID,
			ContainerAppName:       containerAppName,
			ContainerAppPort:       containerAppPort,
			ContainerAppImage:      containerAppImage,
			AutoDeploymentsEnabled: autoDeploymentsEnabled,
			AutoDeploymentsPattern: autoDeploymentsPattern,
			Privileged:             privileged,
			IdleTimeout:            idleTimeout,
			Timeout:                timeout,
			CPU:                    cpu,
			MinInstanceCount:       minInstanceCount,
			MaxInstanceCount:       maxInstanceCount,
			Description:            description,
			PubliclyAccessible:     publiclyAccessible,
			Protocol:               protocol,
			EnvironmentVariables:   environmentVariables,
			Command:                command,
			Args:                   args,
		}

		// Call the service
		operation, err := s.containerAppsService.CreateContainerApp(createRequest)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(operation, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created Container App: %s\n%s", containerAppName, string(result))), nil
	})
}
