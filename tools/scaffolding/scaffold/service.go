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
	ModuleName  string
	RepoOwner   string
}

func GetServiceTemplateValues() (*serviceTemplateValues, string, *RegistrationParams) {
	log.Println("Welcome!")
	fmt.Println("*** Analyzing environment...")

	dest := filepath.Join(os.Getenv("OLDPWD"), "backend", "service")

	fmt.Println("\n*** Based on your environment, we've picked the following destination for your new service:")
	fmt.Println(">", dest)
	okay := promptOrDefault("Is this okay?", "Y/n")
	if !strings.HasPrefix(strings.ToLower(okay), "y") {
		dest = promptOrDefault("Enter the destination folder", dest)
	}

	data := &serviceTemplateValues{}

	data.ServiceName = strings.ToLower(promptOrDefault("Enter the name of this service", "helloworld"))
	description := promptOrDefault("Enter a description of the service", "Greet the world")
	data.Description = strings.ToUpper(description[:1]) + description[1:]
	gitUpstream := determineGitUpstream()
	data.RepoOwner = promptOrDefault("Enter the name of the organization", gitUpstream.RepoOwner)

	destPrefix := filepath.Join(dest, data.ServiceName, data.ServiceName)

	registrationData := &RegistrationParams{}

	registrationData.BackendMainPath = filepath.Join(os.Getenv("OLDPWD"), "backend", "main.go")
	registrationData.ComponentType = "service"
	registrationData.ComponentName = data.ServiceName
	svcPath := filepath.Join(gitUpstream.RepoProvider, gitUpstream.RepoOwner, gitUpstream.RepoName, "backend/service", data.ServiceName)
	registrationData.ComponentPath = svcPath

	return data, destPrefix, registrationData
}

func PostProcessService(flags *Args, tmpFolder string, destPrefix string) {
	MoveTempFilesToDest(filepath.Join(tmpFolder, "service.go"), destPrefix+".go")
	MoveTempFilesToDest(filepath.Join(tmpFolder, "service_test.go"), destPrefix+"_test.go")
}
