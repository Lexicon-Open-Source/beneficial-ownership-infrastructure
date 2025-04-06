package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ConsolidateEnvFiles consolidates environment files from services into a single file
func ConsolidateEnvFiles(outputFile string, forceOverwrite bool, configFile string, autoDiscover bool, serviceDir string) {
	// Get script directory
	scriptDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Set project root (parent of script directory)
	projectRoot := scriptDir

	// Use serviceDir if provided, otherwise use project root
	discoverDir := projectRoot
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

	// If output file path is relative, make it absolute from script directory
	if !filepath.IsAbs(outputFile) {
		if strings.HasPrefix(outputFile, "./") {
			outputFile = filepath.Join(scriptDir, outputFile[2:])
		} else if strings.HasPrefix(outputFile, "../") {
			outputFile = filepath.Join(scriptDir, outputFile)
		} else {
			outputFile = filepath.Join(scriptDir, outputFile)
		}
	}

	// If config file path is relative, make it absolute from script directory
	if !filepath.IsAbs(configFile) {
		if strings.HasPrefix(configFile, "./") {
			configFile = filepath.Join(scriptDir, configFile[2:])
		} else if strings.HasPrefix(configFile, "../") {
			configFile = filepath.Join(scriptDir, configFile)
		} else {
			configFile = filepath.Join(scriptDir, configFile)
		}
	}

	// Debug info
	fmt.Printf("Script directory: %s\n", scriptDir)
	fmt.Printf("Project root: %s\n", projectRoot)
	fmt.Printf("Output file: %s\n", outputFile)
	fmt.Printf("Config file: %s\n", configFile)

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

	// Initialize variables to store services
	var commonServices []ServiceConfig
	var appServices []ServiceConfig

	// Load configuration from YAML file if it exists
	if !autoDiscover {
		fmt.Printf("Loading services configuration from %s\n", configFile)

		// Parse config
		config := getConfig(configFile)

		// Process common services
		for _, service := range config.CommonServices {
			// Always enforce .env in subfolder of project root
			envFile := filepath.Join(discoverDir, service.Name, ".env")

			// Check if the file exists
			if _, err := os.Stat(envFile); err != nil {
				fmt.Printf("Warning: Env file not found at %s for service %s\n", envFile, service.Name)
				continue
			}

			service.EnvFile = envFile
			fmt.Printf("Configured common service env file: %s\n", service.EnvFile)
			commonServices = append(commonServices, service)
		}

		// Process app services
		for _, service := range config.Services {
			// Always enforce .env in subfolder of project root
			envFile := filepath.Join(discoverDir, service.Name, ".env")

			// Check if the file exists
			if _, err := os.Stat(envFile); err != nil {
				fmt.Printf("Warning: Env file not found at %s for service %s\n", envFile, service.Name)
				continue
			}

			service.EnvFile = envFile
			fmt.Printf("Configured application service env file: %s\n", service.EnvFile)
			appServices = append(appServices, service)
		}
	} else if autoDiscover {
		fmt.Println("Auto-discovering services...")

		// Load config for prefixes if available
		var config Config
		if _, err := os.Stat(configFile); err == nil {
			configData, err := os.ReadFile(configFile)
			if err == nil {
				yaml.Unmarshal(configData, &config)
			}
		}

		// Find directories with .env files
		dirs, err := os.ReadDir(discoverDir)
		if err != nil {
			fmt.Printf("Error reading directory: %v\n", err)
			return
		}

		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}

			serviceName := dir.Name()
			envFile := filepath.Join(discoverDir, serviceName, ".env")

			// Log env file search
			fmt.Printf("Checking for .env file in %s\n", filepath.Join(discoverDir, serviceName))

			if _, err := os.Stat(envFile); err == nil {
				fmt.Printf("Found .env file at %s\n", envFile)

				// Try to get prefix from config
				prefix := ""
				isCommon := false

				// Check if it's in common services
				for _, s := range config.CommonServices {
					if s.Name == serviceName {
						prefix = s.Prefix
						isCommon = true
						break
					}
				}

				// If not in common, check app services
				if prefix == "" {
					for _, s := range config.Services {
						if s.Name == serviceName {
							prefix = s.Prefix
							break
						}
					}
				}

				// Use default prefix naming convention if not found
				if prefix == "" {
					prefix = strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")) + "_"
					fmt.Printf("Using auto-generated prefix %s for %s\n", prefix, serviceName)
				}

				// Create service config
				service := ServiceConfig{
					Name:    serviceName,
					EnvFile: envFile,
					Prefix:  prefix,
				}

				// Add to appropriate list
				// First check if it's already identified as a common service through the config
				if isCommon {
					fmt.Printf("Identified %s as a common infrastructure service from config\n", serviceName)
					commonServices = append(commonServices, service)
				} else {
					// Load common service names from config file if available
					isConfiguredCommonService := false

					// Check if we have loaded a valid config
					if _, err := os.Stat(configFile); err == nil {
						configData, err := os.ReadFile(configFile)
						if err == nil {
							var configFromFile Config
							if err := yaml.Unmarshal(configData, &configFromFile); err == nil {
								// Check if service name exists in common_services
								for _, cs := range configFromFile.CommonServices {
									if cs.Name == serviceName {
										isConfiguredCommonService = true
										fmt.Printf("Identified %s as a common infrastructure service from config file\n", serviceName)
										break
									}
								}
							}
						}
					}

					if isConfiguredCommonService {
						commonServices = append(commonServices, service)
					} else {
						appServices = append(appServices, service)
					}
				}
			}
		}

		fmt.Printf("Discovered %d services (%d common, %d application)\n",
			len(commonServices)+len(appServices), len(commonServices), len(appServices))
	} else {
		fmt.Println("Warning: Config file not found and auto-discovery is disabled.")
		fmt.Println("No services will be processed. Use -d or provide a valid config file.")
		return
	}

	// Create consolidated file
	outputFile, err = createConsolidatedFile(outputFile, commonServices, appServices)
	if err != nil {
		fmt.Printf("Error creating consolidated file: %v\n", err)
		return
	}

	fmt.Printf("Consolidated .env file created at %s\n", outputFile)
}

func createConsolidatedFile(outputPath string, commonServices, appServices []ServiceConfig) (string, error) {
	// Create or truncate output file
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	// Write header
	header := fmt.Sprintf("# Consolidated .env file\n"+
		"# Generated on: %s\n"+
		"# This file was automatically generated by consolidating service-specific .env files\n"+
		"# DO NOT EDIT THIS FILE DIRECTLY - Edit individual service .env files instead\n\n",
		time.Now().Format(time.RFC1123))
	file.WriteString(header)

	// Map to track processed variables
	processedVars := make(map[string]bool)

	// Process common services
	file.WriteString("\n# === COMMON INFRASTRUCTURE VARIABLES ===\n\n")
	commonVarCount, commonSuccessCount := 0, 0

	for _, service := range commonServices {
		// Add section header
		file.WriteString(fmt.Sprintf("# %s environment variables\n", service.Name))

		varCount, err := processEnvFile(file, service, processedVars, true)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", service.EnvFile, err)
			continue
		}

		file.WriteString("\n")
		fmt.Printf("Processed common service %s - added %d variables\n", service.EnvFile, varCount)

		if varCount > 0 {
			commonVarCount += varCount
			commonSuccessCount++
		}
	}

	// Process application services
	file.WriteString("\n# === APPLICATION-SPECIFIC VARIABLES ===\n\n")
	appVarCount, appSuccessCount := 0, 0

	for _, service := range appServices {
		// Add section header
		file.WriteString(fmt.Sprintf("# %s environment variables\n", service.Name))

		varCount, err := processEnvFile(file, service, processedVars, false)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", service.EnvFile, err)
			continue
		}

		file.WriteString("\n")
		fmt.Printf("Processed %s - added %d variables\n", service.EnvFile, varCount)

		if varCount > 0 {
			appVarCount += varCount
			appSuccessCount++
		}
	}

	fmt.Printf("Environment file consolidation completed.\n")
	fmt.Printf("%d common infrastructure .env files processed with %d variables.\n", commonSuccessCount, commonVarCount)
	fmt.Printf("%d application-specific .env files processed with %d variables.\n", appSuccessCount, appVarCount)
	fmt.Printf("Consolidated %d unique environment variables.\n", commonVarCount+appVarCount)

	return outputPath, nil
}

func processEnvFile(file *os.File, service ServiceConfig, processedVars map[string]bool, isCommon bool) (int, error) {
	// Check if file exists
	if _, err := os.Stat(service.EnvFile); err != nil {
		return 0, fmt.Errorf("env file not found: %s", service.EnvFile)
	}

	fmt.Printf("Processing .env file: %s\n", service.EnvFile)

	// Open file
	envFile, err := os.Open(service.EnvFile)
	if err != nil {
		return 0, err
	}
	defer envFile.Close()

	// Process lines
	scanner := bufio.NewScanner(envFile)
	varCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Extract variable name
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		varName := strings.TrimSpace(parts[0])
		varValue := strings.TrimSpace(parts[1])

		// Check if it already has prefix
		var prefixedVar string
		if strings.HasPrefix(varName, service.Prefix) {
			prefixedVar = varName
			line = fmt.Sprintf("%s= %s", prefixedVar, varValue)
		} else {
			prefixedVar = service.Prefix + varName
			line = fmt.Sprintf("%s= %s", prefixedVar, varValue)
		}

		// Check for duplicates
		if _, exists := processedVars[prefixedVar]; !exists {
			// New variable, add to consolidated file
			file.WriteString(line + "\n")
			processedVars[prefixedVar] = true
			varCount++
		} else if !isCommon {
			// For application services, show warning when they use variables from common services
			fmt.Printf("Note: Variable %s is already defined in a common service\n", prefixedVar)
		}
	}

	if err := scanner.Err(); err != nil {
		return varCount, err
	}

	return varCount, nil
}
