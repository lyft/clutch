package main

import (
	"github.com/lyft/clutch/backend/cmd/assets"
	"github.com/lyft/clutch/backend/gateway"
	dynamodbmod "github.com/lyft/clutch/backend/module/dynamodb"
)

func main() {
	flags := gateway.ParseFlags()
	components := gateway.CoreComponentFactory
	components.Modules[dynamodbmod.Name] = dynamodbmod.New

	gateway.Run(flags, components, assets.VirtualFS)
}
