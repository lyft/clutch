// +build tools

// This package tracks build dependencies so they are not removed when `go mod tidy` is run.
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/lyft/clutch/tools"
)
