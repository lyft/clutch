package scaffold

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type backendApiTemplateValues struct {
	ApiName          string
	ApiVersion       string
	ModuleName       string
	ServiceName      string
	Description      string
	BackendGoPackage string
	RepoOwner        string
}

type BackendApiScaffoldWorkflow struct {
	Name              string
	ApiDestination    string
	ModuleDestination string
	RepositoryRoot    string
	Data              *backendApiTemplateValues
	RegistrationData  *RegistrationParams
	ServiceWorkflow   *ServiceScaffoldWorkflow
}

func (b *BackendApiScaffoldWorkflow) GetTemplateDirectory() string {
	return "templates/backend-plugin"
}

func (b *BackendApiScaffoldWorkflow) GetTemplateValues() interface{} {
	return b.Data
}

func (b *BackendApiScaffoldWorkflow) GetDestinationDirectories() []string {
	directories := []string{b.ApiDestination, b.ModuleDestination}
	if b.ServiceWorkflow != nil {
		directories = append(directories, b.ServiceWorkflow.GetDestinationDirectories()...)
	}
	return directories
}

func (b *BackendApiScaffoldWorkflow) PromptValues() {
	log.Println("Welcome!")
	fmt.Println("*** Analyzing environment...")

	b.RepositoryRoot = filepath.Join(os.Getenv("OLDPWD")) // $OLDPWD gets overwritten after this

	fmt.Println("\n*** Based on your environment, we've picked the following information:")
	fmt.Println("> Repository root:", b.RepositoryRoot)
	okay := promptOrDefault("Is this okay?", "Y/n")
	if !strings.HasPrefix(strings.ToLower(okay), "y") {
		b.RepositoryRoot = promptOrDefault("Enter the path to the repository root folder", b.RepositoryRoot)
	}
	b.Data = &backendApiTemplateValues{}

	apiName := promptOrDefault("Enter the name of this backend API", "HelloWorld")
	b.Data.ApiName = strings.ToUpper(apiName[:1]) + apiName[1:]
	b.Data.ModuleName = strings.ToLower(b.Data.ApiName)
	b.Data.ApiVersion = promptOrDefault("Enter the API version", "v1")
	description := promptOrDefault("Enter a description for the API", "Greet the world")
	b.Data.Description = strings.ToUpper(description[:1]) + description[1:]
	goPackage := determineGoPackage(filepath.Join(b.RepositoryRoot, "backend"))
	b.Data.BackendGoPackage = promptOrDefault("Enter go backend package path", goPackage)
	gitUpstream := determineGitUpstream(b.RepositoryRoot)
	b.Data.RepoOwner = gitUpstream.RepoOwner
	if b.Data.BackendGoPackage == "github.com/lyft/clutch/backend" {
		b.Data.RepoOwner = "clutch"
	}
	b.Data.RepoOwner = promptOrDefault("Enter the workflow owner", b.Data.RepoOwner)

	b.RegistrationData = &RegistrationParams{
		BackendMainPath: filepath.Join(b.RepositoryRoot, "backend", "main.go"),
		ComponentType:   "module",
		ComponentName:   b.Data.ModuleName,
		ComponentPath:   filepath.Join(b.Data.BackendGoPackage, "module", b.Data.ModuleName),
	}

	b.Name = b.Data.ModuleName
	b.ApiDestination = filepath.Join(b.RepositoryRoot, "api", b.Name, b.Data.ApiVersion)
	b.ModuleDestination = filepath.Join(b.RepositoryRoot, "backend/module", b.Name)
	okay = promptOrDefault("Would you like to generate a new service for this API?", "Y/n")
	if strings.HasPrefix(strings.ToLower(okay), "y") {
		b.ServiceWorkflow = &ServiceScaffoldWorkflow{
			ServiceName: b.Name,
			Destination: filepath.Join(b.RepositoryRoot, "backend/service", b.Name),
			RegistrationData: &RegistrationParams{
				ComponentType:   "service",
				ComponentName:   b.Name,
				ComponentPath:   filepath.Join(b.Data.BackendGoPackage, "service", b.Name),
				BackendMainPath: filepath.Join(b.RepositoryRoot, "backend", "main.go"),
			},
		}
		b.Data.ServiceName = b.Name
	}
}

func (b *BackendApiScaffoldWorkflow) PostProcess(flags *Args, tmpFolder string) {
	MoveTempFilesToDest(filepath.Join(tmpFolder, "api/version/api.proto"), filepath.Join(b.ApiDestination, b.Name+".proto"))
	MoveTempFilesToDest(filepath.Join(tmpFolder, "backend/module/module.go"), filepath.Join(b.ModuleDestination, b.Name+".go"))
	MoveTempFilesToDest(filepath.Join(tmpFolder, "backend/module/module_test.go"), filepath.Join(b.ModuleDestination, b.Name+"_test.go"))
	if b.ServiceWorkflow != nil {
		b.ServiceWorkflow.PostProcess(flags, filepath.Join(tmpFolder, "backend/service"))
	}
	RegisterNewComponent(b.RegistrationData)
	MakeAPI(b.RepositoryRoot)
}
