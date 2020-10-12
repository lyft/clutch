module github.com/lyft/clutch/backend

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-sdk-go v1.35.2
	github.com/bufbuild/buf v0.25.0
	github.com/cactus/go-statsd-client/statsd v0.0.0-20200623234511-94959e3146b2
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/go-control-plane v0.9.7
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/fullstorydev/grpcurl v1.7.0
	github.com/go-git/go-billy/v5 v5.0.0
	github.com/go-git/go-git/v5 v5.2.0
	github.com/go-swagger/go-swagger v0.23.0
	github.com/gobwas/glob v0.2.3
	github.com/golang-migrate/migrate/v4 v4.13.0
	github.com/golang/protobuf v1.4.2
	github.com/google/go-github/v32 v32.1.0
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/grpc-gateway v1.15.0
	github.com/iancoleman/strcase v0.1.2
	github.com/jhump/protoreflect v1.7.1-0.20200723220026-11eaaf73e0ec
	github.com/lib/pq v1.8.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/shurcooL/githubv4 v0.0.0-20200928013246-d292edc3691b
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546
	github.com/slack-go/slack v0.7.1
	github.com/stretchr/testify v1.6.1
	github.com/uber-go/tally v3.3.17+incompatible
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sync v0.0.0-20201008141435-b3e1573b7520
	google.golang.org/genproto v0.0.0-20201009135657-4d944d34d83c
	google.golang.org/grpc v1.33.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
)
