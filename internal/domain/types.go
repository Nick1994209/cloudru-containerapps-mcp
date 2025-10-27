package domain

// Credentials represents the authentication credentials for Cloud.ru
type Credentials struct {
	KeyID     string
	KeySecret string
}

// DockerImage represents a Docker image to be built and pushed
type DockerImage struct {
	RegistryName     string
	RepositoryName   string
	ImageVersion     string
	DockerfilePath   string
	DockerfileTarget string
	DockerfileFolder string
}

// CreateContainerAppRequest represents a request to create a Container App
type CreateContainerAppRequest struct {
	ProjectID              string `json:"projectId"`
	ContainerAppName       string `json:"containerAppName"`
	ContainerAppPort       int    `json:"containerAppPort"`
	ContainerAppImage      string `json:"containerAppImage"`
	AutoDeploymentsEnabled bool   `json:"autoDeploymentsEnabled"`
	AutoDeploymentsPattern string `json:"autoDeploymentsPattern"`
	Privileged             bool   `json:"privileged"`
	IdleTimeout            string `json:"idleTimeout"`
	Timeout                string `json:"timeout"`
	CPU                    string `json:"cpu"`
}

// ContainerApp represents a Cloud.ru Container App
type ContainerApp struct {
	ProjectID     string `json:"projectId"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	Configuration struct {
		Ingress struct {
			PubliclyAccessible bool   `json:"publiclyAccessible"`
			PublicUri          string `json:"publicUri"`
			//InternalUri            string        `json:"internalUri"`
			// AdditionalPortMappings []interface{} `json:"additionalPortMappings"`
		} `json:"ingress"`
		AutoDeployments struct {
			Enabled bool   `json:"enabled"`
			Pattern string `json:"pattern"`
		} `json:"autoDeployments"`
		Privileged bool `json:"privileged"`
	} `json:"configuration"`
	Template struct {
		Timeout     string `json:"timeout"`
		IdleTimeout string `json:"idleTimeout"`
		Protocol    string `json:"protocol"`
		Scaling     struct {
			MinInstanceCount int `json:"minInstanceCount"`
			MaxInstanceCount int `json:"maxInstanceCount"`
			Rule             struct {
				Type  string `json:"type"`
				Value struct {
					Soft int `json:"soft"`
					Hard int `json:"hard"`
				} `json:"value"`
			} `json:"rule"`
		} `json:"scaling"`
		Containers []struct {
			Name      string `json:"name"`
			Image     string `json:"image"`
			Resources struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"resources"`
			ContainerPort int `json:"containerPort"`
			Env           []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Type  string `json:"type,omitempty"`
			} `json:"env"`
			Command      []interface{} `json:"command"`
			Args         []interface{} `json:"args"`
			VolumeMounts []struct {
				Name      string `json:"name"`
				MountPath string `json:"mountPath"`
				ReadOnly  bool   `json:"readOnly"`
			} `json:"volumeMounts"`
		} `json:"containers"`
		InitContainers []interface{} `json:"initContainers"`
		Volumes        []struct {
			Name             string `json:"name"`
			Type             string `json:"type"`
			VolumeAttributes struct {
				BucketName string `json:"bucketName"`
				TenantId   string `json:"tenantId"`
				Region     string `json:"region"`
				ReadOnly   string `json:"readOnly"`
				Entrypoint string `json:"entrypoint"`
			} `json:"volumeAttributes"`
		} `json:"volumes"`
	} `json:"template"`
}

// DockerRegistry represents a Cloud.ru Docker Registry
type DockerRegistry struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	CreatedAt                string `json:"createdAt"`
	UpdatedAt                string `json:"updatedAt"`
	RegistryType             string `json:"registryType"`
	RetentionPolicyIsEnabled bool   `json:"retentionPolicyIsEnabled"`
	RetentionPolicy          string `json:"retentionPolicy"`
	Status                   string `json:"status"`
	IsPublic                 bool   `json:"isPublic"`
	QuarantineMode           string `json:"quarantineMode"`
}

// ContainerAppLogs represents the logs response from Cloud.ru Container App
type ContainerAppLogs struct {
	Data []ContainerAppLogEntry `json:"data"`
}

// ContainerAppLogEntry represents a single log entry
type ContainerAppLogEntry struct {
	Timestamp     string `json:"timestamp"`
	Message       string `json:"message"`
	VersionID     string `json:"versionId"`
	PodName       string `json:"podName"`
	Level         string `json:"level"`
	ContainerName string `json:"containerName"`
}

// ContainerAppSystemLogs represents the system logs response from Cloud.ru Container App
type ContainerAppSystemLogs struct {
	Data []ContainerAppSystemLogEntry `json:"data"`
}

// ContainerAppSystemLogEntry represents a single system log entry
type ContainerAppSystemLogEntry struct {
	EventType    string `json:"eventType"`
	Component    string `json:"component"`
	Reason       string `json:"reason"`
	Message      string `json:"message"`
	RevisionName string `json:"revisionName"`
}

// RegistryImage represents an image in the Docker registry
type RegistryImage struct {
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Digest    string `json:"digest"`
	CreatedAt string `json:"createdAt"`
	Size      int64  `json:"size"`
	MediaType string `json:"mediaType"`
}

// RegistryImagesResponse represents the response from registry images API
type RegistryImagesResponse struct {
	Images []RegistryImage `json:"images"`
}
