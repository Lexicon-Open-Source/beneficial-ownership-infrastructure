package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ServiceConfig represents the structure of a service in the config file
type ServiceConfig struct {
	Name    string `yaml:"name"`
	EnvFile string `yaml:"env_file"`
	Prefix  string `yaml:"prefix"`
}

// Config represents the structure of the services configuration file
type Config struct {
	CommonServices []ServiceConfig `yaml:"common_services"`
	Services       []ServiceConfig `yaml:"services"`
}

// getServicePrefix retrieves the service prefix from services-config.yaml
func getServicePrefix(serviceName string, configFile string) string {
	serviceDefinition := getServiceConfig(serviceName, configFile)

	return serviceDefinition.Prefix
}

func getConfig(configFile string) Config {
	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading services config file: %v\n", err)
		return Config{}
	}

	var config Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		fmt.Printf("Error parsing services config file: %v\n", err)
		return Config{}
	}

	return config
}

func getServiceConfig(serviceName string, configFile string) ServiceConfig {
	config := getConfig(configFile)

	// Look for the service in common services
	for _, service := range config.CommonServices {
		if service.Name == serviceName {
			return service
		}
	}

	// Look for the service in regular services
	for _, service := range config.Services {
		if service.Name == serviceName {
			return service
		}
	}

	return ServiceConfig{}
}

// getServiceEnvFile retrieves the service env file path from services-config.yaml
func getServiceEnvFile(serviceName string, discoverDir string, configFile string) string {
	serviceDefinition := getServiceConfig(serviceName, configFile)
	var envFile string

	if discoverDir != "" {
		envFile = filepath.Join(discoverDir, serviceDefinition.EnvFile)
	} else {
		envFile = serviceDefinition.EnvFile
	}

	return envFile
}

func parseEnvFile(content string, envVars *map[string]string) {
	lines := strings.SplitSeq(content, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		(*envVars)[key] = value
	}
}

// Helper function to convert file paths to absolute paths
func resolveFilePath(path string, scriptDir, projectRoot string) string {
	if filepath.IsAbs(path) {
		return path
	} else if strings.HasPrefix(path, "./") {
		return filepath.Join(scriptDir, path[2:])
	} else if strings.HasPrefix(path, "../") {
		return filepath.Join(scriptDir, path)
	} else {
		// Try different locations
		// First try the project root
		projectPath := filepath.Join(projectRoot, path)
		if _, err := os.Stat(projectPath); err == nil {
			return projectPath
		}

		// Then try the script directory
		scriptPath := filepath.Join(scriptDir, path)
		if _, err := os.Stat(scriptPath); err == nil {
			return scriptPath
		}

		// If not found, default to project root
		return projectPath
	}
}
