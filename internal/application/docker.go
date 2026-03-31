package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/application/cloudru"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/config"
	"github.com/Nick1994209/cloudru-containerapps-mcp/internal/domain"
)

// DockerApplication implements the DockerService interface using actual Docker commands
type DockerApplication struct {
	registryDomain string
	creds          domain.Credentials
	authService    domain.AuthService
}

// NewDockerApplication creates a new DockerApplication with config
func NewDockerApplication(cfg *config.Config) domain.DockerService {
	return &DockerApplication{
		registryDomain: cfg.RegistryDomain,
		creds: domain.Credentials{
			KeyID:     cfg.KeyID,
			KeySecret: cfg.KeySecret,
		},
		authService: cloudru.NewAuthApplication(cfg),
	}
}

// Login logs into the Cloud.ru Docker registry using Docker CLI
func (d *DockerApplication) Login(registryName string) (string, error) {
	loginTarget := fmt.Sprintf("%s.%s", registryName, d.registryDomain)
	cmd := exec.Command("docker", "login", loginTarget, "-u", d.creds.KeyID, "--password-stdin")

	// Create a pipe to send the password to stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %w when trying to login", err)
	}

	// Execute the command
	go func() {
		defer stdin.Close()
		fmt.Fprint(stdin, d.creds.KeySecret)
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker login to %s failed: %w\nOutput: %s\n\nPlease ensure:\n1. The registry exists in Cloud.ru Evolution Artifact Registry\n2. You have created a registry and obtained access keys\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", loginTarget, err, string(output))
	}

	return loginTarget, nil
}

// BuildAndPush builds and pushes a Docker image to Cloud.ru Artifact Registry
func (d *DockerApplication) BuildAndPush(image domain.DockerImage) (string, error) {
	// Login to the Docker registry
	if _, err := d.Login(image.RegistryName); err != nil {
		return "", err
	}

	// Get the command strings
	buildCmdStr, pushCmdStr := d.generateCommands(image)

	// Extract image tag for return value and error messages
	imageTag := d.generateImageTag(image)

	// Build the Docker image
	// Split the build command string and execute it
	buildCmdParts := strings.Fields(buildCmdStr)
	if len(buildCmdParts) > 0 {
		buildCmd := exec.Command(buildCmdParts[0], buildCmdParts[1:]...)
		buildOutput, buildErr := buildCmd.CombinedOutput()

		// Always include build output in the response for visibility
		if len(buildOutput) > 0 {
			fmt.Printf("Docker build output:\n%s\n", string(buildOutput))
		}

		if buildErr != nil {
			return "", fmt.Errorf("failed to build Docker image %s: %w\nOutput: %s", imageTag, buildErr, string(buildOutput))
		}
	}

	// Push the Docker image
	// Split the push command string and execute it
	pushCmdParts := strings.Fields(pushCmdStr)
	if len(pushCmdParts) > 0 {
		pushCmd := exec.Command(pushCmdParts[0], pushCmdParts[1:]...)
		pushOutput, pushErr := pushCmd.CombinedOutput()

		// Always include push output in the response for visibility
		if len(pushOutput) > 0 {
			fmt.Printf("Docker push output:\n%s\n", string(pushOutput))
		}

		if pushErr != nil {
			return "", fmt.Errorf("docker push failed: %w\nOutput: %s\n\nTo resolve this issue:\n1. Ensure you are logged in to the Docker registry\n2. Run the cloudru_docker_login function\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", pushErr, string(pushOutput))
		}
	}

	return imageTag, nil
}

// generateImageTag creates the full image tag for a Docker image
// If ImageVersion is empty, it defaults to "latest"
func (d *DockerApplication) generateImageTag(image domain.DockerImage) string {
	imageVersion := image.ImageVersion
	if imageVersion == "" {
		imageVersion = "latest"
	}
	return fmt.Sprintf("%s.%s/%s:%s", image.RegistryName, d.registryDomain, image.RepositoryName, imageVersion)
}

// generateBuildCommand creates the docker build command string
func (d *DockerApplication) generateBuildCommand(image domain.DockerImage) string {
	imageTag := d.generateImageTag(image)

	// Start with the base command
	buildCommand := fmt.Sprintf("docker build --platform linux/amd64 -t %s", imageTag)

	// Add target if specified
	if image.DockerfileTarget != "" && image.DockerfileTarget != "-" {
		buildCommand = fmt.Sprintf("%s --target %s", buildCommand, image.DockerfileTarget)
	}

	// Handle Dockerfile path - if empty, don't include the -f flag
	if image.DockerfilePath != "" {
		buildCommand = fmt.Sprintf("%s -f %s", buildCommand, image.DockerfilePath)
	}

	// Set build context folder, default to current directory if not specified
	buildContext := "."
	if image.DockerfileFolder != "" && image.DockerfileFolder != "." {
		buildContext = image.DockerfileFolder
	}
	buildCommand = fmt.Sprintf("%s %s", buildCommand, buildContext)

	return buildCommand
}

// generatePushCommand creates the docker push command string
func (d *DockerApplication) generatePushCommand(image domain.DockerImage) string {
	imageTag := d.generateImageTag(image)
	return fmt.Sprintf("docker push --platform linux/amd64 %s", imageTag)
}

// generateCommands returns the docker build and push commands as strings
func (d *DockerApplication) generateCommands(image domain.DockerImage) (string, string) {
	buildCmd := d.generateBuildCommand(image)
	pushCmd := d.generatePushCommand(image)
	return buildCmd, pushCmd
}

// ShowBuildAndPushCommands returns the docker build and push commands as strings without executing them
func (d *DockerApplication) ShowBuildAndPushCommands(image domain.DockerImage) (string, string, error) {
	if _, err := d.Login(image.RegistryName); err != nil {
		return "", "", err
	}

	buildCmd, pushCmd := d.generateCommands(image)
	return buildCmd, pushCmd, nil
}

// GetRegistryImages gets a list of images from a Docker registry using Cloud.ru API token
func (d *DockerApplication) GetRegistryImages(registryName string) ([]domain.RegistryImage, error) {
	// Get access token using AuthApplication
	token, err := d.authService.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// now this is not working =(, we will fix later
	// Use the Cloud.ru API token to access the registry
	registryURL := fmt.Sprintf("https://%s.%s/v2/_catalog", registryName, d.registryDomain)

	client := &http.Client{}
	req, err := http.NewRequest("GET", registryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Use Bearer token for authentication
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "cloudru-containerapps-mcp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to registry API: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the catalog response
	var catalogResponse struct {
		Repositories []string `json:"repositories"`
	}

	if err := json.Unmarshal(body, &catalogResponse); err != nil {
		return nil, fmt.Errorf("failed to parse catalog response: %w", err)
	}

	// Convert to RegistryImage format
	var images []domain.RegistryImage
	for _, repo := range catalogResponse.Repositories {
		image := domain.RegistryImage{
			Name:      repo,
			Tag:       "latest",
			Digest:    "",
			CreatedAt: "",
			Size:      0,
			MediaType: "application/vnd.docker.distribution.manifest.v2+json",
		}
		images = append(images, image)
	}

	return images, nil
}
