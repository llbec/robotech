module github.com/llbec/robotech

go 1.16

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/dfuse-io/client-go v0.0.0-20210526205821-9a3731282240
	github.com/dfuse-io/dgrpc v0.0.0-20210424033943-10e04dd5b19c
	github.com/dfuse-io/pbgo v0.0.6-0.20210429181308-d54fc7723ad3
	github.com/ethereum/go-ethereum v1.10.3
	github.com/fsnotify/fsnotify v1.4.9
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/websocket v1.4.2
	github.com/spf13/viper v1.7.1
	github.com/streamingfast/streamingfast-client v0.0.7
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.35.0
)

replace github.com/gorilla/rpc => github.com/dfuse-io/rpc v1.2.1-0.20201124195002-f9fc01524e38

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

replace github.com/dfuse-io/pbgo => github.com/dfuse-io/pbgo v0.0.6-0.20210108203525-2912af560a1f
