module github.com/MetaRPC/GoMT5/examples

go 1.23.6

replace github.com/MetaRPC/GoMT5/examples/errors => ./errors

require (
	github.com/MetaRPC/GoMT5/package v0.0.0-00010101000000-000000000000
	github.com/MetaRPC/GoMT5/examples/errors v0.0.0-00010101000000-000000000000
	github.com/MetaRPC/GoMT5/mt5 v0.0.0
	github.com/google/uuid v1.6.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250811230008-5f3141c8851a
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)

replace github.com/MetaRPC/GoMT5/mt5 => ./mt5

replace github.com/MetaRPC/GoMT5/package => ../package
