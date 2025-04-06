package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "env":
		consolidateCmd := flag.NewFlagSet("env", flag.ExitOnError)
		outputFile := consolidateCmd.String("o", ".env", "Output file path for consolidated env file")
		forceOverwrite := consolidateCmd.Bool("f", false, "Force overwrite output file if it exists")
		configFile := consolidateCmd.String("c", "services-config.yaml", "Path to services configuration file")
		autoDiscover := consolidateCmd.Bool("d", false, "Auto-discover services in project directory")
		serviceDir := consolidateCmd.String("dir", "", "Directory to discover services (default: current directory)")

		consolidateCmd.Parse(os.Args[2:])
		ConsolidateEnvFiles(*outputFile, *forceOverwrite, *configFile, *autoDiscover, *serviceDir)

	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		consolidatedEnvFile := updateCmd.String("env", ".env", "Path to consolidated env file")
		outputFile := updateCmd.String("o", "", "Output file path (default: docker-compose.yml in project root)")
		forceOverwrite := updateCmd.Bool("f", false, "Force overwrite output file if it exists")
		templateFile := updateCmd.String("t", "", "Path to template file (default: docker-compose.template.yml in project root)")
		discoverDir := updateCmd.String("dir", "", "Directory to discover services (default: current directory)")
		configFile := updateCmd.String("c", "services-config.yaml", "Path to services configuration file")
		updateCmd.Parse(os.Args[2:])
		UpdateDockerCompose(*consolidatedEnvFile, *outputFile, *forceOverwrite, *templateFile, *discoverDir, *configFile)

	case "help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  deployment env [options]  - Consolidate environment files")
	fmt.Println("  deployment update [options]       - Update docker-compose.yml with consolidated env vars")
	fmt.Println("  deployment help                   - Show this help message")
	fmt.Println("\nConsolidate options:")
	fmt.Println("  -o string   Output file path for consolidated env file (default: .env)")
	fmt.Println("  -f          Force overwrite output file if it exists")
	fmt.Println("  -c string   Path to services configuration file (default: services-config.yaml)")
	fmt.Println("  -d          Auto-discover services in project directory")
	fmt.Println("  -dir string Directory to discover services (default: current directory)")
	fmt.Println("              Use this to specify a different directory for service discovery")
	fmt.Println("\nUpdate options:")
	fmt.Println("  -dc string    Path to docker-compose.yml file (default: docker-compose.yml)")
	fmt.Println("  -t string  Path to template file (default: docker-compose.template.yml in project root)")
	fmt.Println("  -env string   Path to consolidated env file (default: .env)")
	fmt.Println("                Services will use this file directly via env_file directive")
	fmt.Println("  -o string     Output file path (default: docker-compose.yml in project root)")
	fmt.Println("  -f            Force overwrite output file if it exists")
	fmt.Println("  -dir string   Directory to discover services (default: current directory)")
	fmt.Println("  -c string     Path to services configuration file (default: services-config.yaml)")
	fmt.Println("\nExamples:")
	fmt.Println("  deployment env -d -o .env -f")
	fmt.Println("  deployment env -d -o .env -dir ./services")
	fmt.Println("  deployment update -t docker-compose.template.yml -env .env -o docker-compose.yml -dir ./services")
}
