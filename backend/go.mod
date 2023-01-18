module github.com/lyft/clutch/backend

go 1.19

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/squirrel v1.5.3
	github.com/aws/aws-sdk-go-v2 v1.16.16
	github.com/aws/aws-sdk-go-v2/config v1.17.8
	github.com/aws/aws-sdk-go-v2/credentials v1.12.21
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.23.16
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.17.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.61.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.18.19
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.15.19
	github.com/aws/aws-sdk-go-v2/service/s3 v1.27.11
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.19
	github.com/aws/smithy-go v1.13.3
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/bufbuild/buf v0.56.0
	github.com/cactus/go-statsd-client/statsd v0.0.0-20200623234511-94959e3146b2
	github.com/coreos/go-oidc/v3 v3.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/go-control-plane v0.10.3-0.20220808135628-785751892618
	github.com/envoyproxy/protoc-gen-validate v0.9.1
	github.com/fullstorydev/grpcurl v1.8.7
	github.com/go-git/go-billy/v5 v5.3.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/gobwas/glob v0.2.3
	github.com/gogo/status v1.1.1
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/golang/protobuf v1.5.2
	github.com/google/go-github/v37 v37.0.0
	github.com/google/go-github/v46 v46.0.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.12.0
	github.com/iancoleman/strcase v0.2.0
	github.com/jhump/protoreflect v1.14.0
	github.com/joho/godotenv v1.4.0
	github.com/lib/pq v1.10.7
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/shurcooL/githubv4 v0.0.0-20221229060216-a8d4a561cc93
	github.com/shurcooL/graphql v0.0.0-20220606043923-3cf50f8a0a29
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546
	github.com/slack-go/slack v0.12.1
	github.com/stretchr/testify v1.8.1
	github.com/uber-go/tally/v4 v4.1.4
	go.temporal.io/sdk v1.19.0
	go.temporal.io/sdk/contrib/tally v0.2.0
	go.uber.org/zap v1.24.0
	golang.org/x/net v0.4.0
	golang.org/x/oauth2 v0.3.0
	golang.org/x/sync v0.1.0
	google.golang.org/genproto v0.0.0-20221227171554-f9683d7f8bef
	google.golang.org/grpc v1.51.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.23.6
	k8s.io/apimachinery v0.23.6
	k8s.io/client-go v0.23.6
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448
)

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.6 // indirect
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cncf/xds/go v0.0.0-20220314180256-7f1daf1720fc // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-github/v29 v29.0.3 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/lyft/protoc-gen-star v0.6.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/twmb/murmur3 v1.1.5 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.temporal.io/api v1.13.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65 // indirect
	golang.org/x/tools v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65 // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
