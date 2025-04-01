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
	var destinations []string
	var templateRoot string
	var templateOverwrites string
	var data interface{}
	var scaffoldWorkflow scaffold.ScaffoldWorkflow

	switch flags.Mode {
	case "gateway":
		scaffoldWorkflow = &scaffold.GatewayScaffoldWorkflow{}
		scaffoldWorkflow.PromptValues()
	case "service":
		scaffoldWorkflow = &scaffold.ServiceScaffoldWorkflow{}
		scaffoldWorkflow.PromptValues()
	case "frontend-plugin":
		feWorkflow := &scaffold.FrontendPluginScaffoldWorkflow{}
		feWorkflow.PromptValues()
		if !strings.Contains(feWorkflow.Destination, "/clutch/frontend/workflows/") && !flags.Internal {
			templateOverwrites = filepath.Join(root, "templates/frontend/workflow/external")
		}
		scaffoldWorkflow = feWorkflow
	default:
		flag.Usage()
		os.Exit(1)
	}

	data = scaffoldWorkflow.GetTemplateValues()
	destinations = scaffoldWorkflow.GetDestinationDirectories()
	templateRoot = filepath.Join(root, scaffoldWorkflow.GetTemplateDirectory())

	for _, dest := range destinations {
		scaffold.FixDestDir(dest)
	}

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

	scaffoldWorkflow.PostProcess(flags, tmpout)
}
