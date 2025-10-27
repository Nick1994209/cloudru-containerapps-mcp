## Example Prompts

### Basic Functions

- "Use cloudru_containerapps_description to tell me about this tool"
- "Run cloudru_docker_login with my registry credentials"
- "Execute cloudru_docker_build_and_push to deploy my application with version v1.2.3"

### Container Apps Management

#### List and Get Container Apps
- "Get list of Container Apps using cloudru_get_list_containerapps"
- "Retrieve details of my specific Container App named 'my-app' with cloudru_get_containerapp"

#### Create, Start, Stop, and Delete Container Apps
- "Create a new Container App called 'my-new-app' using cloudru_create_containerapp on port 8080 with image 'nginx'"
- "Start my Container App 'my-app' with cloudru_start_containerapp"
- "Stop my Container App 'my-app' with cloudru_stop_containerapp"
- "Delete my Container App 'my-old-app' with cloudru_delete_containerapp - be careful as this cannot be undone"

### Docker Registry Management

#### List and Create Docker Registries
- "List all Docker Registries in my project using cloudru_get_list_docker_registries"
- "Create a new private Docker Registry named 'my-private-registry' with cloudru_create_docker_registry"
- "Create a new public Docker Registry named 'my-public-registry' with cloudru_create_docker_registry"
