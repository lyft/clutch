package gateway

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/uber-go/tally"
	tallyprom "github.com/uber-go/tally/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/gateway/mux"
	"github.com/lyft/clutch/backend/gateway/stats"
	"github.com/lyft/clutch/backend/middleware"
	"github.com/lyft/clutch/backend/middleware/accesslog"
	"github.com/lyft/clutch/backend/middleware/errorintercept"
	"github.com/lyft/clutch/backend/middleware/timeouts"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/resolver"
	"github.com/lyft/clutch/backend/service"
)

// The purpose of this identifier is to trim the prefix off clutch component names,
// such as services, modules and resolver.
// This is used when constructing the stats namespace for a given component.
const clutchComponentPrefix = "clutch."

// All available components supply their factory here.
// Whether or not they are used at runtime is dependent on the configuration passed to the Gateway.
type ComponentFactory struct {
	Services   service.Factory
	Resolvers  resolver.Factory
	Middleware middleware.Factory
	Modules    module.Factory
}

func Run(f *Flags, cf *ComponentFactory, assets http.FileSystem) {
	cfg := MustReadOrValidateConfig(f)
	RunWithConfig(f, cfg, cf, assets)
}

func RunWithConfig(f *Flags, cfg *gatewayv1.Config, cf *ComponentFactory, assets http.FileSystem) {
	// Init the server's logger.
	logger, err := newLogger(cfg.Gateway.Logger)
	if err != nil {
		newTmpLogger().Fatal("could not instantiate logger", zap.Error(err))
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()

	logger.Info("using configuration", zap.String("file", f.ConfigPath))

	// Init stats.
	var metricsHandler http.Handler
	var scopeOpts tally.ScopeOptions
	switch t := cfg.Gateway.Stats.Reporter.(type) {
	case nil:
		scopeOpts = tally.ScopeOptions{
			Reporter: tally.NullStatsReporter,
		}
	case *gatewayv1.Stats_LogReporter_:
		scopeOpts = tally.ScopeOptions{
			Reporter: stats.NewDebugReporter(logger),
			Prefix:   "clutch",
		}
	case *gatewayv1.Stats_StatsdReporter_:
		reporter, err := stats.NewStatsdReporter(cfg.Gateway.Stats.GetStatsdReporter())
		if err != nil {
			logger.Fatal("error creating statsd reporter", zap.Error(err))
		}
		scopeOpts = tally.ScopeOptions{
			Reporter: reporter,
			Prefix:   "clutch",
		}
	case *gatewayv1.Stats_PrometheusReporter_:
		reporter, err := stats.NewPrometheusReporter(cfg.Gateway.Stats.GetPrometheusReporter())
		if err != nil {
			logger.Fatal("error creating prometheus reporter", zap.Error(err))
		}
		scopeOpts = tally.ScopeOptions{
			CachedReporter:  reporter,
			Prefix:          "clutch",
			SanitizeOptions: &tallyprom.DefaultSanitizerOpts,
		}
		metricsHandler = reporter.HTTPHandler()
	default:
		logger.Fatal("unsupported reporter", zap.Reflect("type", t))
	}

	scope, scopeCloser := tally.NewRootScope(
		scopeOpts,
		cfg.Gateway.Stats.FlushInterval.AsDuration(),
	)
	defer func() {
		if err := scopeCloser.Close(); err != nil {
			panic(err)
		}
	}()

	initScope := scope.SubScope("gateway")
	initScope.Counter("start").Inc(1)

	// Create the error interceptor so services can register error interceptors if desired.
	errorInterceptMiddleware, err := errorintercept.NewMiddleware(nil, logger, initScope)
	if err != nil {
		logger.Fatal("could not create error interceptor middleware", zap.Error(err))
	}

	// Instantiate and register services.
	for _, svcConfig := range cfg.Services {
		factory, ok := cf.Services[svcConfig.Name]
		logger := logger.With(zap.String("serviceName", svcConfig.Name))
		if !ok {
			logger.Fatal("service not found in registry")
		}
		if factory == nil {
			logger.Fatal("service has nil factory")
		}

		if err := validateAny(svcConfig.TypedConfig); err != nil {
			logger.Fatal("service config validation failed", zap.Error(err))
		}

		logger.Info("registering service")
		svc, err := factory(svcConfig.TypedConfig, logger, scope.SubScope(
			strings.TrimPrefix(svcConfig.Name, clutchComponentPrefix)))
		if err != nil {
			logger.Fatal("service instantiation failed", zap.Error(err))
		}
		service.Registry[svcConfig.Name] = svc

		if ei, ok := svc.(errorintercept.Interceptor); ok {
			logger.Info("service registered an error conversion interceptor")
			errorInterceptMiddleware.AddInterceptor(ei.InterceptError)
		}
	}

	for _, resolverCfg := range cfg.Resolvers {
		factory, ok := cf.Resolvers[resolverCfg.Name]
		logger := logger.With(zap.String("resolverName", resolverCfg.Name))
		if !ok {
			logger.Fatal("resolver not found in registry")
		}
		if factory == nil {
			logger.Fatal("resolver has nil factory")
		}

		if err := validateAny(resolverCfg.TypedConfig); err != nil {
			logger.Fatal("resolver config validation failed", zap.Error(err))
		}

		logger.Info("registering resolver")
		res, err := factory(resolverCfg.TypedConfig, logger, scope.SubScope(
			strings.TrimPrefix(resolverCfg.Name, clutchComponentPrefix)))
		if err != nil {
			logger.Fatal("resolver instantiation failed", zap.Error(err))
		}
		resolver.Registry[resolverCfg.Name] = res
	}

	var interceptors []grpc.UnaryServerInterceptor

	// Error interceptors should be first on the stack (last in chain).
	interceptors = append(interceptors, errorInterceptMiddleware.UnaryInterceptor())

	// Access log.
	if cfg.Gateway.Accesslog != nil {
		a, err := accesslog.New(cfg.Gateway.Accesslog, logger, scope)
		if err != nil {
			logger.Fatal("could not create accesslog interceptor", zap.Error(err))
		}
		interceptors = append(interceptors, a.UnaryInterceptor())
	}

	// Timeouts.
	timeoutInterceptor, err := timeouts.New(cfg.Gateway.Timeouts, logger, scope)
	if err != nil {
		logger.Fatal("could not create timeout interceptor", zap.Error(err))
	}
	interceptors = append(interceptors, timeoutInterceptor.UnaryInterceptor())

	// All other configured middleware.
	for _, mCfg := range cfg.Gateway.Middleware {
		logger := logger.With(zap.String("moduleName", mCfg.Name))

		factory, ok := cf.Middleware[mCfg.Name]
		if !ok {
			logger.Fatal("middleware not found in registry")
		}
		if factory == nil {
			logger.Fatal("middleware has nil factory")
		}

		if err := validateAny(mCfg.TypedConfig); err != nil {
			logger.Fatal("middleware config validation failed", zap.Error(err))
		}

		logger.Info("registering middleware")
		m, err := factory(mCfg.TypedConfig, logger, scope)
		if err != nil {
			logger.Fatal("middleware instatiation failed", zap.Error(err))
		}

		interceptors = append(interceptors, m.UnaryInterceptor())
	}

	// Instantiate and register modules listed in the configuration.
	rpcMux, err := mux.New(interceptors, assets, metricsHandler, cfg.Gateway)
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()

	// Create a client connection for the registrar to make grpc-gateway's handlers available.
	// TODO: stand up a private loopback listener for the grpcServer and connect to that instead.
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	if cfg.Gateway.MaxResponseSizeBytes > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(cfg.Gateway.MaxResponseSizeBytes))))
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Gateway.Listener.GetTcp().Address, cfg.Gateway.Listener.GetTcp().Port), opts...)
	if err != nil {
		logger.Fatal("failed to bring up gRPC transport for grpc-gateway handlers", zap.Error(err))
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				logger.Warn("failed to close gRPC transport connection after err", zap.Error(err))
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				logger.Warn("failed to close gRPC transport connection when done", zap.Error(err))
			}
		}()
	}()

	reg := newRegistrar(ctx, rpcMux.JSONGateway, rpcMux.GRPCServer, conn)
	for _, modCfg := range cfg.Modules {
		logger := logger.With(zap.String("moduleName", modCfg.Name))

		factory, ok := cf.Modules[modCfg.Name]
		if !ok {
			logger.Fatal("module not found in registry")
		}
		if factory == nil {
			logger.Fatal("module has nil factory")
		}

		if err := validateAny(modCfg.TypedConfig); err != nil {
			logger.Fatal("module config validation failed", zap.Error(err))
		}

		logger.Info("registering module")
		mod, err := factory(modCfg.TypedConfig, logger, scope.SubScope(
			strings.TrimPrefix(modCfg.Name, clutchComponentPrefix)))
		if err != nil {
			logger.Fatal("module instantiation failed", zap.Error(err))
		}

		if err := mod.Register(reg); err != nil {
			logger.Fatal("registration to gateway failed", zap.Error(err))
		}
	}

	// Now that everything is registered, enable gRPC reflection.
	rpcMux.EnableGRPCReflection()

	// Save metadata on what RPCs being served for fast-lookup by internal services.
	if err := meta.GenerateGRPCMetadata(rpcMux.GRPCServer); err != nil {
		logger.Fatal("reflection on grpc server failed", zap.Error(err))
	}

	// Instantiate server and listen.
	switch t := cfg.Gateway.Listener.Socket.(type) {
	case *gatewayv1.Listener_Tcp:
		// OK
	default:
		logger.Fatal("socket not supported", zap.String("type", fmt.Sprintf("%T", t)))
	}

	if cfg.Gateway.Listener.GetTcp().Secure {
		logger.Fatal("'secure' set to true but listener security is not currently supported")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Gateway.Listener.GetTcp().Address, cfg.Gateway.Listener.GetTcp().Port)
	logger.Info("listening", zap.Namespace("tcp"), zap.String("addr", addr))

	// Figure out the maximum global timeout and set as a backstop (with 1s buffer).
	timeout := computeMaximumTimeout(cfg.Gateway.Timeouts)
	if timeout > 0 {
		timeout += time.Second
	}

	// Start collecting go runtime stats if enabled
	if cfg.Gateway.Stats != nil && cfg.Gateway.Stats.GoRuntimeStats != nil {
		runtimeStats := stats.NewRuntimeStats(scope, cfg.Gateway.Stats.GoRuntimeStats)
		go runtimeStats.Collect(ctx)
	}

	srv := &http.Server{
		Handler:      mux.InsecureHandler(rpcMux),
		Addr:         addr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	logger.Fatal("error bringing up listener", zap.Error(srv.ListenAndServe()))
}
