package main

import (
	"github.com/lyft/clutch/backend/cmd/assets"
	"github.com/lyft/clutch/backend/gateway"
)

func main() {
	flags := gateway.ParseFlags()
	components := gateway.CoreComponentFactory

	gateway.Run(flags, components, assets.VirtualFS)
}
