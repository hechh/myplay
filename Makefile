
GCFLAGS=-gcflags "all=-N -l"
SYSTEM=$(shell go env GOOS)

# 默认包
import_pb=myplay/common/pb
# 生成路径
gen_path_output=./output
gen_path_pb2redis=./common/dao/gen
gen_path_pb2code=./common/table
gen_path_pb=./common/pb
gen_descriptor=${gen_path_output}/proto.descriptor
# 配置路径
path_data=./configure/data
path_table=./configure/table
path_scripts=./configure/scripts
path_proto=./configure/protocol
path_packet=./framework/packet
# 工具路径
tool_path_pb=./tools/pbtool
tool_path_xlsx=./tools/xlsx

## 需要编译的服务
TARGET=gate db game client
TOOLS=pb pb2redis xlsx xlsx2code xlsx2proto xlsx2data 
LINUX=$(TARGET:%=%_linux)
BUILD=$(TARGET:%=%_build)

.PHONY: ${TARGET} ${TOOLS} docker_stop docker_run develop 

all: clean
	make ${BUILD}

linux: clean
	make ${LINUX}

build: clean ${BUILD}

clean:
	@rm -rf ${gen_path_output}

$(LINUX): %_linux: %
	@echo "Building $*"
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build ${GCFLAGS} -o ${gen_path_output}/ ./server/$*

$(BUILD): %_build: %
	@echo "Building $*"
	go build ${GCFLAGS} -o ${gen_path_output}/$* ./server/$*

pb:
	@echo "Building pb"
	@rm -rf ${gen_path_pb} && ${path_scripts}/pb_gen.sh
	@go run ${tool_path_pb}/main.go -src=${gen_path_pb} -dst=${gen_path_pb} -a=pb 
pb2redis:
	@echo "gen redis code..."
	@rm -rf ${gen_path_pb2redis}
	@go run ${tool_path_pb}/main.go -src=${gen_path_pb} -dst=${gen_path_pb2redis} -a=redis -i=${import_pb}
xlsx:
	@make xlsx2proto && make pb && make xlsx2data && make xlsx2code
xlsx2code:
	@mkdir -p ${gen_path_output}
	@protoc --descriptor_set_out=${gen_descriptor} --include_imports -I${path_packet} -I${path_proto} ${path_proto}/*.proto
	@rm -rf ${gen_path_pb2code}
	@go run ${tool_path_xlsx}/main.go -src=${path_table} -dst=${gen_path_pb2code} -a=code -p=myplay -d=${gen_descriptor} -i=${import_pb}
xlsx2proto:
	@rm -rf ${path_proto}/*.gen.proto
	@go run ${tool_path_xlsx}/main.go -src=${path_table} -dst=${path_proto} -a=proto -o=./pb -p=myplay
xlsx2data:
	@mkdir -p ${gen_path_output}
	@protoc --descriptor_set_out=${gen_descriptor} --include_imports -I${path_packet} -I${path_proto} ${path_proto}/*.proto
	@rm -rf ${path_data}
	@go run ${tool_path_xlsx}/main.go  -src=${path_table} -dst=${path_data} -a=data -p=myplay -d=${gen_descriptor}
docker_stop:
	@echo "停止docker环境"
	-cd ${gen_path_output} && docker-compose -f docker_compose.yaml down
docker_run:
	@echo "启动docker环境"
	-cd ${gen_path_output} && docker-compose -f docker_compose.yaml up -d
develop:
	@mkdir -p ${gen_path_output}
	@cp -rf ./configure/env/develop/* ./configure/data ./output


