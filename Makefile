OS := $(shell uname -s)
IS_LINUX := $(shell echo $(OS)|grep -i linux)
IS_MAC := $(shell echo $(OS)|grep -i Drawin)

PWD=$(shell pwd)
remote=$(shell cat .devhost)

BIN_DIR=$(PWD)/build

LDFLAGS = -X main.buildPath=${PWD}

# 编译时间
COMPILE_TIME = $(shell date +"%Y-%m-%dT%H:%M:%S%z")
LDFLAGS += -X main.buildAt=$(COMPILE_TIME)

# GIT版本号
GIT_REVISION = $(shell git show -s --pretty=format:%h)
LDFLAGS += -X main.buildVer=$(GIT_REVISION)

WEB_OUTER = ${LDFLAGS}
WEB_INNER = ${LDFLAGS}
WEB_TEST = ${LDFLAGS}

cw:
	go build -ldflags "$(WEB_INNER)" -o ./build/push ./src/app

mac_cw:
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(WEB_INNER)" -o ./build/bin/mom ./src/

rg:
	./build/push gateway --app_path="$(PWD)" --app_mode="local"

rc:
	./build/push client --app_path="$(PWD)" --app_mode="local"

rl:
	./build/push logic --app_path="$(PWD)" --app_mode="local"

