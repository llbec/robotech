# install protoc

[protoc](https://github.com/protocolbuffers/protobuf.git)


# install protoc-gen-go
```
go get -u github.com/golang/protobuf/protoc-gen-go
```

# generate go file
```
protoc --go_out=. *.proto
```