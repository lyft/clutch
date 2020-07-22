// +build ignore

package main

import (
	"github.com/lyft/clutch/backend/gateway"
	amiibomod "github.com/lyft/clutch/backend/module/amiibo"
	amiiboservice "github.com/lyft/clutch/backend/service/amiibo"
)

func main() {
	flags := gateway.ParseFlags()
	components := gateway.CoreComponentFactory
	components.Modules[amiibomod.Name] = amiibomod.New
	components.Services[amiiboservice.Name] = amiiboservice.New
	gateway.Run(flags, components)
}
