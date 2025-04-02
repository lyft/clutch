package scaffold

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type RegistrationParams struct {
	BackendMainPath string
	ComponentType   string
	ComponentName   string
	ComponentPath   string
}

func RegisterNewComponent(params *RegistrationParams) {
	log.Printf("Adding new %s into %s: %s\n", params.ComponentType, params.BackendMainPath, params.ComponentName)

	importLine := generateImportLine(params)
	registrationLine := generateRegistrationLine(params)

	input, err := os.ReadFile(params.BackendMainPath)
	if err != nil {
		log.Fatalf("failed to read backend main file: %v", err)
	}

	lines := strings.Split(string(input), "\n")

	// Find the import section and add the import line
	importInserted := insertImportLine(&lines, importLine)
	if !importInserted {
		log.Fatalf("failed to find import section in backend main file")
	}

	// Find the registration section and add the registration line
	registrationInserted := insertRegistrationLine(&lines, registrationLine, params.ComponentType)
	if !registrationInserted {
		log.Fatalf("failed to find registration section in backend main file")
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(params.BackendMainPath, []byte(output), 0644)
	if err != nil {
		log.Fatalf("failed to write updated backend main file: %v", err)
	}
}

func generateImportLine(params *RegistrationParams) string {
	// Generates: sampleservice "github.com/lyft/clutch/backend/service/sample"
	importLine := `%s%s%s "%s"`
	indent := "    "
	return fmt.Sprintf(
		importLine,
		indent,
		params.ComponentName,
		params.ComponentType,
		params.ComponentPath,
	)
}

func generateRegistrationLine(params *RegistrationParams) string {
	// Generates: components.Services[sampleservice.Name] = sampleservice.New
	registrationLine := `%scomponents.%s[%s%s.Name] = %s%s.New`
	indent := "    "
	compTitlePlural := strings.ToUpper(params.ComponentType[:1]) + params.ComponentType[1:] + "s"
	return fmt.Sprintf(
		registrationLine,
		indent,
		compTitlePlural,
		params.ComponentName,
		params.ComponentType,
		params.ComponentName,
		params.ComponentType,
	)
}

func insertImportLine(lines *[]string, importLine string) bool {
	for i, line := range *lines {
		if strings.HasPrefix(line, "import (") {
			for j := i + 1; j < len(*lines); j++ {
				if strings.HasPrefix((*lines)[j], ")") {
					*lines = slices.Insert(*lines, j, importLine)
					return true
				}
			}
			break
		}
	}
	return false
}

func insertRegistrationLine(lines *[]string, registrationLine, componentType string) bool {
	compTitlePlural := strings.ToUpper(componentType[:1]) + componentType[1:] + "s"
	registrationPattern := fmt.Sprintf(`components.%s[`, compTitlePlural)

	for i, line := range *lines {
		if strings.Contains(line, registrationPattern) {
			for j := i + 1; j < len(*lines); j++ {
				if !strings.Contains((*lines)[j], registrationPattern) {
					*lines = slices.Insert(*lines, j, registrationLine)
					return true
				}
			}
			break
		} else if strings.Contains(line, "gateway.Run(") {
			*lines = slices.Insert(*lines, i, registrationLine, "")
			return true
		}
	}

	return false
}
