module github.com/dolthub/go-mysql-server

require (
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/dolthub/sqllogictest/go v0.0.0-20201107003712-816f3ae12d81
	github.com/dolthub/vitess v0.0.0-20221121184553-8d519d0bbb91
	github.com/go-kit/kit v0.13.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gocraft/dbr/v2 v2.7.2
	github.com/google/flatbuffers v24.3.25+incompatible
	github.com/google/uuid v1.6.0
	github.com/hashicorp/golang-lru v1.0.2
	github.com/lestrrat-go/strftime v1.1.0
	github.com/mitchellh/hashstructure v1.1.0
	github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/shopspring/decimal v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel v1.29.0
	go.opentelemetry.io/otel/trace v1.29.0
	golang.org/x/sync v0.8.0
	golang.org/x/text v0.17.0
	golang.org/x/tools v0.24.0
	gopkg.in/src-d/go-errors.v1 v1.0.0
)

require (
	golang.org/x/crypto v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240827150818-7e3bb234dfed // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	google.golang.org/grpc v1.66.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/oliveagle/jsonpath => github.com/dolthub/jsonpath v0.0.0-20210609232853-d49537a30474

go 1.21

toolchain go1.21.11
