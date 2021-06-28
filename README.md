# robotech

# golang
- Download [Golang](https://golang.org/dl/)
- `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.4.linux-amd64.tar.gz`
- Add to `PATH` environment variable
```
# set go env
# go path
if [ -d "$HOME/.go" ] ; then
        GOPATH="$HOME/.go"
        PATH="$GOPATH/bin:$PATH"
fi
# go root
if [ -d "/usr/local/go" ] ; then
        GOROOT="/usr/local/go"
        PATH="$GOROOT/bin:$PATH"
fi
# go proxy
export GOPROXY=https://goproxy.io,direct
```

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