package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// APIURLs holds the API endpoints for Cloud.ru services
type APIURLs struct {
	ContainersAPI string
	IAMAPI        string
	ArtifactAPI   string
}

// Config holds the configuration for the MCP
type Config struct {
	RegistryName     string
	RegistryDomain   string
	KeyID            string
	KeySecret        string
	RepositoryName   string
	Dockerfile       string
	DockerfileTarget string
	DockerfileFolder string
	ProjectID        string
	ContainerAppName string
	CurrentDir       string
	API              APIURLs
}

// EnvVarNames contains the names of environment variables
const (
	EnvRegistryName     = "CLOUDRU_REGISTRY_NAME"
	EnvRegistryDomain   = "CLOUDRU_REGISTRY_DOMAIN"
	EnvKeyID            = "CLOUDRU_KEY_ID"
	EnvKeySecret        = "CLOUDRU_KEY_SECRET"
	EnvRepositoryName   = "CLOUDRU_REPOSITORY_NAME"
	EnvProjectID        = "CLOUDRU_PROJECT_ID"
	EnvContainerAppName = "CLOUDRU_CONTAINERAPP_NAME"
	Dockerfile          = "CLOUDRU_DOCKERFILE"
	DockerfileTarget    = "CLOUDRU_DOCKERFILE_TARGET"
	DockerfileFolder    = "CLOUDRU_DOCKERFILE_FOLDER"
	EnvContainersAPI    = "CLOUDRU_CONTAINERS_API"
	EnvIAMAPI           = "CLOUDRU_IAM_API"
	EnvArtifactAPI      = "CLOUDRU_ARTIFACT_API"
)

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Overload()
	if err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	// Check for required environment variables
	keyID := os.Getenv(EnvKeyID)
	keySecret := os.Getenv(EnvKeySecret)

	if keyID == "" || keySecret == "" {
		log.Fatal(`CLOUDRU_KEY_ID and CLOUDRU_KEY_SECRET environment variables must be set.
		
To obtain access keys for authentication, please follow the instructions at:
https://cloud.ru/docs/console_api/ug/topics/quickstart

You will need a Key ID and Key Secret to use this service.`)
	}

	dir, err := os.Getwd()
	if err != nil {
		dir = "default"
	}
	projectDirName := filepath.Base(dir)

	// Set default API URLs if environment variables are not provided
	containersAPI := os.Getenv(EnvContainersAPI)
	if containersAPI == "" {
		containersAPI = "https://containers.api.cloud.ru"
	}

	iamAPI := os.Getenv(EnvIAMAPI)
	if iamAPI == "" {
		iamAPI = "https://iam.api.cloud.ru"
	}

	artifactAPI := os.Getenv(EnvArtifactAPI)
	if artifactAPI == "" {
		artifactAPI = "https://ar.api.cloud.ru"
	}

	// Set default registry domain if environment variable is not provided
	registryDomain := os.Getenv(EnvRegistryDomain)
	if registryDomain == "" {
		registryDomain = "cr.cloud.ru"
	}

	return &Config{
		RegistryName:     os.Getenv(EnvRegistryName),
		RegistryDomain:   registryDomain,
		KeyID:            keyID,
		KeySecret:        keySecret,
		RepositoryName:   os.Getenv(EnvRepositoryName),
		ProjectID:        os.Getenv(EnvProjectID),
		ContainerAppName: os.Getenv(EnvContainerAppName),
		Dockerfile:       os.Getenv(Dockerfile),
		DockerfileTarget: os.Getenv(DockerfileTarget),
		DockerfileFolder: os.Getenv(DockerfileFolder),
		CurrentDir:       projectDirName,
		API: APIURLs{
			ContainersAPI: containersAPI,
			IAMAPI:        iamAPI,
			ArtifactAPI:   artifactAPI,
		},
	}
}
