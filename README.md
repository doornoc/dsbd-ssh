# dsbd-ssh
dsbd-ssh

## protoc
### Install
```
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

### Generate from proto file
```
protoc --go_out=. --go-grpc_out=. *.proto
```