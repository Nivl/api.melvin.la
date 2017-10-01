#!/bin/bash

alias ddc-build="docker-compose build" # builds the services
alias ddc-up="docker-compose up -d" # starts the services
alias ddc-rm="docker-compose stop && docker-compose rm -f" # Removes the services
alias ddc-stop="docker-compose stop" # Stops the running services

# Execute any command in the container
function ml-exec {
  CMD="cd /go/src/github.com/melvin-laplanche/ml-api && $@"
  docker-compose exec api /bin/bash -ic $CMD
}

# Execute any command in the container
function ml-psql {
  docker-compose exec database psql -U $POSTGRES_USER
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

# Remove and rebuild the containers
function ml-reset {
  source config/api.env

  ddc-rm
  ddc-up

  until docker-compose exec database psql "$API_POSTGRES_URI_STR" -c "select 1" > /dev/null 2>&1; do sleep 2; done
}

# Execute a test
function ml-test {
  echo "Restart services..."
  ddc-stop &> /dev/null
  ddc-up &> /dev/null

  echo "Start testings"
  ml-exec "go test -tags=integration $@"
}

# Execute a test
function ml-tests {
  echo "Restart services..."
  ddc-stop &> /dev/null
  ddc-up &> /dev/null

  echo "Start testings"
  ml-exec "cd src && go test -tags=integration ./..."
}

# Execute a test
function ml-coverage {
  echo "Restart services..."
  ddc-stop &> /dev/null
  ddc-up &> /dev/null

  echo "Start testings"
  ml-exec "cd src && ../go.test.sh"
}
