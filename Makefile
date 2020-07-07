DOCKER_BUILD=./docker_build
BINARY_DIR=./bin
PROJECTNAME=$(shell basename "$(PWD)")
DOCKER_BINARY=$(DOCKER_BUILD)/$(PROJECTNAME)
OSMAC=mac
OSWIN=windows
OSUNIX=linux

SERVER_UNIX_BINARY=$(BINARY_DIR)/$(OSUNIX)/$(PROJECTNAME)
SERVER_MAC_BINARY=$(BINARY_DIR)/$(OSMAC)/$(PROJECTNAME)
SERVER_WIN_BINARY=$(BINARY_DIR)/$(OSWIN)/$(PROJECTNAME).exe

CLIENT_UNIX_BINARY=$(BINARY_DIR)/$(OSUNIX)/$(PROJECTNAME)-cli
CLIENT_MAC_BINARY=$(BINARY_DIR)/$(OSMAC)/$(PROJECTNAME)-cli
CLIENT_WIN_BINARY=$(BINARY_DIR)/$(OSWIN)/$(PROJECTNAME)-cli.exe

HTTP_WEB_UNIX_BINARY=$(BINARY_DIR)/$(OSUNIX)/http-web
HTTP_WEB_MAC_BINARY=$(BINARY_DIR)/$(OSMAC)/http-web
HTTP_WEB_WIN_BINARY=$(BINARY_DIR)/$(OSWIN)/http-web.exe

.PHONY:all test image clean build

all: build

test:
	go test  -cover ./...

build: clean build-server build-client build-web-tool

build-server:
	echo $(PROJECTNAME)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(SERVER_UNIX_BINARY)  ./apps/server/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(SERVER_MAC_BINARY)  ./apps/server/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(SERVER_WIN_BINARY)  ./apps/server/main.go
build-client:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(CLIENT_UNIX_BINARY)  ./apps/client/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(CLIENT_MAC_BINARY)  ./apps/client/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(CLIENT_WIN_BINARY)  ./apps/client/main.go

build-web-tool:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(HTTP_WEB_UNIX_BINARY)  ./apps/http-web/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(HTTP_WEB_MAC_BINARY)  ./apps/http-web/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(HTTP_WEB_WIN_BINARY)  ./apps/http-web/main.go

docker-ready:
	sudo rm -rf $(DOCKER_BUILD)
	mkdir -p $(DOCKER_BUILD)
	cp docker/* $(DOCKER_BUILD)
	cp -r $(SERVER_UNIX_BINARY) $(DOCKER_BUILD)/$(PROJECTNAME)
	echo $(PROJECTNAME)

docker-build: docker-ready
	sudo docker build -t $(PROJECTNAME) $(DOCKER_BUILD)

docker_stop:
	sudo docker stop $(sudo docker container ls -a | grep "tstl" | awk '{print $1}')
docker_rm:
	sudo docker rm $(sudo docker container ls -a | grep "tstl" | awk '{print $1}')
docker_rmi:
	sudo docker rmi $(sudo docker images | grep "none" | awk '{print $3}')

clean:
	rm -rf bin
	rm -rf docker_build
