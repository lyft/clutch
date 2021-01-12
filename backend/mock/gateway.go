package main

import (
	"github.com/lyft/clutch/backend/cmd/assets"
	"github.com/lyft/clutch/backend/gateway"
	"github.com/lyft/clutch/backend/mock/service/auditmock"
	"github.com/lyft/clutch/backend/mock/service/awsmock"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
	"github.com/lyft/clutch/backend/mock/service/envoyadminmock"
	"github.com/lyft/clutch/backend/mock/service/githubmock"
	"github.com/lyft/clutch/backend/mock/service/k8smock"
	"github.com/lyft/clutch/backend/mock/service/topologymock"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/audit"
	"github.com/lyft/clutch/backend/service/aws"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"github.com/lyft/clutch/backend/service/envoyadmin"
	"github.com/lyft/clutch/backend/service/github"
	"github.com/lyft/clutch/backend/service/k8s"
	"github.com/lyft/clutch/backend/service/topology"
)

var MockServiceFactory = service.Factory{
	aws.Name:             awsmock.NewAsService,
	audit.Name:           auditmock.NewAsService,
	envoyadmin.Name:      envoyadminmock.NewAsService,
	experimentstore.Name: experimentstoremock.NewMock,
	github.Name:          githubmock.NewAsService,
	k8s.Name:             k8smock.NewAsService,
	topology.Name:        topologymock.NewAsService,
}

func main() {
	cf := gateway.CoreComponentFactory

	// Replace core services with any available mocks.
	cf.Services = MockServiceFactory

	gateway.Run(gateway.ParseFlags(), cf, assets.VirtualFS)
}
