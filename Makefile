# Build info
VERSION=1.0.0
BUILD_INFO=`git rev-parse HEAD`

# Commands
GO_CMD=go
INSTALL_CMD=$(GO_CMD) install

# Input/Output
SOURCES_DIR=./api

# Flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD_INFO)"

install:
	$(INSTALL_CMD) $(LDFLAGS) $(SOURCES_DIR)

.PHONY:
	install