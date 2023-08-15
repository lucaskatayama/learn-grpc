# Quickstart Helloworld

A simple client/server application with gRPC

## Protobuf

Install `protoc`

```sh
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Compile `.proto`

```sh
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto
```


## References

- [HelloWorld from Github gRPC repository](https://github.com/grpc/grpc-go/tree/master/examples/helloworld)
