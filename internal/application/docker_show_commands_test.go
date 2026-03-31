package application

import (
	"testing"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDockerApplication_ShowBuildAndPushCommands(t *testing.T) {
	// Create a mock config with a test registry domain
	cfg := &config.Config{
		RegistryDomain: "cr.cloud.ru",
		KeyID:          "test-key-id",
		KeySecret:      "test-key-secret",
	}

	// Create a DockerApplication instance
	dockerApp := &DockerApplication{
		registryDomain: cfg.RegistryDomain,
		creds: domain.Credentials{
			KeyID:     cfg.KeyID,
			KeySecret: cfg.KeySecret,
		},
		// We won't actually use the auth service in this test
		authService: nil,
	}

	tests := []struct {
		name             string
		image            domain.DockerImage
		expectedBuildCmd string
		expectedPushCmd  string
	}{
		{
			name: "Basic case with all required fields",
			image: domain.DockerImage{
				RegistryName:   "test-registry",
				RepositoryName: "test-repo",
				ImageVersion:   "latest",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "With custom Dockerfile path",
			image: domain.DockerImage{
				RegistryName:   "test-registry",
				RepositoryName: "test-repo",
				ImageVersion:   "v1.0.0",
				DockerfilePath: "Dockerfile.prod",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:v1.0.0 -f Dockerfile.prod .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:v1.0.0",
		},
		{
			name: "With Dockerfile target",
			image: domain.DockerImage{
				RegistryName:     "test-registry",
				RepositoryName:   "test-repo",
				ImageVersion:     "latest",
				DockerfilePath:   "Dockerfile",
				DockerfileTarget: "production",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest --target production -f Dockerfile .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "With custom build context folder",
			image: domain.DockerImage{
				RegistryName:     "test-registry",
				RepositoryName:   "test-repo",
				ImageVersion:     "latest",
				DockerfilePath:   "Dockerfile",
				DockerfileFolder: "./build",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest -f Dockerfile ./build",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "With all optional fields",
			image: domain.DockerImage{
				RegistryName:     "test-registry",
				RepositoryName:   "test-repo",
				ImageVersion:     "v2.0.0",
				DockerfilePath:   "Dockerfile.custom",
				DockerfileTarget: "release",
				DockerfileFolder: "./deploy",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:v2.0.0 --target release -f Dockerfile.custom ./deploy",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:v2.0.0",
		},
		{
			name: "With empty Dockerfile path",
			image: domain.DockerImage{
				RegistryName:   "test-registry",
				RepositoryName: "test-repo",
				ImageVersion:   "latest",
				DockerfilePath: "",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "With empty target",
			image: domain.DockerImage{
				RegistryName:     "test-registry",
				RepositoryName:   "test-repo",
				ImageVersion:     "latest",
				DockerfilePath:   "Dockerfile",
				DockerfileTarget: "",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest -f Dockerfile .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "With dash target (should be treated as empty)",
			image: domain.DockerImage{
				RegistryName:     "test-registry",
				RepositoryName:   "test-repo",
				ImageVersion:     "latest",
				DockerfilePath:   "Dockerfile",
				DockerfileTarget: "-",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest -f Dockerfile .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "With dot folder (should be treated as current directory)",
			image: domain.DockerImage{
				RegistryName:     "test-registry",
				RepositoryName:   "test-repo",
				ImageVersion:     "latest",
				DockerfilePath:   "Dockerfile",
				DockerfileFolder: ".",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest -f Dockerfile .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "With empty folder (should be treated as current directory)",
			image: domain.DockerImage{
				RegistryName:     "test-registry",
				RepositoryName:   "test-repo",
				ImageVersion:     "latest",
				DockerfilePath:   "Dockerfile",
				DockerfileFolder: "",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest -f Dockerfile .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "Empty Dockerfile path with target",
			image: domain.DockerImage{
				RegistryName:     "test-registry",
				RepositoryName:   "test-repo",
				ImageVersion:     "latest",
				DockerfilePath:   "",
				DockerfileTarget: "production",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest --target production .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
		{
			name: "With empty ImageVersion (should default to latest)",
			image: domain.DockerImage{
				RegistryName:   "test-registry",
				RepositoryName: "test-repo",
				ImageVersion:   "",
			},
			expectedBuildCmd: "docker build --platform linux/amd64 -t test-registry.cr.cloud.ru/test-repo:latest .",
			expectedPushCmd:  "docker push --platform linux/amd64 test-registry.cr.cloud.ru/test-repo:latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the helper functions directly to avoid authentication issues
			buildCmd := dockerApp.generateBuildCommand(tt.image)
			pushCmd := dockerApp.generatePushCommand(tt.image)

			assert.Equal(t, tt.expectedBuildCmd, buildCmd)
			assert.Equal(t, tt.expectedPushCmd, pushCmd)
		})
	}
}
