package scaffold

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type workflowTemplateValues struct {
	Name             string
	PackageName      string
	Description      string
	DeveloperName    string
	DeveloperEmail   string
	URLRoot          string
	URLPath          string
	IsWizardTemplate bool
}

type FrontendPluginScaffoldWorkflow struct {
	Destination string
	Data        *workflowTemplateValues
}

func (f *FrontendPluginScaffoldWorkflow) PromptValues() {
	templateValues, dest := GetFrontendPluginTemplateValues()
	f.Destination = dest
	f.Data = templateValues
}

func (f *FrontendPluginScaffoldWorkflow) GetTemplateDirectory() string {
	return "templates/frontend/workflow/external"
}

func (f *FrontendPluginScaffoldWorkflow) GetTemplateValues() interface{} {
	return f.Data
}

func (f *FrontendPluginScaffoldWorkflow) GetDestinationDirectories() []string {
	return []string{f.Destination}
}

func (f *FrontendPluginScaffoldWorkflow) PostProcess(flags *Args, tmpFolder string) {
	if strings.Contains(f.Destination, "/clutch/frontend/workflows/") || flags.Internal {
		PostProcessFrontendInternal(flags, tmpFolder, f.Destination)
	} else {
		PostProcessFrontend(flags, tmpFolder, f.Destination)
	}
}

func GetFrontendPluginTemplateValues() (*workflowTemplateValues, string) {
	log.Println("Welcome!")
	fmt.Println("*** Analyzing environment...")

	dest := filepath.Join(os.Getenv("OLDPWD"), "frontend", "workflows")

	fmt.Println("\n*** Based on your environment, we've picked the following destination for your new workflow:")
	fmt.Println(">", dest)
	okay := promptOrDefault("Is this okay?", "Y/n")
	if !strings.HasPrefix(strings.ToLower(okay), "y") {
		dest = promptOrDefault("Enter the destination folder", dest)
	}

	data := &workflowTemplateValues{}

	data.IsWizardTemplate = true

	wizard := promptOrDefault("Is this a wizard workflow?", "Y/n")
	if strings.HasPrefix(strings.ToLower(wizard), "n") {
		data.IsWizardTemplate = false
	}

	data.Name = strings.Title(promptOrDefault("Enter the name of this workflow", "Hello World"))
	description := promptOrDefault("Enter a description of the workflow", "Greet the world")
	data.Description = strings.ToUpper(description[:1]) + description[1:]
	data.DeveloperName = promptOrDefault("Enter the developer's name", determineUsername())
	data.DeveloperEmail = promptOrDefault("Enter the developer's email", determineUserEmail())

	// n.b. transform workflow name into package name, e.g. foo bar baz -> fooBarBaz
	packageName := strings.ToLower(data.Name[:1]) + strings.Title(data.Name)[1:]
	data.PackageName = strings.Replace(packageName, " ", "", -1)

	data.URLRoot = strings.Replace(strings.ToLower(data.Name), " ", "", -1)
	data.URLPath = "/"

	return data, filepath.Join(dest, data.PackageName)
}

func GenerateFrontend(args *Args, tmpFolder string, dest string) {
	fmt.Println("Compiling workflow, this may take a few minutes...")
	yarn := DetermineYarnPath()
	fmt.Println("cd", tmpFolder, "&&", yarn, "install &&", yarn, "tsc && ", yarn, "compile")
	if err := os.Chdir(tmpFolder); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running yarn install in", tmpFolder, "with yarn path", yarn)
	installCmd := exec.Command(yarn, "install")
	if out, err := installCmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Println("`yarn install` returned the above error")

		YarnInstall(yarn, args.YarnPin)
	}

	fmt.Println("Running yarn tsc in", tmpFolder, "with yarn path", yarn)
	compileTypesCmd := exec.Command(yarn, "tsc")
	if out, err := compileTypesCmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`yarn tsc` returned the above error")
	}

	fmt.Println("Running yarn compile in", tmpFolder, "with yarn path", yarn)
	compileDevCmd := exec.Command(yarn, "compile")
	if out, err := compileDevCmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`yarn compile` returned the above error")
	}

	fmt.Println("*** All done!")
	fmt.Printf("\n*** Your new workflow can be found here: %s\n", dest)
	fmt.Println("For information on how to register this new workflow see our configuration guide: https://clutch.sh/docs/configuration")
}

func GenerateInternalFrontend(args *Args, tmpFolder string, dest string) {
	root := os.Getenv("OLDPWD")
	feRoot := filepath.Join(root, "frontend")

	fmt.Println("Compiling workflow, this may take a few minutes...")

	fmt.Println("Running `make frontend-install` in", root)
	installCmd := exec.Command("make", "frontend-install")
	installCmd.Dir = root
	if out, err := installCmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`make frontend-install` returned the above error")
	}

	// Move tmpdir contents to destination.
	MoveTempFilesToDest(tmpFolder, dest)

	fmt.Println("Running `yarn install` in", feRoot)
	yarnInstallCmd := exec.Command("yarn", "install")
	yarnInstallCmd.Dir = feRoot
	if out, err := yarnInstallCmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`yarn install` returned the above error")
	}

	fmt.Println("Running `make frontend-compile` in", root)
	compileDevCmd := exec.Command("make", "frontend-compile")
	compileDevCmd.Dir = root
	if out, err := compileDevCmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`make frontend-compile` returned the above error")
	}

	fmt.Println("*** All done!")
	fmt.Printf("\n*** Your new workflow can be found here: %s\n", dest)
	fmt.Println("For information on how to register this new workflow see our configuration guide: https://clutch.sh/docs/configuration")
}

func PostProcessFrontend(flags *Args, tmpFolder string, dest string) {
	GenerateFrontend(flags, tmpFolder, dest)

	MoveTempFilesToDest(tmpFolder, dest)
}

func PostProcessFrontendInternal(flags *Args, tmpFolder string, dest string) {
	GenerateInternalFrontend(flags, tmpFolder, dest)
}
