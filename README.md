# ml-api

## Master badges
[![Build Status](https://travis-ci.org/melvin-laplanche/ml-api.svg?branch=master)](https://travis-ci.org/melvin-laplanche/ml-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/melvin-laplanche/ml-api)](https://goreportcard.com/report/github.com/melvin-laplanche/ml-api)
[![codebeat badge](https://codebeat.co/badges/111cf407-0776-4331-96d2-da2e4df9c4f5)](https://codebeat.co/projects/github-com-melvin-laplanche-ml-api)

## Staging badges
[![Build Status](https://travis-ci.org/melvin-laplanche/ml-api.svg?branch=staging)](https://travis-ci.org/melvin-laplanche/ml-api)

## Run the API using docker

```
docker-compose build
docker-compose up -d
```

Bash helpers can be found in `tools/docker-helpers.sh`

## travis

```
travis encrypt HEROKU_API_KEY=$(heroku auth:token) --add
```

## Documentation

The documentation is using blueprint and the generated file can be uploaded to apiary

### Install hercule

```
  yarn global add hercule
  # OR
  npm install hercule -g
```

### Update the documentation

```
hercule doc/main.apib -o doc.apib
```