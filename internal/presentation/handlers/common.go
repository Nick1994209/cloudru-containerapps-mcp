package handlers

import (
	"encoding/json"
	"fmt"
	"reflect"

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
				required:     true,
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
				description:  "CPU allocation (0.1 CPU - 256 Mi RAM, 0.2 CPU - 512 Mi RAM, 0.3 CPU - 768 RAM, 0.5 CPU - 1Gb RAM ...)",
				defaultValue: "0.1",
				required:     false,
				title:        "Options: 0.1, 0.2, 0.3, 0.5, 1",
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
				defaultValue: "This ContainerApps created via MCP",
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
		description := fieldData.description
		if fieldData.envValue != "" {
			description = fmt.Sprintf("%s (default: %s)", fieldData.description, fieldData.envValue)
		} else if fieldData.defaultValue != "" {
			description = fmt.Sprintf("%s (default: %s)", fieldData.description, fieldData.defaultValue)
		}
		opts := []mcp.PropertyOption{
			mcp.Description(description),
		}
		if fieldData.required && fieldData.envValue == "" {
			opts = append(opts, mcp.Required())
		}
		if fieldData.title != "" {
			opts = append(opts, mcp.Title(fieldData.title))
		}
		if fieldData.envValue != "" {
			opts = append(opts, mcp.DefaultString(fieldData.envValue))
		} else if fieldData.defaultValue != "" {
			opts = append(opts, mcp.DefaultString(fieldData.defaultValue))
		}
		result = append(result, mcp.WithString(field, opts...))
	}
	return result
}

func (s *MCPServer) getMCPFieldValue(field string, request mcp.CallToolRequest) (string, error) {
	fieldData := s.mappedFields[field]

	// Try to get the value from the request
	result, err := request.RequireString(field)
	if err != nil && fieldData.defaultValue == "" && fieldData.envValue == "" && fieldData.required {
		// If there's an error and no default or env value, and the field is required, return the error
		return "", err
	}

	// If we got a value from the request, use it
	if result != "" {
		return result, nil
	}

	// If we have an environment variable value, use it
	if fieldData.envValue != "" {
		return fieldData.envValue, nil
	}

	// Otherwise, use the default value if available
	if fieldData.defaultValue != "" {
		return fieldData.defaultValue, nil
	}

	// If the field is not required, return empty string
	if !fieldData.required {
		return "", nil
	}

	// If we get here and the field is required but has no value, return an error
	if fieldData.required && result == "" && fieldData.defaultValue == "" && fieldData.envValue == "" {
		return "", fmt.Errorf("field %s is empty: %s", field, fieldData.description)
	}

	// If we get here, return whatever we have (likely empty)
	return result, nil
}

func (s *MCPServer) getMCPBooleanFieldValue(field string, request mcp.CallToolRequest) (bool, error) {
	fieldValueStr, err := s.getMCPFieldValue(field, request)
	if err != nil {
		return false, err
	}

	if fieldValueStr == "true" || fieldValueStr == "1" {
		return true, nil
	} else if fieldValueStr == "false" || fieldValueStr == "0" {
		return false, nil
	} else {
		return false, fmt.Errorf("field %s must be 'true', 'false', '1', or '0', got: %s", field, fieldValueStr)
	}
}

// checkRequestHasKey checks if a request has a specific key in its arguments
func checkRequestHasKey(r mcp.CallToolRequest, key string) bool {
	args := r.GetArguments()
	_, ok := args[key]
	return ok
}

// RegisterAllTools registers all tools with the MCP server
func (s *MCPServer) RegisterAllTools(mcpServer *server.MCPServer) {
	s.RegisterDescriptionTool(mcpServer)
	s.RegisterDockerLoginTool(mcpServer)
	s.RegisterDockerBuildAndPushTool(mcpServer)
	s.RegisterGetListContainerAppsTool(mcpServer)
	s.RegisterGetContainerAppTool(mcpServer)
	s.RegisterPatchContainerAppTool(mcpServer)
	s.RegisterCreateContainerAppTool(mcpServer)
	s.RegisterDeleteContainerAppTool(mcpServer)
	s.RegisterStartContainerAppTool(mcpServer)
	s.RegisterStopContainerAppTool(mcpServer)
	s.RegisterGetContainerAppLogsTool(mcpServer)
	s.RegisterGetContainerAppSystemLogsTool(mcpServer)
	s.RegisterGetListDockerRegistriesTool(mcpServer)
	s.RegisterCreateDockerRegistryTool(mcpServer)
	s.RegisterGetRegistryImagesTool(mcpServer)
}
