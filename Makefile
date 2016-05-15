VERSION=1.0.0
BUILD_INFO=`git rev-parse HEAD`

GO_CMD=go
BUILD_CMD=$(GO_CMD) build
INSTALL_CMD=$(GO_CMD) install
LINT_CMD=golint

SOURCES_DIR=./src
OUTPUT=./bin/api-melvin-la

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD_INFO)"

.DEFAULT_GOAL: $(OUTPUT)
$(OUTPUT): build

build:
	$(BUILD_CMD) -v $(LDFLAGS) -o $(OUTPUT) $(SOURCES_DIR)

install:
	$(INSTALL_CMD) $(LDFLAGS) -o $(OUTPUT) $(SOURCES_DIR)

lint:
	$(LINT_CMD) $(SOURCES_DIR)

run:
	forego start

clean:
	if [ -f $(OUTPUT) ] ; then rm $(OUTPUT) ; fi

.PHONY:
	all build run clean install
