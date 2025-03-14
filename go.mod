module github.com/kenliu/peer-acks-v2

go 1.21
toolchain go1.24.1

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.9.1
	github.com/lib/pq v1.10.9
	github.com/slack-go/slack v0.11.2
	github.com/stretchr/testify v1.9.0
)

require (
	cloud.google.com/go/functions v1.19.3 // indirect
	github.com/cloudevents/sdk-go/v2 v2.15.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Force consistent version
replace github.com/ugorji/go => github.com/ugorji/go v1.2.12
