APP=Citrix-NetScaler-Exporter
VERSION=3.2.1
BINARY-WINDOWS-X86=${APP}_${VERSION}_Windows_32bit.exe
BINARY-WINDOWS-X64=${APP}_${VERSION}_Windows_amd64.exe
BINARY-LINUX=${APP}_${VERSION}_amd64

BUILD=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags="-X main.version=${VERSION} -X main.build=${BUILD}"

build-binary-windows-x86 :
	@echo "=== Building Windows x86 ==="
	(GOOS=windows; GOARCH=386; go build -o ${BINARY-WINDOWS-X86} ${LDFLAGS})

build-binary-windows-x64 :
	@echo "=== Building Windows x64 ==="
	(GOOS=windows; GOARCH=amd64; go build -o ${BINARY-WINDOWS-X64} ${LDFLAGS})

build-binary-linux :
	@echo "=== Building Linux ==="
	(GOOS=linux; GOARCH=amd64; go build -o ${BINARY-LINUX} ${LDFLAGS})

build: build-binary-linux build-binary-windows-x64 build-binary-windows-x86
