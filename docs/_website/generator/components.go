package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var clutchDocRegexp = regexp.MustCompile(`(?s)<!-- START clutchdoc -->(.*)<!-- END clutchdoc -->`)

func clutchDocFromComment(s string) (*ClutchDoc, error) {
	m := clutchDocRegexp.FindStringSubmatch(s)
	if len(m) != 2 {
		return nil, nil
	}

	cd := &ClutchDoc{}
	if err := yaml.Unmarshal([]byte(m[1]), cd); err != nil {
		return nil, err
	}

	return cd, nil
}

func getClutchComponentFromFile(path string) (*Component, bool) {
	if strings.HasSuffix(path, "_test.go") {
		return nil, false
	}

	fileset := token.NewFileSet()
	node, err := parser.ParseFile(fileset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	cd := &ClutchDoc{}
	for _, c := range node.Comments {
		d, err := clutchDocFromComment(c.Text())
		if err != nil {
			log.Fatal(err)
		}
		if d != nil {
			cd = d
			break
		}
	}

	var name *ast.Object
	for k, v := range node.Scope.Objects {
		if strings.HasSuffix(k, "Name") && v.Kind == ast.Con {
			name = v
			break
		}
	}

	if name == nil {
		return nil, false
	}

	nameStr := name.Decl.(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value
	nameStr = strings.Trim(nameStr, "\"")

	relPath := regexp.MustCompile(`^.*?(backend/.*?)$`).FindStringSubmatch(path)[1]
	url := fmt.Sprintf("https://github.com/lyft/clutch/blob/main/%s", relPath)

	t := regexp.MustCompile(`clutch\.(\w+)\.`).FindStringSubmatch(nameStr)
	if len(t) == 0 {
		return nil, false
	}

	component := &Component{
		Name:      nameStr,
		URL:       url,
		Type:      t[1],
		ClutchDoc: cd,
	}
	return component, true
}

type ClutchDoc struct {
	Description string `yaml:"description"`
}

type Component struct {
	Name      string
	Type      string
	URL       string
	ClutchDoc *ClutchDoc
}
