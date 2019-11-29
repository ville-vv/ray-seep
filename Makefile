DOCKER_BUILD=./docker_build
BINARY_DIR=./bin
PROJECTNAME=$(shell basename "$(PWD)")
DOCKER_BINARY=$(DOCKER_BUILD)/$(PROJECTNAME)

SERVER_UNIX_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_linux
SERVER_MAC_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_mac
SERVER_WIN_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_win.exe

CLIENT_UNIX_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_linux_cli
CLIENT_MAC_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_mac_cli
CLIENT_WIN_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_win_cli.exe

.PHONY:all test image clean build 

all: build

test:
	go test  -cover ./...

build: clean build-server build-client

build-server:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(SERVER_UNIX_BINARY)  $(PROJECTNAME)/run/server/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(SERVER_MAC_BINARY)  $(PROJECTNAME)/run/server/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(SERVER_WIN_BINARY)  $(PROJECTNAME)/run/server/main.go
build-client:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(CLIENT_UNIX_BINARY)  $(PROJECTNAME)/run/client/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(CLIENT_MAC_BINARY)  $(PROJECTNAME)/run/client/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(CLIENT_WIN_BINARY)  $(PROJECTNAME)/run/client/main.go

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
