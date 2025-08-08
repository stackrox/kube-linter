package main

import (
	"fmt"
	"log"
	"os"

	"golang.stackrox.io/kube-linter/pkg/run"
)

func main() {
	// Example 1: Basic validation with a simple policy
	exampleBasicValidation()

	// Example 2: Kubernetes manifest validation
	exampleKubernetesValidation()

	// Example 3: Custom policy validation
	exampleCustomPolicy()
}

func exampleBasicValidation() {
	fmt.Println("=== Example 1: Basic Validation ===")

	// Create a simple policy
	policy := `
package example

allow {
    input.value == "allowed"
}

deny[msg] {
    input.value != "allowed"
    msg := "Value must be 'allowed'"
}
`

	// Create test data
	yamlData := `
value: "allowed"
`

	// Create validator
	validator, err := run.NewOPAValidator()
	if err != nil {
		log.Fatal(err)
	}
	defer validator.Close()

	// Load policy from string
	err = validator.LoadPolicyFromString(policy, "example.rego")
	if err != nil {
		log.Fatal(err)
	}

	// Validate
	result, err := validator.ValidateYAML([]byte(yamlData), "data.example")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %+v\n", result)
	fmt.Println()
}

func exampleKubernetesValidation() {
	fmt.Println("=== Example 2: Kubernetes Manifest Validation ===")

	// Create a Kubernetes security policy
	policy := `
package kubernetes.admission

deny[msg] {
    input.kind == "Pod"
    not input.spec.securityContext.runAsNonRoot

    msg := "Pods must not run as root"
}

deny[msg] {
    input.kind == "Pod"
    container := input.spec.containers[_]
    not container.resources.limits.memory

    msg := sprintf("Container %v must have memory limits", [container.name])
}

allow {
    input.kind != "Pod"
}

allow {
    input.kind == "Pod"
    input.spec.securityContext.runAsNonRoot
    all_containers_have_limits
}

all_containers_have_limits {
    container := input.spec.containers[_]
    container.resources.limits.memory
}
`

	// Create a Kubernetes Pod manifest
	podYAML := `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  securityContext:
    runAsNonRoot: true
  containers:
  - name: nginx
    image: nginx:latest
    resources:
      limits:
        memory: "128Mi"
        cpu: "500m"
      requests:
        memory: "64Mi"
        cpu: "250m"
`

	// Create validator
	validator, err := run.NewOPAValidator()
	if err != nil {
		log.Fatal(err)
	}
	defer validator.Close()

	// Load policy
	err = validator.LoadPolicyFromString(policy, "k8s-policy.rego")
	if err != nil {
		log.Fatal(err)
	}

	// Validate
	result, err := validator.ValidateYAML([]byte(podYAML), "data.kubernetes.admission")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pod validation result: %+v\n", result)

	// Test with an invalid pod
	invalidPodYAML := `
apiVersion: v1
kind: Pod
metadata:
  name: invalid-pod
spec:
  containers:
  - name: nginx
    image: nginx:latest
`

	result, err = validator.ValidateYAML([]byte(invalidPodYAML), "data.kubernetes.admission")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Invalid pod validation result: %+v\n", result)
	fmt.Println()
}

func exampleCustomPolicy() {
	fmt.Println("=== Example 3: Custom Policy Validation ===")

	// Create a custom policy for application configuration
	policy := `
package app.config

deny[msg] {
    input.environment == "production"
    not input.database.backup_enabled

    msg := "Production environments must have database backups enabled"
}

deny[msg] {
    input.environment == "production"
    input.database.max_connections > 100

    msg := "Production database connections should not exceed 100"
}

allow {
    input.environment != "production"
}

allow {
    input.environment == "production"
    input.database.backup_enabled
    input.database.max_connections <= 100
}
`

	// Create application configuration
	configYAML := `
environment: "production"
database:
  backup_enabled: true
  max_connections: 50
  host: "db.example.com"
  port: 5432
`

	// Create validator
	validator, err := run.NewOPAValidator()
	if err != nil {
		log.Fatal(err)
	}
	defer validator.Close()

	// Load policy
	err = validator.LoadPolicyFromString(policy, "app-policy.rego")
	if err != nil {
		log.Fatal(err)
	}

	// Validate
	result, err := validator.ValidateYAML([]byte(configYAML), "data.app.config")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("App config validation result: %+v\n", result)

	// Test with invalid configuration
	invalidConfigYAML := `
environment: "production"
database:
  backup_enabled: false
  max_connections: 150
  host: "db.example.com"
  port: 5432
`

	result, err = validator.ValidateYAML([]byte(invalidConfigYAML), "data.app.config")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Invalid app config validation result: %+v\n", result)
	fmt.Println()
}

// Helper function to create a policy file
func createPolicyFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// Helper function to create a YAML file
func createYAMLFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
