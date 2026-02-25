package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// marshalJSONWithLimit marshals sequential data with a limit on the number of records
func marshalJSONWithLimit(data interface{}, maxRecords int) ([]byte, error) {
	// Use reflection to check if data is a slice or array
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		// If it's a slice/array, limit the number of elements
		if val.Len() > maxRecords {
			// Create a new slice with limited elements
			limitedData := make([]interface{}, maxRecords)
			for i := 0; i < maxRecords; i++ {
				limitedData[i] = val.Index(i).Interface()
			}
			return json.MarshalIndent(limitedData, "", "  ")
		}
	}

	// If not a slice/array or length is within limit, marshal normally
	return json.MarshalIndent(data, "", "  ")
}

// MCPServer holds the application services
type MCPServer struct {
	descriptionService    domain.DescriptionService
	dockerService         domain.DockerService
	containerAppsService  domain.ContainerAppsService
	dockerRegistryService domain.ArtifactRegistryService

	mappedFields map[string]struct {
		envValue     string
		description  string
		defaultValue string
		title        string
		required     bool
	}
	cfg *config.Config
}

// NewMCPServer creates a new MCP server with the required services
func NewMCPServer(descriptionService domain.DescriptionService, dockerService domain.DockerService, containerAppsService domain.ContainerAppsService, dockerRegistryService domain.ArtifactRegistryService) *MCPServer {
	cfg := config.LoadConfig()

	defaultRepoName := cfg.CurrentDir
	if cfg.DockerfileTarget != "" && cfg.DockerfileTarget != "-" {
		defaultRepoName = defaultRepoName + "-" + cfg.DockerfileTarget
	}

	containerappImage := fmt.Sprintf("%s.%s/%s:%s", cfg.RegistryName, cfg.RegistryDomain, cfg.RepositoryName, "latest")
	return &MCPServer{
		descriptionService:    descriptionService,
		dockerService:         dockerService,
		containerAppsService:  containerAppsService,
		dockerRegistryService: dockerRegistryService,
		cfg:                   cfg,

		mappedFields: map[string]struct {
			envValue     string
			description  string
			defaultValue string
			title        string
			required     bool
		}{
			"project_id": {
				envValue:    cfg.ProjectID,
				description: "Project ID for Container Apps (can be set via PROJECT_ID environment variable)",
				required:    true,
			},
			"registry_name": {
				envValue:    cfg.RegistryName,
				description: "Registry name",
				required:    true,
			},
			"registry_is_public": {
				description:  "Make registry public",
				required:     false,
				defaultValue: "false",
			},
			"repository_name": {
				envValue:     cfg.RepositoryName,
				description:  "Repository name",
				defaultValue: defaultRepoName,
				required:     true,
			},
			"image_version": {
				description:  "Image version",
				title:        "For example: latest or v0.0.1",
				required:     true,
				defaultValue: "latest",
			},
			"show_commands": {
				description:  "If true, return Docker build and push commands without executing them",
				defaultValue: "true",
				required:     false,
			},
			"dockerfile_path": {
				envValue:     cfg.Dockerfile,
				description:  "Repository name",
				defaultValue: "Dockerfile",
				required:     false,
			},
			"dockerfile_target": {
				envValue:     cfg.DockerfileTarget,
				description:  "Dockerfile target stage",
				defaultValue: "-",
				required:     false,
			},
			"dockerfile_folder": {
				envValue:     cfg.DockerfileFolder,
				description:  "Dockerfile folder (build context)",
				defaultValue: ".",
				required:     false,
			},
			"containerapp_name": {
				envValue:     cfg.ContainerAppName,
				description:  "Container App name (can be set via CONTAINERAPP_NAME environment variable)",
				required:     false,
				defaultValue: cfg.CurrentDir,
				title:        "You can use example: " + cfg.CurrentDir,
			},
			"containerapp_port": {
				description: "Container App port number",
				required:    true,
				title:       "You can use example: 8000",
			},
			"containerapp_image": {
				description: "Container App image",
				required:    true,
				title:       "Example image: " + containerappImage,
			},
			"containerapp_auto_deployments_enabled": {
				description:  "Enable auto deployments",
				defaultValue: "false",
				required:     false,
			},
			"containerapp_auto_deployments_pattern": {
				description:  "Auto deployments pattern",
				defaultValue: "latest",
				required:     false,
			},
			"containerapp_privileged": {
				description:  "Run container in privileged mode",
				defaultValue: "false",
				required:     false,
			},
			"containerapp_idle_timeout": {
				description:  "Parameter defines how long a service stays active without receiving any requests before being shut down.",
				defaultValue: "600s",
				required:     false,
			},
			"containerapp_timeout": {
				description:  "Parameter that defines the maximum amount of time allowed for processing a request. If a complete response is not generated and sent within this period, the request is terminated.",
				defaultValue: "60s",
				required:     false,
			},
			"containerapp_cpu": {
				description:  "CPU allocation (0.1 CPU - 256 Mi RAM, 0.2 CPU - 512 Mi RAM, ...)",
				defaultValue: "0.1",
				required:     false,
				title:        "Options: 0.1, 0.2, 0.5, 1",
			},
			"containerapp_min_instance_count": {
				description:  "Minimum number of instances for scaling",
				defaultValue: "0",
				required:     false,
			},
			"containerapp_max_instance_count": {
				description:  "Maximum number of instances for scaling",
				defaultValue: "1",
				required:     false,
			},
			"containerapp_description": {
				description:  "Description of the container app",
				defaultValue: "",
				required:     false,
			},
			"containerapp_publicly_accessible": {
				description:  "Whether the container app is publicly accessible",
				defaultValue: "true",
				required:     false,
			},
			"containerapp_protocol": {
				description:  "Protocol for the container app",
				defaultValue: "http_1",
				required:     false,
				title:        "Options: http_1, http_2",
			},
			"containerapp_environment_variables": {
				description:  "Environment variables in format <name>='<value>';<next_name>='value2'",
				defaultValue: "",
				required:     false,
			},
			"containerapp_command": {
				description:  "Command to run in the container (comma-separated values)",
				defaultValue: "",
				required:     false,
			},
			"containerapp_args": {
				description:  "Arguments for the command (comma-separated values)",
				defaultValue: "",
				required:     false,
			},
		},
	}
}

func (s *MCPServer) getMCPFieldsOptions(description string, fields ...string) []mcp.ToolOption {
	result := []mcp.ToolOption{
		mcp.WithDescription(description),
	}
	for _, field := range fields {
		fieldData := s.mappedFields[field]
		if fieldData.envValue == "" {
			description := fieldData.description
			if fieldData.defaultValue != "" {
				description = fmt.Sprintf("%s (default: %s)", fieldData.description, fieldData.defaultValue)
			}
			opts := []mcp.PropertyOption{
				mcp.Description(description),
			}
			if fieldData.required {
				opts = append(opts, mcp.Required())
			}
			if fieldData.title != "" {
				opts = append(opts, mcp.Title(fieldData.title))
			}
			if fieldData.defaultValue != "" {
				opts = append(opts, mcp.DefaultString(fieldData.defaultValue))
			}
			result = append(result, mcp.WithString(field, opts...))
		}
	}
	return result
}

func (s *MCPServer) getMCPFieldValue(field string, request mcp.CallToolRequest) (string, error) {
	fieldData := s.mappedFields[field]
	// If we have an environment variable value, use it
	if fieldData.envValue != "" {
		return fieldData.envValue, nil
	}

	// Try to get the value from the request
	result, err := request.RequireString(field)
	if err != nil && fieldData.defaultValue == "" {
		// If there's an error and no default or env value, return the error
		return "", err
	}

	// If we got a value from the request, use it
	if result != "" {
		return result, nil
	}

	// Otherwise, use the default value if available
	if fieldData.defaultValue != "" {
		return fieldData.defaultValue, nil
	}

	// Otherwise, use the default value if available
	if !fieldData.required {
		return "", nil
	}

	// If we get here, return whatever we have (likely empty)
	return "", fmt.Errorf("field %s is empty: %s", field, fieldData.description)
}

func (s *MCPServer) getMCPBooleanFieldValue(field string, request mcp.CallToolRequest) (bool, error) {
	fieldValueStr, _ := s.getMCPFieldValue(field, request)

	if fieldValueStr == "true" || fieldValueStr == "1" {
		return true, nil
	} else if fieldValueStr == "false" || fieldValueStr == "0" {
		return false, nil
	} else {
		return false, fmt.Errorf("field %s must be 'true', 'false', '1', or '0', got: %s", field, fieldValueStr)
	}
}

// RegisterDescriptionTool registers the description tool with the MCP server
func (s *MCPServer) RegisterDescriptionTool(server *server.MCPServer) {
	descriptionTool := mcp.NewTool("cloudru_containerapps_description",
		mcp.WithDescription("Returns usage instructions for Cloud.ru Container Apps MCP"),
	)

	server.AddTool(descriptionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(s.descriptionService.GetDescription()), nil
	})
}

// RegisterDockerLoginTool registers the docker login tool with the MCP server
func (s *MCPServer) RegisterDockerLoginTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions("Login to Cloud.ru Artifact registry (Docker registry)", "registry_name")
	dockerLoginTool := mcp.NewTool("cloudru_docker_login", toolOptions...)

	server.AddTool(dockerLoginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Using helper functions for type-safe argument access
		registryName, err := s.getMCPFieldValue("registry_name", request)
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

// RegisterDockerBuildAndPushTool registers the docker build and push tool with the MCP server
func (s *MCPServer) RegisterDockerBuildAndPushTool(server *server.MCPServer) {
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

	server.AddTool(dockerPushTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
			combined := fmt.Sprintf("Run Docker build command:\n%s\n and then run docker push command:\n%s", buildCmd, pushCmd)
			return mcp.NewToolResultText(combined), nil
		}

		result, err := s.dockerService.BuildAndPush(image)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully built and pushed Docker image to Cloud.ru Artifact Registry: %s", result)), nil
	})
}

// RegisterGetListContainerAppsTool registers the get list container apps tool with the MCP server
func (s *MCPServer) RegisterGetListContainerAppsTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get list of Container Apps from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
	)
	getListContainerAppsTool := mcp.NewTool("cloudru_get_list_containerapps", toolOptions...)

	server.AddTool(getListContainerAppsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get project ID
		projectID, err := s.getMCPFieldValue("project_id", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the service
		containerApps, err := s.containerAppsService.GetListContainerApps(projectID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(containerApps, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}

// RegisterGetContainerAppTool registers the get container app tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"containerapp_name",
	)
	getContainerAppTool := mcp.NewTool("cloudru_get_containerapp", toolOptions...)

	server.AddTool(getContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Call the service
		containerApp, err := s.containerAppsService.GetContainerApp(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(containerApp, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}

// RegisterCreateContainerAppTool registers the create container app tool with the MCP server
func (s *MCPServer) RegisterCreateContainerAppTool(server *server.MCPServer) {
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

	server.AddTool(createContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		containerApp, err := s.containerAppsService.CreateContainerApp(createRequest)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(containerApp, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully created Container App: %s\n%s", containerAppName, string(result))), nil
	})
}

// RegisterDeleteContainerAppTool registers the delete container app tool with the MCP server
func (s *MCPServer) RegisterDeleteContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Delete a Container App from Cloud.ru. WARNING: This action cannot be undone!",
		"project_id",
		"containerapp_name",
	)
	deleteContainerAppTool := mcp.NewTool("cloudru_delete_containerapp", toolOptions...)

	server.AddTool(deleteContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Confirmation prompt - in MCP context, we'll add a warning in the description
		// but the actual confirmation would typically happen in the client UI

		// Call the service
		err = s.containerAppsService.DeleteContainerApp(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted Container App: %s", containerAppName)), nil
	})
}

// RegisterStartContainerAppTool registers the start container app tool with the MCP server
func (s *MCPServer) RegisterStartContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Start a Container App in Cloud.ru",
		"project_id",
		"containerapp_name",
	)
	startContainerAppTool := mcp.NewTool("cloudru_start_containerapp", toolOptions...)

	server.AddTool(startContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Call the service
		err = s.containerAppsService.StartContainerApp(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully started Container App: %s", containerAppName)), nil
	})
}

// RegisterStopContainerAppTool registers the stop container app tool with the MCP server
func (s *MCPServer) RegisterStopContainerAppTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Stop a Container App in Cloud.ru",
		"project_id",
		"containerapp_name",
	)
	stopContainerAppTool := mcp.NewTool("cloudru_stop_containerapp", toolOptions...)

	server.AddTool(stopContainerAppTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Call the service
		err = s.containerAppsService.StopContainerApp(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully stopped Container App: %s", containerAppName)), nil
	})
}

// RegisterGetContainerAppLogsTool registers the get container app logs tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppLogsTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get logs for a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"containerapp_name",
	)
	getContainerAppLogsTool := mcp.NewTool("cloudru_get_containerapp_logs", toolOptions...)

	server.AddTool(getContainerAppLogsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Call the service
		containerAppLogs, err := s.containerAppsService.GetContainerAppLogs(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(containerAppLogs, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}

// RegisterGetContainerAppSystemLogsTool registers the get container app system logs tool with the MCP server
func (s *MCPServer) RegisterGetContainerAppSystemLogsTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get system logs for a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
		"containerapp_name",
	)
	getContainerAppSystemLogsTool := mcp.NewTool("cloudru_get_containerapp_system_logs", toolOptions...)

	server.AddTool(getContainerAppSystemLogsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		// Call the service
		containerAppSystemLogs, err := s.containerAppsService.GetContainerAppSystemLogs(projectID, containerAppName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output with limit of 200 records
		limitedData, err := marshalJSONWithLimit(containerAppSystemLogs.Data, 200)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		// Create the response structure with limited data
		response := map[string]interface{}{
			"data": limitedData,
		}

		// Marshal the final response
		result, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}

// RegisterGetListDockerRegistriesTool registers the get list docker registries tool with the MCP server
func (s *MCPServer) RegisterGetListDockerRegistriesTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get list of Docker Registries from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru",
		"project_id",
	)
	getListDockerRegistriesTool := mcp.NewTool("cloudru_get_list_docker_registries", toolOptions...)

	server.AddTool(getListDockerRegistriesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

// RegisterCreateDockerRegistryTool registers the create docker registry tool with the MCP server
func (s *MCPServer) RegisterCreateDockerRegistryTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Create a new Docker Registry in Cloud.ru",
		"project_id",
		"registry_name",
		"registry_is_public",
	)
	createDockerRegistryTool := mcp.NewTool("cloudru_create_docker_registry", toolOptions...)

	server.AddTool(createDockerRegistryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

// RegisterGetRegistryImagesTool registers the get registry images tool with the MCP server
func (s *MCPServer) RegisterGetRegistryImagesTool(server *server.MCPServer) {
	// Prepare tool options including description and fields
	toolOptions := s.getMCPFieldsOptions(
		"Get list of images from a Docker registry in Cloud.ru",
		"registry_name",
	)
	getRegistryImagesTool := mcp.NewTool("cloudru_get_registry_images", toolOptions...)

	server.AddTool(getRegistryImagesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get registry name
		registryName, err := s.getMCPFieldValue("registry_name", request)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the service
		images, err := s.dockerService.GetRegistryImages(registryName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert to JSON for output
		result, err := json.MarshalIndent(images, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}
