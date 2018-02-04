#!/bin/sh

# 资源文件打包
go-bindata-assetfs -ignore=\\.git -ignore=\\.gitignore -ignore=\\.DS_Store -pkg=main static/...

# 编译到 $GOPATH/bin
go install -ldflags "-s -w" .