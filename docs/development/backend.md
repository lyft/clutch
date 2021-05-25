---
title: Backend
{{ .EditURL }}
---

Familiarizing yourself with concepts in the [Architecture Reference](/docs/about/architecture) is recommended prior to starting development on the backend.

## Structure

```
backend
├─ api         # generated code from proto
├─ gateway     # JSON/gRPC server for components
├─ middleware  # request/response interceptors
├─ mock        # mock gateway and components
├─ module      # implementations of proto stubs for API endpoints
├─ resolver    # resolver implementations
└─ service     # interfaces to external services or resources
```

## Components

Components consist primarily of a `Name` string for registration purposes and a `New` factory method, which returns an implementation of the interface required for that component.

### New

All `New` methods for components have the following signature:

```go
New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (<component type>, error)
```

Configuration must be unmarshaled from the Any type into the native protobuf type for further consumption using [`anypb`](https://pkg.go.dev/google.golang.org/protobuf/types/known/anypb#hdr-Unmarshaling_an_Any).

References to the logger and scope can be saved in the component implementation's `struct` for later use. Note that the gateway adds a component namespace for the logger and stats scope so it does not have to be added by the developer.

### Gateway and Middleware

- [`backend/gateway/`](https://github.com/lyft/clutch/tree/main/backend/gateway)
- [`backend/middleware/`](https://github.com/lyft/clutch/tree/main/backend/middleware)

The gateway is extensible via middleware, which operates on the incoming request and the outgoing response in a single context.

Currently the gateway only supports unary requests (as opposed to streaming). Therefore middleware returns a [server unary interceptor](https://github.com/grpc/grpc-go/tree/master/examples/features/interceptor#unary-interceptor-1).

### Modules

- [`backend/module/`](https://github.com/lyft/clutch/tree/main/backend/module)

Modules are implementations of API stubs generated from proto. The `Register` method is a callback used by the gateway to add the gRPC and JSON servers to the main gateway server.

### Resolvers

- [`backend/resolver/`](https://github.com/lyft/clutch/tree/main/backend/resolver)

Resolvers abstract resource location into an extensible framework. Finding resources is different at almost every organization. Changing the first step of locating a resource should not require wholesale replacement of a feature. That's where resolvers come in.

Resolvers start in protobuf with an object definition representing the user input needed to locate a resource. The fields are additionally annotated with schema information (such as `placeholder`, `required`, `display_name`) at the object and field level so that the frontend knows how to render the appropriate form.

The resolver has two ways of locating resources, `Search` takes a free-form text query and locates resources by interpreting the query and making the appropriate call. `Resolve` takes the filled-out schemas and makes a query based on that structured information.

One additional nicety in the resolver is the [`FanoutHandler`](https://github.com/lyft/clutch/blob/main/backend/resolver/fanouthandler.go) which makes managing concurrent requests easier. It has `limit` handling to return early once the limit is satisfied, cancelling any remaining outstanding requests.

In the future, the resolver will handle autocomplete and asynchronous validation of form input for the frontend.

More docs are coming on developing resolvers. For now look at other resolvers as an example.

### Services

- [`backend/service/`](https://github.com/lyft/clutch/tree/main/backend/service)

Interactions with third-party services or networked dependencies such as a database should always hide behind a service. Services have no other restrictions and the underlying type is simply an empty `interface`.

All services are placed in a registry for use by other components. Other components that depend on services fetch them in their factory method. For example,

```go
// Create a new instance of my module that uses the foo service.
func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	si, ok := service.Registry["clutch.service.foo"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	c, ok := si.(foo.Client)
	if !ok {
		return nil, errors.New("service was not the correct type")
    }

    return &mod{client: c}, nil
}
```

#### Custom Code

Because services return interfaces, it is easy to substitute in a custom implementation by creating a new custom component that implements the interface type.

It is also possible to override specific methods of a service.
