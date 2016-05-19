# Build info
VERSION=1.0.0
BUILD_INFO=`git rev-parse HEAD`

# Commands
GO_CMD=go
BUILD_CMD=$(GO_CMD) build
INSTALL_CMD=$(GO_CMD) install
GET_CMD=$(GO_CMD) get
LINT_CMD=golint

# Input/Output
SOURCES_DIR=./src
OUTPUT=./bin/api-melvin-la

# Flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD_INFO)"

.DEFAULT_GOAL: $(OUTPUT)
$(OUTPUT): build

get-dep:
	$(GET_CMD) github.com/ddollar/forego
	$(GET_CMD) github.com/onsi/ginkgo/ginkgo
	$(GET_CMD) github.com/onsi/gomega

build:
	$(BUILD_CMD) -v $(LDFLAGS) -o $(OUTPUT) $(SOURCES_DIR)

install:
	$(INSTALL_CMD) $(LDFLAGS) -o $(OUTPUT) $(SOURCES_DIR)

lint:
	$(LINT_CMD) $(SOURCES_DIR)

test:
	ginkgo -r -cover $(SOURCES_DIR)

run:
	forego start

clean:
	if [ -f $(OUTPUT) ] ; then rm $(OUTPUT) ; fi

.PHONY:
	all build run clean install get-dep test lint
