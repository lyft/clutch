package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type WorkflowPackage struct {
	PackageName string
	Description string
	URL         string
	Workflows   []string
}

type packageJSON struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var routes = regexp.MustCompile(`^\s+routes\W`)
var displayName = regexp.MustCompile(`displayName:\s*(.*?)$`)

func getWorkflowMetadata(indexPath string) []string {
	fh, err := os.Open(indexPath)
	if err != nil {
		return nil
	}
	defer fh.Close()

	var ret []string

	s := bufio.NewScanner(fh)
	atRoutes := false
	for s.Scan() {
		t := s.Text()
		if m := routes.FindString(t); m != "" {
			atRoutes = true
		}

		if atRoutes {
			if m := displayName.FindStringSubmatch(t); len(m) == 2 {
				ret = append(ret, strings.Trim(m[1], ",\" "))
			}
		}
	}
	return ret
}

func getWorkflowPackage(pathToPackageJSON string) *WorkflowPackage {
	raw, err := ioutil.ReadFile(pathToPackageJSON)
	if err != nil {
		log.Fatal(err)
	}

	pj := &packageJSON{}
	if err := json.Unmarshal(raw, &pj); err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(pathToPackageJSON)

	indexPath := filepath.Join(dir, "src/index.jsx")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		indexPath = filepath.Join(dir, "src/index.tsx")
	}

	return &WorkflowPackage{
		PackageName: pj.Name,
		Description: pj.Description,
		URL:         fmt.Sprintf("https://github.com/lyft/clutch/blob/main/frontend/workflows/%s", filepath.Base(dir)),
		Workflows:   getWorkflowMetadata(indexPath),
	}
}
