module github.com/MetaRPC/GoMT5/examples

go 1.23.6

require (
	github.com/MetaRPC/GoMT5/mt5 v0.0.0
	github.com/google/uuid v1.6.0
)

require (
	git.mtapi.io/root/mrpc-proto v0.0.0-20250812093834-58b4119a2c55 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250811230008-5f3141c8851a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/grpc v1.74.2 // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)

replace github.com/MetaRPC/GoMT5/mt5 => ./mt5
