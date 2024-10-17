package scaffold

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

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
