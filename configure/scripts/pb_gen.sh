#!/usr/bin/env bash

GO_BIN="$(go env GOPATH)/bin"
SYSTEM=$(go env GOOS)
PacketPath=./framework/packet
ProtoPath=./configure/protocol
PbGoPath=./common/pb

rm -rf ${PbGoPath} && mkdir -p ${PbGoPath}

if [ "${SYSTEM}" == "windows" ]; then
    protoc.exe -I${ProtoPath} -I${PacketPath} ${ProtoPath}/*.proto --go_opt paths=source_relative --go_out=${PbGoPath}
    protoc-go-inject-tag.exe -input=${PbGoPath}/*.pb.go -XXX_skip="state,sizeCache,unknownFields"
else
	protoc -I${ProtoPath} -I${PacketPath} ${ProtoPath}/*.proto --go_opt paths=source_relative --go_out=${PbGoPath}
    protoc-go-inject-tag -input=${PbGoPath}/*.pb.go -XXX_skip="state,sizeCache,unknownFields"
fi

# 使用 sed 批量添加忽略标签
if [ "${SYSTEM}" == "darwin" ]; then
    sed -i '' -E 's/(^[[:space:]]*(state|sizeCache|unknownFields)[[:space:]]+protoimpl\.[[:alpha:]]+)/\1 `xorm:"-"`/' ${PbGoPath}/*.pb.go
    sed -i '' 's/`protogen:"open.v1"`//g' ${PbGoPath}/*.pb.go
else 
    sed -i -E 's/(^\s*(state|sizeCache|unknownFields)\s+protoimpl\.[A-Za-z]+)/\1 `xorm:"-"`/' ${PbGoPath}/*.pb.go
    sed -i 's/`protogen:"open.v1"`//g' ${PbGoPath}/*.pb.go
fi
