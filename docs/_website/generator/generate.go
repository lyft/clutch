package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

const repoRoot = "../../../"

const editBasePath = "https://github.com/lyft/clutch/edit/main/docs/"

type templateData struct {
	Middleware []*Component
	Modules    []*Component
	Resolvers  []*Component
	Services   []*Component

	WorkflowPackages []*WorkflowPackage

	EditURL string
}

var funcMap = sprig.TxtFuncMap()

func generate(source, dest string, td *templateData) error {
	t := template.New(filepath.Base(source))
	t = t.Funcs(funcMap)

	td.EditURL = fmt.Sprintf("custom_edit_url: %s", strings.Replace(source, "../../", editBasePath, 1))

	if _, err := t.ParseFiles(source); err != nil {
		return err
	}

	fh, err := os.Create(dest)
	if err != nil {
		return err
	}

	if err := t.Execute(fh, td); err != nil {
		return err
	}

	return nil
}

func getFiles(root, extension string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relpath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if !strings.HasSuffix(relpath, extension) || strings.HasPrefix(relpath, "_website") || strings.HasPrefix(relpath, "README") || info.IsDir() {
			// Exclude directories or website dir.
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func getTemplateData() *templateData {
	td := &templateData{}

	// Backend.

	files, err := getFiles(filepath.Join(repoRoot, "backend"), ".go")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if c, ok := getClutchComponentFromFile(file); ok {
			switch c.Type {
			case "service":
				td.Services = append(td.Services, c)
			case "module":
				td.Modules = append(td.Modules, c)
			case "middleware":
				td.Middleware = append(td.Middleware, c)
			case "resolver":
				td.Resolvers = append(td.Resolvers, c)
			}
		}
	}

	// Frontend.
	files, err = getFiles(filepath.Join(repoRoot, "frontend/workflows"), "package.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file, "node_modules") {
			// ignore sub packages from builds.
			continue
		}
		w := getWorkflowPackage(file)
		td.WorkflowPackages = append(td.WorkflowPackages, w)
	}

	return td
}

func main() {

	docsRoot := "../../"
	destRoot := filepath.Join(docsRoot, "_website/generated/docs")

	// Delete any old files.
	subdirs, _ := ioutil.ReadDir(destRoot)
	for _, f := range subdirs {
		if err := os.RemoveAll(filepath.Join(destRoot, f.Name())); err != nil {
			log.Fatal(err)
		}
	}

	// Parse protos and register funcs to map.
	ps, err := newProtoScope(filepath.Join(repoRoot, "api"))
	if err != nil {
		log.Fatal(err)
	}
	funcMap["simpleProtoYAML"] = ps.getSimpleMessageYAML

	// Get template data.
	td := getTemplateData()

	// Get all files to be interpolated and generated.
	files, err := getFiles(docsRoot, ".md")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		relpath, _ := filepath.Rel(docsRoot, f)
		dest := filepath.Join(destRoot, relpath)

		// Make directory if it doesn't exist.
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			log.Fatal(err)
		}

		fmt.Println(f)
		if err := generate(f, dest, td); err != nil {
			log.Fatal(err)
		}
	}
}
