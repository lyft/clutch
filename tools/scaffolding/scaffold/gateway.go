package scaffold

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

const (
	defaultSourceControlProvider = "github.com"
	defaultRepoName              = "clutch-custom-gateway"
)

type gatewayTemplateValues struct {
	RepoOwner    string
	RepoName     string
	RepoProvider string
}

func GetGatewayTemplateValues() (*gatewayTemplateValues, string) {
	// Ask the user if assumptions are correct or a new destination is needed.
	log.Println("Welcome!")
	fmt.Println("*** Analyzing environment...")

	gopath := determineGoPath()
	username := determineUsername()

	fmt.Println("> GOPATH:", gopath)
	fmt.Println("> User:", username)

	dest := filepath.Join(gopath, "src", defaultSourceControlProvider, username, defaultRepoName)
	fmt.Println("\n*** Based on your environment, we've picked the following destination for your new repo:")
	fmt.Println(">", dest)
	fmt.Println("\nNote: please pay special attention to see if the username matches your provider's username.")
	okay := promptOrDefault("Is this okay?", "Y/n")
	data := &gatewayTemplateValues{
		RepoOwner:    username,
		RepoName:     defaultRepoName,
		RepoProvider: defaultSourceControlProvider,
	}
	if !strings.HasPrefix(strings.ToLower(okay), "y") {
		data.RepoProvider = promptOrDefault("Enter the name of the source control provider", data.RepoProvider)
		data.RepoOwner = promptOrDefault("Enter the name of the repository owner or org", data.RepoOwner)
		data.RepoName = promptOrDefault("Enter the desired repository name", data.RepoName)
		dest = promptOrDefault(
			"Enter the destination folder",
			filepath.Join(gopath, "src", data.RepoProvider, data.RepoOwner, data.RepoName),
		)
	}

	return data, dest
}

func PostProcessGateway(flags *Args, tmpFolder string, dest string) {
	GenerateAPI(flags, tmpFolder, dest)

	MoveTempFilesToDest(tmpFolder, dest)
}
