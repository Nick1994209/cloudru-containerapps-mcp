package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

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
	if _, err := d.Login(image.RegistryName); err != nil {
		return "", err
	}

	imageTag := fmt.Sprintf("%s.%s/%s:%s", image.RegistryName, d.registryDomain, image.RepositoryName, image.ImageVersion)

	// Build the Docker image
	var buildCmd *exec.Cmd
	// Set build context folder, default to current directory if not specified
	buildContext := "."
	if image.DockerfileFolder != "" && image.DockerfileFolder != "." {
		buildContext = image.DockerfileFolder
	}

	if image.DockerfileTarget != "" && image.DockerfileTarget != "-" {
		buildCmd = exec.Command("docker", "build", "--platform", "linux/amd64", "-t", imageTag, "--target", image.DockerfileTarget, "-f", image.DockerfilePath, buildContext)
	} else {
		buildCmd = exec.Command("docker", "build", "--platform", "linux/amd64", "-t", imageTag, "-f", image.DockerfilePath, buildContext)
	}
	buildOutput, buildErr := buildCmd.CombinedOutput()

	// Always include build output in the response for visibility
	if len(buildOutput) > 0 {
		fmt.Printf("Docker build output:\n%s\n", string(buildOutput))
	}

	if buildErr != nil {
		return "", fmt.Errorf("failed to build Docker image %s: %w\nOutput: %s", imageTag, buildErr, string(buildOutput))
	}

	// Push the Docker image
	pushCmd := exec.Command("docker", "push", "--platform", "linux/amd64", imageTag)
	pushOutput, pushErr := pushCmd.CombinedOutput()

	// Always include push output in the response for visibility
	if len(pushOutput) > 0 {
		fmt.Printf("Docker push output:\n%s\n", string(pushOutput))
	}

	if pushErr != nil {
		return "", fmt.Errorf("docker push failed: %w\nOutput: %s\n\nTo resolve this issue:\n1. Ensure you are logged in to the Docker registry\n2. Run the cloudru_docker_login function\n3. See documentation: https://cloud.ru/docs/container-apps-evolution/ug/topics/tutorials__before-work", pushErr, string(pushOutput))
	}

	return imageTag, nil
}

// ShowBuildAndPushCommands returns the docker build and push commands as strings without executing them
func (d *DockerApplication) ShowBuildAndPushCommands(image domain.DockerImage) (string, string, error) {
	if _, err := d.Login(image.RegistryName); err != nil {
		return "", "", err
	}

	imageTag := fmt.Sprintf("%s.%s/%s:%s", image.RegistryName, d.registryDomain, image.RepositoryName, image.ImageVersion)

	// Build the Docker build command
	var buildCmd string
	// Set build context folder, default to current directory if not specified
	buildContext := "."
	if image.DockerfileFolder != "" && image.DockerfileFolder != "." {
		buildContext = image.DockerfileFolder
	}

	if image.DockerfileTarget != "" && image.DockerfileTarget != "-" {
		buildCmd = fmt.Sprintf("docker build --platform linux/amd64 -t %s --target %s -f %s %s", imageTag, image.DockerfileTarget, image.DockerfilePath, buildContext)
	} else {
		buildCmd = fmt.Sprintf("docker build --platform linux/amd64 -t %s -f %s %s", imageTag, image.DockerfilePath, buildContext)
	}

	// Build the Docker push command
	pushCmd := fmt.Sprintf("docker push --platform linux/amd64 %s", imageTag)

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
