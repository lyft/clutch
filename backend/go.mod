module github.com/lyft/clutch/backend

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/squirrel v1.5.0
	github.com/aws/aws-sdk-go-v2 v1.3.2
	github.com/aws/aws-sdk-go-v2/config v1.1.5
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.3.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.4.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.2.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.5.0
	github.com/aws/smithy-go v1.3.1
	github.com/bufbuild/buf v0.37.0
	github.com/cactus/go-statsd-client/statsd v0.0.0-20200623234511-94959e3146b2
	github.com/coreos/go-oidc/v3 v3.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/go-control-plane v0.9.9-0.20210304204206-e2b50f82e48e
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/fullstorydev/grpcurl v1.8.0
	github.com/go-git/go-billy/v5 v5.1.0
	github.com/go-git/go-git/v5 v5.3.0
	github.com/gobwas/glob v0.2.3
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/golang/protobuf v1.5.2
	github.com/google/go-github/v35 v35.1.0
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.3.0
	github.com/iancoleman/strcase v0.1.3
	github.com/jhump/protoreflect v1.8.1
	github.com/lib/pq v1.10.1
	github.com/mitchellh/hashstructure/v2 v2.0.1
	github.com/shurcooL/githubv4 v0.0.0-20201206200315-234843c633fa
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546
	github.com/slack-go/slack v0.9.0
	github.com/stretchr/testify v1.7.0
	github.com/uber-go/tally v3.3.17+incompatible
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210415231046-e915ea6b2b7d
	golang.org/x/oauth2 v0.0.0-20210413134643-5e61552d6c78
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/genproto eac7f76ac494
	google.golang.org/grpc v1.37.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/square/go-jose.v2 v2.5.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
)
