.PHONY: all build run gotool  clean help
BINARY=gva-server
all: clean npmBuild gotool  build
npmBuild:
	cd ../web && npm run build && cp -r ./dist/* ../server/fe
build:
	# 将静态文件打包到二进制文件中
	go generate -tags=packfile ./packfile/usePackFile.go
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o ${BINARY}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=packfile  -o ${BINARY}
clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
run:
	@go run ./

gotool:
	cd ../server
	go fmt ./
	go vet ./
help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"