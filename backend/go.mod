module github.com/lyft/clutch/backend

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/squirrel v1.5.1
	github.com/aws/aws-sdk-go-v2 v1.11.1
	github.com/aws/aws-sdk-go-v2/config v1.10.1
	github.com/aws/aws-sdk-go-v2/credentials v1.6.1
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.15.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.8.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.22.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.13.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.9.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.19.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.10.0
	github.com/aws/smithy-go v1.9.0
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/bufbuild/buf v0.56.0
	github.com/cactus/go-statsd-client/statsd v0.0.0-20200623234511-94959e3146b2
	github.com/coreos/go-oidc/v3 v3.1.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/go-control-plane v0.10.1
	github.com/envoyproxy/protoc-gen-validate v0.6.2
	github.com/fullstorydev/grpcurl v1.8.5
	github.com/go-git/go-billy/v5 v5.3.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/gobwas/glob v0.2.3
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/golang/protobuf v1.5.2
	github.com/google/go-github/v29 v29.0.3 // indirect
	github.com/google/go-github/v37 v37.0.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/iancoleman/strcase v0.2.0
	github.com/jhump/protoreflect v1.10.1
	github.com/joho/godotenv v1.4.0
	github.com/lib/pq v1.10.4
	github.com/m3db/prometheus_client_golang v0.8.1 // indirect
	github.com/m3db/prometheus_client_model v0.1.0 // indirect
	github.com/m3db/prometheus_common v0.1.0 // indirect
	github.com/m3db/prometheus_procfs v0.8.1 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/shurcooL/githubv4 v0.0.0-20211117020012-5800b9de5b8b
	github.com/shurcooL/graphql v0.0.0-20200928012149-18c5c3165e3a
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546
	github.com/slack-go/slack v0.10.0
	github.com/stretchr/testify v1.7.0
	github.com/twmb/murmur3 v1.1.5 // indirect
	github.com/uber-go/tally v3.4.2+incompatible
	go.uber.org/zap v1.19.1
	golang.org/x/net v0.0.0-20211013171255-e13a2654a71e
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/genproto v0.0.0-20211013025323-ce878158c4d4
	google.golang.org/grpc v1.42.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.22.3
	k8s.io/apimachinery v0.22.3
	k8s.io/client-go v0.22.3
)
