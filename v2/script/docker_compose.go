package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DockerComposeConfig represents the structure of a docker-compose.yml file
type DockerComposeConfig struct {
	Version  string                          `yaml:"version,omitempty"`
	Services map[string]DockerComposeService `yaml:"services"`
	Networks map[string]any                  `yaml:"networks,omitempty"`
	Volumes  map[string]any                  `yaml:"volumes,omitempty"`
}

// DockerComposeService represents a service in the docker-compose.yml file
type DockerComposeService struct {
	Image         string         `yaml:"image,omitempty"`
	ContainerName string         `yaml:"container_name,omitempty"`
	Build         any            `yaml:"build,omitempty"`
	EnvFile       any            `yaml:"env_file,omitempty"`
	Environment   any            `yaml:"environment,omitempty"`
	Ports         any            `yaml:"ports,omitempty"`
	Labels        any            `yaml:"labels,omitempty"`
	Volumes       any            `yaml:"volumes,omitempty"`
	Networks      any            `yaml:"networks,omitempty"`
	DependsOn     any            `yaml:"depends_on,omitempty"`
	Command       any            `yaml:"command,omitempty"`
	WorkingDir    string         `yaml:"working_dir,omitempty"`
	Restart       string         `yaml:"restart,omitempty"`
	Expose        any            `yaml:"expose,omitempty"`
	ExtraFields   map[string]any `yaml:",inline"`
}

// UpdateDockerCompose updates a docker-compose.yml file with environment variables from a consolidated .env file
func UpdateDockerCompose(consolidatedEnvFile string, outputFile string, forceOverwrite bool, customTemplate string, discoverDir string, configFile string) {
	// Get script directory
	scriptDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Set project root - find the actual project root instead of just assuming it's one level up
	projectRoot := scriptDir

	serviceDir := discoverDir

	// Use serviceDir if provided, otherwise use project root
	if serviceDir != "" {
		// If service dir path is relative, make it absolute from script directory
		if !filepath.IsAbs(serviceDir) {
			if strings.HasPrefix(serviceDir, "./") {
				discoverDir = filepath.Join(scriptDir, serviceDir[2:])
			} else if strings.HasPrefix(serviceDir, "../") {
				discoverDir = filepath.Join(scriptDir, serviceDir)
			} else {
				discoverDir = filepath.Join(scriptDir, serviceDir)
			}
		} else {
			discoverDir = serviceDir
		}
		fmt.Printf("Using specified service directory: %s\n", discoverDir)
	}

	fmt.Printf("Script directory: %s\n", scriptDir)
	fmt.Printf("Determined project root: %s\n", projectRoot)

	// Determine template file to use
	var templateFile string

	if customTemplate != "" {
		// Use the specified custom template
		templateFile = resolveFilePath(customTemplate, scriptDir, projectRoot)
		fmt.Printf("Using specified template file: %s\n", templateFile)
	} else {
		// Check if default template file exists, if not try to use regular docker-compose.yml
		templateFile = filepath.Join(projectRoot, "docker-compose.template.yml")

		if _, err := os.Stat(templateFile); err == nil {
			fmt.Printf("Using template file: %s\n", templateFile)
		} else {
			fmt.Printf("Template file is not specified")
			return
		}
	}

	// Resolve absolute paths
	consolidatedEnvFile = resolveFilePath(consolidatedEnvFile, scriptDir, projectRoot)
	if outputFile == "" {
		outputFile = filepath.Join(projectRoot, "docker-compose.yml")
	} else {
		outputFile = resolveFilePath(outputFile, scriptDir, projectRoot)
	}

	// Get the relative path to .env file from the project root
	relEnvPath, err := filepath.Rel(projectRoot, consolidatedEnvFile)
	if err != nil {
		relEnvPath = consolidatedEnvFile // Use absolute path if unable to get relative path
	}

	// Debug info
	fmt.Printf("Project root: %s\n", projectRoot)
	fmt.Printf("Template file: %s\n", templateFile)
	fmt.Printf("Consolidated env file: %s\n", consolidatedEnvFile)
	fmt.Printf("Output file: %s\n", outputFile)

	// Check if template file exists
	if _, err := os.Stat(templateFile); err != nil {
		fmt.Printf("Template file not found: %s\n", templateFile)
		return
	}

	// Check if consolidated env file exists
	if _, err := os.Stat(consolidatedEnvFile); err != nil {
		fmt.Printf("Consolidated env file not found: %s\n", consolidatedEnvFile)
		return
	}

	// Check if output file exists
	if _, err := os.Stat(outputFile); err == nil && !forceOverwrite {
		fmt.Printf("Output file %s already exists. Overwrite? (y/n): ", outputFile)
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(answer) != "y" {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	// Read template file
	templateBytes, err := os.ReadFile(templateFile)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		return
	}

	// Read consolidated env file
	consolidatedEnvBytes, err := os.ReadFile(consolidatedEnvFile)
	if err != nil {
		fmt.Printf("Error reading consolidated env file: %v\n", err)
		return
	}

	// Parse template file
	var dockerCompose DockerComposeConfig
	if err := yaml.Unmarshal(templateBytes, &dockerCompose); err != nil {
		fmt.Printf("Error parsing template file: %v\n", err)
		return
	}

	// Get environment variables from consolidated env file
	envVars := make(map[string]string)
	parseEnvFile(string(consolidatedEnvBytes), &envVars)

	// Check for services that are referenced in depends_on but aren't defined
	// First, gather all defined service names
	definedServices := make(map[string]bool)
	for serviceName := range dockerCompose.Services {
		definedServices[serviceName] = true
	}

	// Process Docker compose services
	for serviceName, service := range dockerCompose.Services {
		fmt.Printf("Processing service: %s\n", serviceName)

		serviceEnvFile := getServiceEnvFile(serviceName, discoverDir, configFile)

		// Check if service env file exists
		if _, err := os.Stat(serviceEnvFile); err != nil {
			fmt.Printf("Service env file not found: %s\n", serviceEnvFile)
			continue
		}
		serviceEnvBytes, err := os.ReadFile(serviceEnvFile)
		if err != nil {
			fmt.Printf("Error reading service env file: %v\n", err)
			continue
		}
		serviceEnvVars := make(map[string]string)
		parseEnvFile(string(serviceEnvBytes), &serviceEnvVars)

		// Update environment variables in service
		updateServiceEnvironment(serviceName, &service, &envVars, &serviceEnvVars, configFile)

		// Update ports in service
		updateServicePorts(serviceName, &service, &envVars, configFile)

		// Update service in Docker compose
		dockerCompose.Services[serviceName] = service
	}

	// Write updated Docker compose file
	updatedDockerComposeBytes, err := yaml.Marshal(dockerCompose)
	if err != nil {
		fmt.Printf("Error marshalling updated Docker compose: %v\n", err)
		return
	}

	if err := os.WriteFile(outputFile, updatedDockerComposeBytes, 0644); err != nil {
		fmt.Printf("Error writing updated Docker compose file: %v\n", err)
		return
	}

	fmt.Printf("Generated Docker compose file written to %s\n", outputFile)
	fmt.Printf("Services will use environment variables from %s\n", relEnvPath)
}

func updateServiceEnvironment(serviceName string, service *DockerComposeService, envVars *map[string]string, serviceEnvVars *map[string]string, configFile string) {
	servicePrefix := getServicePrefix(serviceName, configFile)
	fmt.Printf("  Looking for environment variables with prefix %s\n", servicePrefix)
	serviceEnvVarsKeys := []string{}

	for key := range *serviceEnvVars {
		serviceEnvVarsKeys = append(serviceEnvVarsKeys, key)
	}

	consolidatedEnvVarsKeys := []string{}

	for key := range *envVars {
		if strings.HasPrefix(key, servicePrefix) {
			consolidatedEnvVarsKeys = append(consolidatedEnvVarsKeys, key)
		}
	}

	// Create a list to store service environment variables
	envList := []string{}

	// Map service env vars to their corresponding consolidated env vars
	for _, key := range serviceEnvVarsKeys {
		// Find matching consolidated env var by stripping the service prefix
		// and checking if the remaining part matches the service env var
		found := false
		for _, consolidatedKey := range consolidatedEnvVarsKeys {
			varNameWithoutPrefix := strings.TrimPrefix(consolidatedKey, servicePrefix)
			if varNameWithoutPrefix == key {
				// Add the mapping using the original service var name and the consolidated var reference
				envList = append(envList, fmt.Sprintf("%s=${%s}", key, consolidatedKey))
				found = true
				break
			}
		}

		if !found {
			// If no specific match found, still keep the original environment variable
			envList = append(envList, fmt.Sprintf("%s=${%s}", key, key))
		}
	}

	// Set the updated environment list
	if len(envList) > 0 {
		service.Environment = envList
		fmt.Printf("  Updated environment variables for service %s with references to prefixed variables\n", serviceName)
	} else {
		fmt.Printf("  No matching environment variables found for service %s\n", serviceName)
	}
}

func updateServicePorts(serviceName string, service *DockerComposeService, envVars *map[string]string, configFile string) {

	config := getConfig(configFile)
	// Get the prefix from services-config.yaml
	servicePrefix := getServicePrefix(serviceName, configFile)
	fmt.Printf("  Looking for port variables for service %s\n", serviceName)

	// Create a list to store port mappings
	portMappings := []string{}
	envVarsKeys := []string{}

	for key := range *envVars {
		if strings.HasPrefix(key, servicePrefix) {
			envVarsKeys = append(envVarsKeys, key)
		}
	}
	var externalServices []string
	for _, service := range config.CommonServices {
		externalServices = append(externalServices, service.Name)
	}
	externalServices = append(externalServices, "db", "mongo", "elastic", "mail")

	// Look for port variables using regex-like matching
	for _, envKey := range envVarsKeys {
		// Check if the variable has the service prefix and contains _PORT
		if strings.Contains(envKey, "_PORT") {
			// Skip ports that belong to external services
			lowerKey := strings.ToLower(envKey)

			isExternalService := false
			for _, service := range externalServices {
				if strings.Contains(lowerKey, service) && serviceName != service {
					isExternalService = true
					break
				}
			}

			if isExternalService {
				fmt.Printf("  Skipping external service port: %s\n", envKey)
				continue
			}

			// Create port mapping - use external port if available, otherwise use the same port
			portMapping := fmt.Sprintf("${%s}:${%s}", envKey, envKey)

			portMappings = append(portMappings, portMapping)
			fmt.Printf("  Added port mapping %s for variable %s\n", portMapping, envKey)
		}
	}

	// Set the updated port mappings
	if len(portMappings) > 0 {
		service.Ports = portMappings
		fmt.Printf("  Set %d port mappings for service %s\n", len(portMappings), serviceName)
	} else {
		fmt.Printf("  No port variables found for service %s\n", serviceName)
	}
}
