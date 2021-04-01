package gateway

import (
	"github.com/lyft/clutch/backend/middleware"
	"github.com/lyft/clutch/backend/middleware/audit"
	"github.com/lyft/clutch/backend/middleware/authn"
	"github.com/lyft/clutch/backend/middleware/authz"
	"github.com/lyft/clutch/backend/middleware/stats"
	"github.com/lyft/clutch/backend/middleware/validate"
	"github.com/lyft/clutch/backend/module"
	assetsmod "github.com/lyft/clutch/backend/module/assets"
	auditmod "github.com/lyft/clutch/backend/module/audit"
	authnmod "github.com/lyft/clutch/backend/module/authn"
	authzmod "github.com/lyft/clutch/backend/module/authz"
	awsmod "github.com/lyft/clutch/backend/module/aws"
	experimentationapi "github.com/lyft/clutch/backend/module/chaos/experimentation/api"
	"github.com/lyft/clutch/backend/module/chaos/redisexperimentation"
	"github.com/lyft/clutch/backend/module/chaos/serverexperimentation"
	xdsmod "github.com/lyft/clutch/backend/module/chaos/serverexperimentation/xds"
	"github.com/lyft/clutch/backend/module/envoytriage"
	"github.com/lyft/clutch/backend/module/healthcheck"
	k8smod "github.com/lyft/clutch/backend/module/k8s"
	kinesismod "github.com/lyft/clutch/backend/module/kinesis"
	resolvermod "github.com/lyft/clutch/backend/module/resolver"
	"github.com/lyft/clutch/backend/module/sourcecontrol"
	topologymod "github.com/lyft/clutch/backend/module/topology"
	"github.com/lyft/clutch/backend/resolver"
	awsresolver "github.com/lyft/clutch/backend/resolver/aws"
	coreresolver "github.com/lyft/clutch/backend/resolver/core"
	k8sresolver "github.com/lyft/clutch/backend/resolver/k8s"
	"github.com/lyft/clutch/backend/service"
	auditservice "github.com/lyft/clutch/backend/service/audit"
	loggingsink "github.com/lyft/clutch/backend/service/auditsink/logger"
	"github.com/lyft/clutch/backend/service/auditsink/slack"
	authnservice "github.com/lyft/clutch/backend/service/authn"
	authzservice "github.com/lyft/clutch/backend/service/authz"
	awsservice "github.com/lyft/clutch/backend/service/aws"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/terminator"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
	"github.com/lyft/clutch/backend/service/envoyadmin"
	"github.com/lyft/clutch/backend/service/github"
	k8sservice "github.com/lyft/clutch/backend/service/k8s"
	topologyservice "github.com/lyft/clutch/backend/service/topology"
)

var Middleware = middleware.Factory{
	audit.Name:    audit.New,
	authn.Name:    authn.New,
	authz.Name:    authz.New,
	stats.Name:    stats.New,
	validate.Name: validate.New,
}

var Modules = module.Factory{
	assetsmod.Name:             assetsmod.New,
	auditmod.Name:              auditmod.New,
	authnmod.Name:              authnmod.New,
	authzmod.Name:              authzmod.New,
	awsmod.Name:                awsmod.New,
	envoytriage.Name:           envoytriage.New,
	experimentationapi.Name:    experimentationapi.New,
	k8smod.Name:                k8smod.New,
	kinesismod.Name:            kinesismod.New,
	healthcheck.Name:           healthcheck.New,
	resolvermod.Name:           resolvermod.New,
	xdsmod.Name:                xdsmod.New,
	serverexperimentation.Name: serverexperimentation.New,
	redisexperimentation.Name:  redisexperimentation.New,
	sourcecontrol.Name:         sourcecontrol.New,
	topologymod.Name:           topologymod.New,
}

var Services = service.Factory{
	auditservice.Name:        auditservice.New,
	authnservice.Name:        authnservice.New,
	authnservice.StorageName: authnservice.NewStorage,
	authzservice.Name:        authzservice.New,
	awsservice.Name:          awsservice.New,
	envoyadmin.Name:          envoyadmin.New,
	experimentstore.Name:     experimentstore.New,
	terminator.Name:          terminator.New,
	github.Name:              github.New,
	k8sservice.Name:          k8sservice.New,
	loggingsink.Name:         loggingsink.New,
	pgservice.Name:           pgservice.New,
	slack.Name:               slack.New,
	topologyservice.Name:     topologyservice.New,
}

var Resolvers = resolver.Factory{
	awsresolver.Name:  awsresolver.New,
	k8sresolver.Name:  k8sresolver.New,
	coreresolver.Name: coreresolver.New,
}

var CoreComponentFactory = &ComponentFactory{
	Services:   Services,
	Resolvers:  Resolvers,
	Middleware: Middleware,
	Modules:    Modules,
}
