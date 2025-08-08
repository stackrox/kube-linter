package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.stackrox.io/kube-linter/pkg/run"
)

func main() {
	var (
		policyFile = flag.String("policy", "", "Path to the Rego policy file")
		yamlFile   = flag.String("yaml", "", "Path to the YAML file to validate")
		query      = flag.String("query", "", "OPA query path (e.g., data.kubernetes.admission)")
		output     = flag.String("output", "text", "Output format: text, json")
		verbose    = flag.Bool("verbose", false, "Enable verbose output")
	)

	flag.Parse()

	if *policyFile == "" || *yamlFile == "" || *query == "" {
		fmt.Println("Usage: opa-validator -policy <policy.rego> -yaml <file.yaml> -query <query.path>")
		fmt.Println("Example: opa-validator -policy security.rego -yaml pod.yaml -query data.kubernetes.admission")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create validator
	validator, err := run.NewOPAValidator()
	if err != nil {
		log.Fatalf("Failed to create OPA validator: %v", err)
	}
	defer validator.Close()

	// Load policy
	if *verbose {
		fmt.Printf("Loading policy from: %s\n", *policyFile)
	}

	err = validator.LoadPolicy(*policyFile)
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}

	// Validate YAML
	if *verbose {
		fmt.Printf("Validating YAML file: %s\n", *yamlFile)
		fmt.Printf("Using query path: %s\n", *query)
	}

	result, err := validator.ValidateYAMLFile(*yamlFile, *query)
	if err != nil {
		log.Fatalf("Failed to validate YAML: %v", err)
	}

	// Output results
	switch *output {
	case "json":
		outputJSON(result)
	case "text":
		outputText(result)
	default:
		log.Fatalf("Unknown output format: %s", *output)
	}

	// Exit with appropriate code
	if result.Allowed {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func outputJSON(result *run.ValidationResult) {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

func outputText(result *run.ValidationResult) {
	if result.Allowed {
		fmt.Println("✅ Validation passed")
	} else {
		fmt.Println("❌ Validation failed")

		if len(result.Results) > 0 {
			fmt.Println("\nViolations:")
			for i, violation := range result.Results {
				fmt.Printf("  %d. %v\n", i+1, violation)
			}
		}

		if len(result.Errors) > 0 {
			fmt.Println("\nErrors:")
			for i, err := range result.Errors {
				fmt.Printf("  %d. %s\n", i+1, err)
			}
		}
	}

	if len(result.Data) > 0 {
		fmt.Println("\nAdditional data:")
		for key, value := range result.Data {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}
