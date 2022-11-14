//go:build ignore

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/vfsgen"
)

func packageAssets() {
	if len(os.Args) != 2 {
		log.Fatal("usage: go run generate.go <dir>")
	}

	dir := os.Args[1]
	assets := http.Dir(dir)
	err := vfsgen.Generate(assets, vfsgen.Options{
		Filename:     "cmd/assets/generated_assets.go",
		PackageName:  "assets",
		VariableName: "VirtualFS",
		BuildTags:    "withAssets",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	packageAssets()
}
