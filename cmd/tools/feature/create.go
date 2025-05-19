package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type FeatureData struct {
	Name            string
	UpperName       string
	PluralName      string
	UpperPluralName string
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/tools/feature/main.go <feature-name>")
		os.Exit(1)
	}

	featureName := os.Args[1]
	featureName = strings.ToLower(featureName)
	upperName := strings.ToUpper(featureName)

	pluralName := featureName
	if !strings.HasSuffix(featureName, "s") {
		pluralName = featureName + "s"
	}
	upperPluralName := strings.ToUpper(pluralName)

	data := FeatureData{
		Name:            featureName,
		UpperName:       upperName,
		PluralName:      pluralName,
		UpperPluralName: upperPluralName,
	}

	featureDir := filepath.Join("internal", featureName)
	err := os.MkdirAll(featureDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Generate each file
	templateFiles := map[string]string{
		"model.go":      modelTemplate,
		"repository.go": repositoryTemplate,
		"routes.go":     routesTemplate,
		"service.go":    serviceTemplate,
	}

	for filename, templateContent := range templateFiles {
		filePath := filepath.Join(featureDir, filename)

		// Check if file already exists
		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("File %s already exists. Skipping.\n", filePath)
			continue
		}

		tmpl, err := template.New(filename).Parse(templateContent)
		if err != nil {
			fmt.Printf("Error parsing template for %s: %v\n", filename, err)
			continue
		}

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filePath, err)
			continue
		}

		err = tmpl.Execute(file, data)
		if err != nil {
			fmt.Printf("Error executing template for %s: %v\n", filename, err)
		}

		file.Close()
		fmt.Printf("Created %s\n", filePath)
	}

	fmt.Printf("\nFeature '%s' generated successfully!\n", featureName)
	fmt.Printf("Next steps:\n")
	fmt.Printf("1. Define your %s model in internal/%s/model.go\n", featureName, featureName)
	fmt.Printf("2. Implement repository methods in internal/%s/repository.go\n", featureName)
	fmt.Printf("3. Create database migration with: make migration create_%s\n", featureName)
	fmt.Printf("4. Update the main app to include your new feature\n")
}

const modelTemplate = `package {{.Name}}


`

const repositoryTemplate = `package {{.Name}}

import (
    "database/sql"

    "github.com/phsaurav/echo_prod_blueprint/internal/database"
)

// Repo is the concrete implementation of the {{.Name}} repository.
type Repo struct {
    DB *sql.DB
}

// NewRepo creates a new {{.Name}} repository instance.
func NewRepo(db database.Service) *Repo {
    return &Repo{DB: db.DB()}
}

var _ Repository = (*Repo)(nil)


`

const routesTemplate = `package {{.Name}}

import (
    "github.com/labstack/echo/v4"
    "github.com/phsaurav/echo_prod_blueprint/internal/database"
)

type {{.UpperName}}Service interface {

}

func Register(g *echo.Group, db database.Service, authMiddleware echo.MiddlewareFunc) {
    repo := NewRepo(db)
    service := NewService(repo)
    RegisterRoutes(g, service, authMiddleware)
}

func RegisterRoutes(g *echo.Group, service {{.UpperName}}Service, authMiddleware echo.MiddlewareFunc) {

}
`

const serviceTemplate = `package {{.Name}}

import (

)

type Repository interface {

}

// Service implements the consumer-side {{.UpperName}}Service interface.
type Service struct {
    Repo Repository
}

// NewService creates a new {{.Name}} service instance.
func NewService(repo Repository) *Service {
    return &Service{Repo: repo}
}


`
