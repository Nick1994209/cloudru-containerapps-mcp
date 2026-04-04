package utils

import "strings"

// ParseEnvironmentVariables parses environment variables from format <name>='<value>';<next_name>='value2'
func ParseEnvironmentVariables(environmentVariables string) []map[string]interface{} {
	var envVars []map[string]interface{}
	if environmentVariables != "" {
		// Split by semicolon
		variables := strings.Split(environmentVariables, ";")
		for _, variable := range variables {
			// Split by first equals sign
			parts := strings.SplitN(variable, "=", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
					value = value[1 : len(value)-1]
				}
				envVars = append(envVars, map[string]interface{}{
					"name":  name,
					"value": value,
				})
			}
		}
	}
	return envVars
}

// ParseCPU maps CPU allocation to memory allocation
func ParseCPU(cpu string) (string, string) {
	var memory string
	switch cpu {
	case "0.1":
		memory = "256Mi"
	case "0.2":
		memory = "512Mi"
	case "0.3":
		memory = "768Mi"
	case "0.5":
		memory = "1024Mi"
	case "1":
		memory = "4096Mi"
	default:
		// Default to 0.1 CPU and 256Mi memory for unknown values
		cpu = "0.1"
		memory = "256Mi"
	}

	return cpu, memory
}
