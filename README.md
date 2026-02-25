# Cloud.ru Container Apps MCP

A Model Context Protocol (MCP) server for interacting with Cloud.ru Container Apps and Artifact Registry.

## What is MCP?

Model Context Protocol (MCP) is an open standard that enables seamless integration between AI assistants and external tools or data sources. It allows AI models to interact with your applications, services, and data in a secure and controlled manner.

With MCP, you can:
- Extend the capabilities of AI assistants beyond their training data
- Enable real-time interaction with live systems and APIs
- Provide contextually relevant information from your own data sources
- Execute complex workflows without leaving your AI assistant interface

### Example Usage

Instead of manually performing tasks, you can simply ask your AI assistant:

```
"Build image, push and deploy to Cloud.ru Container Apps via MCP"
```

Your AI assistant, using this MCP server, can then:
1. Build a Docker image of your application
2. Push it to Cloud.ru Artifact Registry
3. Update your Container App with the new image
4. Report back the status of the deployment

All of this happens automatically through natural language commands, making complex DevOps tasks accessible to everyone.


## Installation cloudru-containerapps-mcp to your system
```bash
go install github.com/Nick1994209/cloudru-containerapps-mcp/cmd/cloudru-containerapps-mcp@latest
```
[docs/INSTALLATION.md](docs/INSTALLATION.md)

## Add cloudru-containerapps-mcp to your IDE. For example VisualStudioCode or Cursor
Added MCP Setup
```json
{
  "mcpServers": {
    "cloudru-containerapps-mcp": {
      "command": "cloudru-containerapps-mcp",
      "args": [],
      "env": {
        "CLOUDRU_KEY_ID": "********",
        "CLOUDRU_KEY_SECRET": "********",
        "CLOUDRU_PROJECT_ID": "********",
      },
      "timeout": 900,
      "disabledTools": []
    }
  }
}
```
[docs/HOW_ADD_TO_IDE.md](docs/HOW_ADD_TO_IDE.md)

The server will listen for JSON-RPC messages on stdin/stdout.

## MCP Environment variables
[docs/ENVIRONMENT_VARIABLES.md](docs/ENVIRONMENT_VARIABLES.md)

### Functions

#### cloudru_containerapps_description()

Returns usage instructions for this MCP.

#### cloudru_docker_login(registry_name)

Logs into the Cloud.ru Docker registry using the provided credentials.

Parameters:
- `registry_name`: Name of the registry (falls back to CLOUDRU_REGISTRY_NAME env var)

If login fails, you'll need to:
1. Go to Cloud.ru Evolution Artifact Registry
2. Create a registry
3. Obtain access keys
4. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work

#### cloudru_docker_build_and_push(registry_name, repository_name, image_version, dockerfile_path, dockerfile_target, dockerfile_folder, show_commands)

Builds a Docker image and pushes it to Cloud.ru Artifact Registry.

Parameters:
- `registry_name`: Name of the registry (falls back to CLOUDRU_REGISTRY_NAME env var)
- `repository_name`: Name of the repository (falls back to CLOUDRU_REPOSITORY_NAME env var, then to current directory name)
- `image_version`: Version/tag for the image (optional, defaults to 'latest')
- `dockerfile_path`: Path to Dockerfile (optional, defaults to 'Dockerfile')
- `dockerfile_target`: Target stage in a multi-stage Dockerfile (optional, defaults to '-' which means no target)
- `dockerfile_folder`: Dockerfile folder (build context, defaults to '.' which means current directory)
- `show_commands`: If true, return Docker build and push commands without executing them (optional, defaults to 'true')

If Docker push fails due to authentication issues and CLOUDRU_KEY_ID/CLOUDRU_KEY_SECRET environment variables are set, the function will attempt to re-login and retry the push operation.

#### cloudru_get_list_containerapps(project_id)

Gets a list of Container Apps from Cloud.ru. Project ID can be set via CLOUDRU_PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)

#### cloudru_get_containerapp(project_id, containerapp_name)

Gets a specific Container App from Cloud.ru by name. Project ID can be set via CLOUDRU_PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to retrieve

#### cloudru_create_containerapp(project_id, containerapp_name, containerapp_port, containerapp_image, containerapp_auto_deployments_enabled, containerapp_auto_deployments_pattern, containerapp_privileged, containerapp_idle_timeout, containerapp_timeout, containerapp_cpu, containerapp_min_instance_count, containerapp_max_instance_count, containerapp_description, containerapp_publicly_accessible, containerapp_protocol, containerapp_environment_variables, containerapp_command, containerapp_args)

Creates a new Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to create
- `containerapp_port`: Port number for the Container App
- `containerapp_image`: Image for the Container App
- `containerapp_auto_deployments_enabled`: Enable auto deployments (optional, defaults to "false")
- `containerapp_auto_deployments_pattern`: Auto deployments pattern (optional, defaults to "latest")
- `containerapp_privileged`: Run container in privileged mode (optional, defaults to "false")
- `containerapp_idle_timeout`: Container idle timeout (optional, defaults to "600s")
- `containerapp_timeout`: Request timeout (optional, defaults to "60s")
- `containerapp_cpu`: CPU allocation (optional, defaults to "0.1", options: 0.1, 0.2, 0.5, 1)
- `containerapp_min_instance_count`: Minimum number of instances for scaling (optional, defaults to "0")
- `containerapp_max_instance_count`: Maximum number of instances for scaling (optional, defaults to "1")
- `containerapp_description`: Description of the container app (optional, defaults to empty string)
- `containerapp_publicly_accessible`: Whether the container app is publicly accessible (optional, defaults to "true")
- `containerapp_protocol`: Protocol for the container app (optional, defaults to "http_1", options: http_1, http_2)
- `containerapp_environment_variables`: Environment variables in format <name>='<value>';<next_name>='value2' (optional)
- `containerapp_command`: Command to run in the container (comma-separated values) (optional)
- `containerapp_args`: Arguments for the command (comma-separated values) (optional)

#### cloudru_delete_containerapp(project_id, containerapp_name)

Deletes a Container App from Cloud.ru. WARNING: This action cannot be undone!

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to delete

#### cloudru_start_containerapp(project_id, containerapp_name)

Starts a Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to start

#### cloudru_stop_containerapp(project_id, containerapp_name)

Stops a Container App in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to stop

#### cloudru_get_containerapp_logs(project_id, containerapp_name)

Gets logs for a specific Container App from Cloud.ru by name. Project ID can be set via CLOUDRU_PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to get logs from

#### cloudru_get_list_docker_registries(project_id)

Gets a list of Docker Registries from Cloud.ru. Project ID can be set via CLOUDRU_PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)

#### cloudru_create_docker_registry(project_id, registry_name, registry_is_public)

Creates a new Docker Registry in Cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to CLOUDRU_PROJECT_ID env var)
- `registry_name`: Name of the Docker Registry to create
- `registry_is_public`: Boolean flag indicating if the registry should be public (true) or private (false)

#### cloudru_get_registry_images(registry_name)

Gets a list of images from a Docker registry in Cloud.ru.

Parameters:
- `registry_name`: Name of the registry (falls back to CLOUDRU_REGISTRY_NAME env var)

Note: This function is currently disabled in the main.go file (line 56 is commented out).

#### cloudru_get_containerapp_system_logs(project_id, containerapp_name)

Gets system logs for a specific Container App from Cloud.ru by name. Project ID can be set via PROJECT_ID environment variable and obtained from console.cloud.ru.

Parameters:
- `project_id`: Project ID in Cloud.ru (falls back to PROJECT_ID env var)
- `containerapp_name`: Name of the Container App to get system logs from

Note: This function is currently disabled in the main.go file (line 53 is commented out).

## Currently Disabled Functions

The following functions are implemented but currently disabled in the main.go file:

1. `cloudru_get_containerapp_system_logs()` - Get system logs for a specific Container App (line 53 is commented out)
2. `cloudru_get_registry_images()` - Get list of images from a Docker registry (line 56 is commented out)

To enable these functions, uncomment the respective registration lines in [`cmd/cloudru-containerapps-mcp/main.go`](cmd/cloudru-containerapps-mcp/main.go).

## Development Guidelines

When contributing to this project, please follow our [Git Commit Guidelines](docs/GIT_COMMIT_GUIDELINES.md) for consistent commit messages in English using the Conventional Commits format.

## Documentation

For more information about Cloud.ru Container Apps, see:
https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work
