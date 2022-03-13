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

```

```
go env -w GOPATH=$HOME/go
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=*.gitlab.com,*.gitee.com,*.coding.net
```

unknown
```
go env -w GOSUMDB=“sum.golang.google.cn”
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