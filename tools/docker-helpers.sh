#!/bin/bash
alias ddc-build="ML_BUILD_ENV=travis docker-compose build" # builds the services
alias ddc-up="ML_BUILD_ENV=travis docker-compose up -d" # starts the services
alias ddc-rm="docker-compose stop && docker-compose rm -f" # Removes the services
alias ddc-stop="docker-compose stop" # Stops the running services

alias ml-log-mongo="docker logs apimelvinla_database_1" # print mongo logs

# Execute any command in the container
function ml-exec {
  CMD="cd /go/src/github.com/Nivl/api.melvin.la && $@"
  docker exec -i -t apimelvinla_api_1 /bin/bash -ic $CMD
}

# Open a bash session
function ml-bash {
  ml-exec bash
}

# Execute a make command
function ml-make {
  ml-exec make "$@"
}

# Execute a go command
function ml-go {
  ml-exec go "$@"
}

# Execute a test
function ml-test {
  echo "Restart services..."
  ddc-stop &> /dev/null
  ddc-build &> /dev/null
  ddc-up &> /dev/null

  echo "Start testings"
  ml-exec go test "$@"
}

# Execute a test
function ml-tests {
  echo "Restart services..."
  ddc-stop &> /dev/null
  ddc-build &> /dev/null
  ddc-up &> /dev/null

  echo "Start testings"
  ml-exec "cd api && go test ./..."
}