package main

import (
	"github.com/lyft/clutch/backend/gateway"

	"{{ .RepoProvider}}/{{ .RepoOwner }}/{{ .RepoName }}/backend/module/echo"
)

func main() {
	flags := gateway.ParseFlags()

	components := gateway.CoreComponentFactory

	// Add custom components.
	components.Modules[echo.Name] = echo.New

	gateway.Run(flags, components)
}
