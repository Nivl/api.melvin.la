FROM golang:1.7

# install depedencies
RUN go get bitbucket.org/liamstask/goose/cmd/goose

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/Nivl/api.melvin.la

# Install api binary globally within container
RUN cd /go/src/github.com/Nivl/api.melvin.la && make install

# Set binary as entrypoint
ENTRYPOINT /go/bin/api

EXPOSE 5000