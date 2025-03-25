package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	scaffold "github.com/lyft/clutch/tools/scaffolding/scaffold"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	flags := scaffold.ParseArgs()

	// Collect info from user based on mode and determine template root.
	var dest string
	var templateRoot string
	var templateOverwrites string
	var data interface{}
	var postProcessFunction func(flags *scaffold.Args, tmpFolder string, dest string)
	var registrationData *scaffold.RegistrationParams
	var newComponent bool = false

	switch flags.Mode {
	case "gateway":
		templateRoot = filepath.Join(root, "templates/gateway")
		data, dest = scaffold.GetGatewayTemplateValues()
		postProcessFunction = scaffold.PostProcessGateway
	case "service":
		newComponent = true
		templateRoot = filepath.Join(root, "templates/backend/service")
		data, dest, registrationData = scaffold.GetServiceTemplateValues()
		postProcessFunction = scaffold.PostProcessService
	case "frontend-plugin":
		templateRoot = filepath.Join(root, "templates/frontend/workflow/internal")
		data, dest = scaffold.GetFrontendPluginTemplateValues()
		if strings.Contains(dest, "/clutch/frontend/workflows/") || flags.Internal {
			postProcessFunction = scaffold.PostProcessFrontendInternal
		} else {
			templateOverwrites = filepath.Join(root, "templates/frontend/workflow/external")
			postProcessFunction = scaffold.PostProcessFrontend
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	scaffold.FixDestDir(dest)

	// Make a tmpdir for output.
	fmt.Println("\n*** Generating...")
	fmt.Println("Using templates in", templateRoot)
	tmpout := scaffold.TemplateTempDir(templateRoot, data)
	if flags.TemplateOverwrite != "" {
		fmt.Println("\n*** Overwriting files with custom templates...", tmpout)
		scaffold.TemplateFiles(flags.TemplateOverwrite, tmpout, data)
	} else if templateOverwrites != "" {
		fmt.Println("\n*** Overwriting files with custom templates...", tmpout)
		scaffold.TemplateFiles(templateOverwrites, tmpout, data)
	}
	defer os.RemoveAll(tmpout)

	if newComponent {
		scaffold.RegisterNewComponent(registrationData)
	}

	postProcessFunction(flags, tmpout, dest)
}
