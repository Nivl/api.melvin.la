#!/bin/bash
alias ddc-build="ML_BUILD_ENV=test docker-compose build" # builds the services
alias ddc-up="ML_BUILD_ENV=test docker-compose up -d" # starts the services
alias ddc-rm="ML_BUILD_ENV=test docker-compose stop && docker-compose rm -f" # Removes the services
alias ddc-stop="ML_BUILD_ENV=test docker-compose stop" # Stops the running services

alias ml-log-mongo="docker logs ml_api_mongodb" # print mongo logs

# Execute any command in the container
function ml-exec {
  CMD="cd /go/src/github.com/Nivl/api.melvin.la && $@"
  docker exec -i -t ml_api /bin/bash -ic $CMD
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
  ml-exec "go test $@"
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

function ml-reset-mongo {
  CMD="mongo api-melvin --eval \"printjson(db.dropDatabase())\""
  docker exec -i -t ml_api_mongodb /bin/bash -ic $CMD
}