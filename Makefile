# Target configuration
TARGET_OS		:= $(go env GOOS)
TARGET_ARCH		:= amd64
CGO_ENABLED		:= 0
EXTRA_FLAGS     :=

VERSION			= $(shell git rev-list -1 HEAD)
APP_NAME		:= dq
VERSION_SETTING	= dq/cmd._version
NAME_SETTING	= dq/cmd._name

build:
	go get ./...
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) CGO_ENABLED=$(CGO_ENABLED) \
	go build $(EXTRA_FLAGS) -ldflags "-X $(VERSION_SETTING)=$(VERSION) -X $(NAME_SETTING)=$(APP_NAME)" -o $(APP_NAME)
