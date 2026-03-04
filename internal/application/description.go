package application

import (
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// DescriptionApplication implements the DescriptionService interface
type DescriptionApplication struct{}

// NewDescriptionApplication creates a new DescriptionApplication
func NewDescriptionApplication() domain.DescriptionService {
	return &DescriptionApplication{}
}

// GetDescription returns usage instructions for this MCP
func (d *DescriptionApplication) GetDescription() string {
	cfg := config.LoadConfig()

	return `Cloud.ru Container Apps MCP provides functions to interact with Cloud.ru Artifact Registry:

1. cloudru_containerapps_description() - Returns usage instructions for this MCP
2. cloudru_docker_login(registry_name) - Login to Cloud.ru Artifact registry (Docker registry)
3. cloudru_docker_build_and_push(registry_name, repository_name, image_version, dockerfile_path, dockerfile_target, dockerfile_folder, show_commands) - Build and push Docker image to Cloud.ru Artifact Registry (Docker registry)
4. cloudru_get_list_containerapps(project_id) - Get list of Container Apps from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
5. cloudru_get_containerapp(project_id, containerapp_name) - Get a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
6. cloudru_create_containerapp(project_id, containerapp_name, containerapp_port, containerapp_image, containerapp_auto_deployments_enabled, containerapp_auto_deployments_pattern, containerapp_privileged, containerapp_idle_timeout, containerapp_timeout, containerapp_cpu) - Create a new Container App in Cloud.ru
7. cloudru_patch_containerapp(project_id, containerapp_name, containerapp_port, containerapp_image, containerapp_auto_deployments_enabled, containerapp_auto_deployments_pattern, containerapp_privileged, containerapp_idle_timeout, containerapp_timeout, containerapp_cpu, containerapp_min_instance_count, containerapp_max_instance_count, containerapp_description, containerapp_publicly_accessible, containerapp_protocol, containerapp_environment_variables, containerapp_command, containerapp_args) - Patch an existing Container App in Cloud.ru. This function gets the current state, merges it with the new values, and updates the container app. Note: The privileged field is read-only and cannot be modified through this function.
8. cloudru_delete_containerapp(project_id, containerapp_name) - Delete a Container App from Cloud.ru. WARNING: This action cannot be undone!
9. cloudru_start_containerapp(project_id, containerapp_name) - Start a Container App in Cloud.ru
10. cloudru_stop_containerapp(project_id, containerapp_name) - Stop a Container App in Cloud.ru
11. cloudru_get_containerapp_logs(project_id, containerapp_name) - Get logs for a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
12. cloudru_get_list_docker_registries(project_id) - Get list of Docker Registries from Cloud.ru. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru
13. cloudru_create_docker_registry(project_id, registry_name, registry_is_public) - Create a new Docker Registry in Cloud.ru

Environment variables can be used as fallbacks for parameters:

**Current configuration values:**
- Current directory: (` + cfg.CurrentDir + `) (Name of the current working directory)
- CLOUDRU_KEY_ID: (` + maskSensitiveInfo(cfg.KeyID) + `) [required] (Authentication key identifier)
- CLOUDRU_KEY_SECRET: (` + maskSensitiveInfo(cfg.KeySecret) + `) [required] (Authentication key secret)
- CLOUDRU_PROJECT_ID: (` + cfg.ProjectID + `) (Project ID for Container Apps)
- CLOUDRU_REGISTRY_NAME: (` + cfg.RegistryName + `) (Registry for storing Docker images)
- CLOUDRU_REPOSITORY_NAME: (` + cfg.RepositoryName + `) (Name of the repository in the registry)
- CLOUDRU_DOCKERFILE: (` + cfg.Dockerfile + `) (Path to the Dockerfile to build the image, by default Dockerfile)
- CLOUDRU_DOCKERFILE_TARGET: (` + cfg.DockerfileTarget + `) (Dockerfile target stage, defaults to "-" which means no target)
- CLOUDRU_DOCKERFILE_FOLDER: (` + cfg.DockerfileFolder + `) (Dockerfile folder (build context), defaults to "." which means current directory)
- CLOUDRU_CONTAINERAPP_NAME: (` + cfg.ContainerAppName + `) (Container App name)

For more details see: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work`
}

// maskSensitiveInfo replaces the middle of a string with asterisks for sensitive data
func maskSensitiveInfo(value string) string {
	if len(value) == 0 {
		return ""
	}

	if len(value) <= 4 {
		return "***"
	}

	// Show first 2 and last 2 characters, replace the rest with asterisks
	start := value[:3]
	end := value[len(value)-3:]
	middle := "***"

	return start + middle + end
}
