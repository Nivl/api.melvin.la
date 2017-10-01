FROM golang:1.9

# install depedencies
RUN go get github.com/pressly/goose/cmd/goose

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/melvin-laplanche/ml-api

# Install api binary globally within container
RUN cd /go/src/github.com/melvin-laplanche/ml-api && make install

# Set binary as entrypoint
CMD /go/bin/ml-api

EXPOSE 5000