package scaffold

import (
	"flag"
)

const yarnInstallVersion = "4.3.1"

type Args struct {
	Internal          bool
	Mode              string
	GoPin             string
	Org               string
	TemplateOverwrite string
	YarnPin           string
}

func ParseArgs() *Args {
	f := &Args{}
	flag.BoolVar(&f.Internal, "i", false, "when creating workflows in clutch")
	flag.StringVar(&f.Mode, "m", "gateway", "oneof gateway, frontend-plugin")
	flag.StringVar(&f.GoPin, "p", "main", "sha or other github ref to version of tools used in scaffolding")
	flag.StringVar(&f.Org, "o", "lyft", "overrides the github organization (for use in fork testing)")
	flag.StringVar(&f.TemplateOverwrite, "templates", "", "directory to use for template overwrites")
	flag.StringVar(&f.YarnPin, "y", "4.3.1", "version of yarn to use")
	flag.Parse()
	return f
}
