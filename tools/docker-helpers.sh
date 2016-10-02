#!/bin/bash
alias ddc-build="docker-compose build" # builds the container
alias ddc-up="docker-compose up -d" # starts the container
alias ddc-rm="docker-compose stop && docker-compose rm -f" # Removes the renning container
alias ddc-stop="docker-compose stop" # Stops the renning container

alias ml-log-mongo="docker logs apimelvinla_database_1" # print mongo logs

# Execute any command in the container
function ml-exec {
  docker exec -i -t apimelvinla_api_1 /bin/bash -ic "cd /go/src/github.com/Nivl/api.melvin.la && $@"
}

# Open a bash session
function ml-bash {
  ddc-ml-exec bash
}

# Execute a make command
function ml-make {
  ddc-ml-exec "make $@"
}

# Execute a go command
function ml-go {
  ddc-ml-exec "go $@"
}