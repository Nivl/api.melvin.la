# Build info
VERSION=1.0.0
BUILD_INFO=`git rev-parse HEAD`

# Flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD_INFO)"

install:
	go install $(LDFLAGS) .

migration:
	goose up

.PHONY:
	install