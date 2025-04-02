package scaffold

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type serviceTemplateValues struct {
	ServiceName string
	Description string
	RepoOwner   string
}

type ServiceScaffoldWorkflow struct {
	ServiceName      string
	Destination      string
	RegistrationData *RegistrationParams
	Data             *serviceTemplateValues
}

func (s *ServiceScaffoldWorkflow) GetTemplateDirectory() string {
	return "templates/backend-plugin/backend/service"
}

func (s *ServiceScaffoldWorkflow) GetDestinationDirectories() []string {
	return []string{s.Destination}
}

func (s *ServiceScaffoldWorkflow) PromptValues() {
	log.Println("Welcome!")
	fmt.Println("*** Analyzing environment...")

	repositoryRoot := os.Getenv("OLDPWD")
	dest := filepath.Join(repositoryRoot, "backend", "service")
	fmt.Println("\n*** Based on your environment, we've picked the following destination for your new service:")
	fmt.Println(">", dest)
	okay := promptOrDefault("Is this okay?", "Y/n")
	if !strings.HasPrefix(strings.ToLower(okay), "y") {
		dest = promptOrDefault("Enter the destination folder", dest)
	}

	s.Data = &serviceTemplateValues{}

	s.Data.ServiceName = strings.ToLower(promptOrDefault("Enter the name of this service", "helloworld"))
	description := promptOrDefault("Enter a description of the service", "Greet the world")
	s.Data.Description = strings.ToUpper(description[:1]) + description[1:]
	gitUpstream := determineGitUpstream(repositoryRoot)
	s.Data.RepoOwner = gitUpstream.RepoOwner
	if determineGoPackage(filepath.Join(repositoryRoot, "backend")) == "github.com/lyft/clutch/backend" {
		s.Data.RepoOwner = "clutch"
	}
	s.Data.RepoOwner = promptOrDefault("Enter the name of the service owner", s.Data.RepoOwner)
	s.ServiceName = s.Data.ServiceName
	s.Destination = filepath.Join(dest, s.Data.ServiceName)

	registrationData := &RegistrationParams{}

	registrationData.BackendMainPath = filepath.Join(repositoryRoot, "backend", "main.go")
	registrationData.ComponentType = "service"
	registrationData.ComponentName = s.Data.ServiceName
	svcPath := filepath.Join(gitUpstream.RepoProvider, gitUpstream.RepoOwner, gitUpstream.RepoName, "backend/service", s.Data.ServiceName)
	registrationData.ComponentPath = svcPath
	s.RegistrationData = registrationData
}

func (s *ServiceScaffoldWorkflow) GetTemplateValues() interface{} {
	return s.Data
}

func (s *ServiceScaffoldWorkflow) PostProcess(_ *Args, tmpFolder string) {
	destPrefix := filepath.Join(s.Destination, s.ServiceName)
	MoveTempFilesToDest(filepath.Join(tmpFolder, "service.go"), destPrefix+".go")
	MoveTempFilesToDest(filepath.Join(tmpFolder, "service_test.go"), destPrefix+"_test.go")

	RegisterNewComponent(s.RegistrationData)
}
