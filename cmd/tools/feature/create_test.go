package main

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestFeatureGeneration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "feature-test")
	if err != nil {
		t.Fatalf("Could not create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save current directory and change to temp directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	os.Chdir(tempDir)

	// Create internal directory
	internalDir := filepath.Join(tempDir, "internal")
	if err := os.MkdirAll(internalDir, 0755); err != nil {
		t.Fatalf("Could not create internal directory: %v", err)
	}

	// Set up test feature data
	featureName := "test"
	data := FeatureData{
		Name:            featureName,
		UpperName:       "TEST",
		PluralName:      "tests",
		UpperPluralName: "TESTS",
	}

	// Create feature directory
	featureDir := filepath.Join(internalDir, featureName)
	if err := os.MkdirAll(featureDir, 0755); err != nil {
		t.Fatalf("Could not create feature directory: %v", err)
	}

	// Create templates
	templates := map[string]string{
		"model.go":      modelTemplate,
		"repository.go": repositoryTemplate,
		"routes.go":     routesTemplate,
		"service.go":    serviceTemplate,
	}

	// Execute templates and verify files were created
	for filename, templateContent := range templates {
		filePath := filepath.Join(featureDir, filename)

		// Create file from template
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Could not create file %s: %v", filename, err)
		}

		tmpl, err := template.New(filename).Parse(templateContent)
		if err != nil {
			t.Fatalf("Could not parse template %s: %v", filename, err)
		}

		err = tmpl.Execute(file, data)
		file.Close()
		if err != nil {
			t.Fatalf("Could not execute template %s: %v", filename, err)
		}

		// Verify file exists
		_, err = os.Stat(filePath)
		assert.NoError(t, err, "File should exist: %s", filename)

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Could not read file %s: %v", filename, err)
		}

		// Verify content contains expected strings
		contentStr := string(content)
		assert.Contains(t, contentStr, "package test", "File should have correct package name")
	}
}
