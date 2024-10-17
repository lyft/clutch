package scaffold

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

// Replace any special characters, that might appear in repository owners or organizations names, but are not supported
// in proto3 package names.
var sanitizeProto3Identifier = strings.NewReplacer(
	"-", "_",
	"/", "_",
)

// SanitizeProto3 identifiers by replacing unsupported characters.
func (gatewayTemplateValues) SanitizeProto3(s string) string {
	return sanitizeProto3Identifier.Replace(s)
}

func promptOrDefault(prompt string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	if defaultValue == "" {
		fmt.Printf("%s: ", prompt)
	} else {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	}
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	text = strings.TrimSpace(text)

	if text == "" && defaultValue != "" {
		return defaultValue
	}
	if text == "" {
		log.Fatal("no input provided")
	}

	return text
}

func TemplateFiles(templateRoot string, templateOut string, data interface{}) error {
	return filepath.Walk(templateRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relpath, err := filepath.Rel(templateRoot, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			if relpath != "." {
				err := os.MkdirAll(filepath.Join(templateOut, relpath), 0755)
				return err
			}
			return nil
		}

		log.Println(relpath)

		t, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		out := strings.TrimSuffix(filepath.Join(templateOut, relpath), ".tmpl")
		fh, err := os.Create(out)
		if err != nil {
			return err
		}

		return t.Execute(fh, data)
	})
}

func FixDestDir(dest string) (string, error) {
	// Fix tilde HOME path.
	if strings.HasPrefix(dest, "~/") {
		homeUser, err := user.Current()
		if err != nil {
			log.Fatal("could not get user's information", err)
		}
		dest = filepath.Join(homeUser.HomeDir, dest[2:])
	}

	// Check if dest exists.
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		log.Fatal("ERROR destination folder exists")
	}

	return dest, nil
}

func TemplateTempDir(templateRoot string, data interface{}) string {
	tmpout, err := os.MkdirTemp(os.TempDir(), "clutch-scaffolding-")
	if err != nil {
		log.Fatal("could not create temp dir", err)
	}
	log.Println("Using tmpdir", tmpout)

	// Walk files and template them.
	err = TemplateFiles(templateRoot, tmpout, data)
	if err != nil {
		log.Fatal(err)
	}

	return tmpout
}

func MoveTempFilesToDest(tmpout string, dest string) {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.Rename(tmpout, dest); err != nil {
		if os.IsExist(err) {
			log.Fatal(fmt.Sprintf("Failed moving %s to %s destination folder already exists", tmpout, dest))
		}
		log.Fatal(err)
	}
}

func YarnInstall(yarn string, pinned string) {
	installVersion := pinned
	if len(installVersion) == 0 {
		installVersion = yarnInstallVersion
	}

	yarnVersionCmd := exec.Command(yarn, "--version")
	if out, err := yarnVersionCmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`yarn --version` returned the above error")
	} else if strings.TrimRight(string(out), "\n") != installVersion {
		log.Println("Yarn version is not equal to", installVersion)
		corepack := promptOrDefault("Would you like to attempt installation of Yarn via corepack?", "Y/n")
		if !strings.HasPrefix(strings.ToLower(corepack), "y") {
			log.Fatal("Corepack installation not granted, exiting")
		} else {
			corepackCmd := exec.Command("corepack", "enable")
			if out, err := corepackCmd.CombinedOutput(); err != nil {
				fmt.Println(string(out))
				log.Fatal("`corepack enable` returned the above error. Please make sure you are on Node > 18")
			}

			corepackPrepareCmd := exec.Command("corepack", "prepare", "yarn@"+installVersion, "--activate")
			if out, err := corepackPrepareCmd.CombinedOutput(); err != nil {
				fmt.Println(string(out))
				log.Fatal("`corepack prepare", "yarn@"+installVersion, "--activate` returned the above error. Please correct the error and attempt again.")
			}
		}
	}
}
