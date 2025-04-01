package scaffold

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

type GitUpstream struct {
	RepoProvider string
	RepoOwner    string
	RepoName     string
}

func determineGitUpstream(root string) *GitUpstream {
	oldDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(root)
	if err != nil {
		log.Fatal(err)
	}

	url, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		log.Fatal(err)
	}
	urlStr := strings.TrimSpace(string(url))
	// This assumes SSH format
	r := regexp.MustCompile(`^git@(?<RepoProvider>[A-Za-z.]+):(?<RepoOwner>.+)/(?<RepoName>.+)\.git$`)
	matches := r.FindStringSubmatch(urlStr)
	if r == nil || len(matches) != 4 {
		log.Fatal(fmt.Errorf("unable to determine git upstream from %s", urlStr))
	}

	err = os.Chdir(oldDir)
	if err != nil {
		log.Fatal(err)
	}

	return &GitUpstream{
		RepoProvider: matches[1],
		RepoOwner:    matches[2],
		RepoName:     matches[3],
	}
}

func determineGoPackage(workdir string) string {
	oldDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(workdir)
	if err != nil {
		log.Fatal(err)
	}

	goPackageOut, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(oldDir)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(goPackageOut))
}

func determineGoPath() string {
	goPathOut, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(goPathOut))
}

func DetermineYarnPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	yarnScript := filepath.Join(path, "../..", "build", "bin", "yarn.sh")
	scriptPath, err := exec.LookPath(yarnScript)
	if err != nil {
		cmd := exec.Command("make", "yarn-ensure")
		cmd.Dir = filepath.Join(path, "../..")
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Println(string(out))
			log.Fatal("`make yarn-ensure` returned the above error")
		}

		return yarnScript
	}
	return scriptPath
}

func determineUsername() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return u.Username
}

func determineUserEmail() string {
	gitEmail, err := exec.Command("git", "config", "user.email").Output()
	if err != nil {
		log.Fatal(err)
	}
	email := strings.TrimSpace(string(gitEmail))
	if email == "" {
		email = "unknown@example.com"
	}
	return email
}
