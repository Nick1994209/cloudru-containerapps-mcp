package domain

// DescriptionService provides usage instructions for the MCP
type DescriptionService interface {
	GetDescription() string
}

// DockerService handles Docker operations
type DockerService interface {
	Login(registryName string) (string, error)
	BuildAndPush(image DockerImage) (string, error)
	ShowBuildAndPushCommands(image DockerImage) (string, string, error)
	GetRegistryImages(registryName string) ([]RegistryImage, error)
}

// AuthService handles authentication operations
type AuthService interface {
	GetAccessToken() (string, error)
}

// ContainerAppsService handles Cloud.ru Container Apps API operations
type ContainerAppsService interface {
	GetListContainerApps(projectID string) ([]ContainerApp, error)
	GetContainerApp(projectID string, containerAppName string) (*ContainerApp, error)
	CreateContainerApp(request CreateContainerAppRequest) (*Operation, error)
	PatchContainerApp(projectID string, containerAppName string, request PatchContainerAppRequest) (*Operation, error)
	DeleteContainerApp(projectID string, containerAppName string) (*Operation, error)
	StartContainerApp(projectID string, containerAppName string) (*Operation, error)
	StopContainerApp(projectID string, containerAppName string) (*Operation, error)
	GetContainerAppLogs(projectID string, containerAppName string) (*ContainerAppLogs, error)
	GetContainerAppSystemLogs(projectID string, containerAppName string) (*ContainerAppSystemLogs, error)
}

// ArtifactRegistryService handles Cloud.ru Artifact Registry API operations
type ArtifactRegistryService interface {
	GetListDockerRegistries(projectID string) ([]DockerRegistry, error)
	CreateDockerRegistry(projectID string, registryName string, isPublic bool) (*DockerRegistry, error)
}

// JobsService handles Cloud.ru Jobs API operations
type JobsService interface {
	GetListJobs(projectID string, pageSize string) ([]Job, error)
	GetJob(projectID string, jobName string) (*Job, error)
	CreateJob(request CreateJobRequest) (*Operation, error)
	PatchJob(projectID string, jobName string, request PatchJobRequest) (*Operation, error)
	DeleteJob(projectID string, jobName string) (*Operation, error)
	ExecuteJob(projectID string, jobName string, params map[string]interface{}) (*JobExecution, error)
	GetListExecutions(projectID string, jobName string, pageSize string) ([]JobExecution, error)
}
