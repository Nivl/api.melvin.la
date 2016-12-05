#!/usr/bin/env bash

set -ex

heroku plugins:install heroku-container-registry
heroku container:login
heroku container:push web --app melvin-laplanche
heroku run "cd /go/src/github.com/melvin-laplanche/ml-api && goose up" --app melvin-laplanche