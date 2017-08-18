# ml-api

## Master badges
[![Build Status](https://travis-ci.org/melvin-laplanche/ml-api.svg?branch=master)](https://travis-ci.org/melvin-laplanche/ml-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/melvin-laplanche/ml-api)](https://goreportcard.com/report/github.com/melvin-laplanche/ml-api)
[![codecov](https://codecov.io/gh/melvin-laplanche/ml-api/branch/master/graph/badge.svg)](https://codecov.io/gh/melvin-laplanche/ml-api)
[![BCH compliance](https://bettercodehub.com/edge/badge/melvin-laplanche/ml-api?branch=master)](https://bettercodehub.com/results/melvin-laplanche/ml-api)

## Staging badges
[![Build Status](https://travis-ci.org/melvin-laplanche/ml-api.svg?branch=staging)](https://travis-ci.org/melvin-laplanche/ml-api)
[![codecov](https://codecov.io/gh/melvin-laplanche/ml-api/branch/staging/graph/badge.svg)](https://codecov.io/gh/melvin-laplanche/ml-api)
[![BCH compliance](https://bettercodehub.com/edge/badge/melvin-laplanche/ml-api?branch=staging)](https://bettercodehub.com/results/melvin-laplanche/ml-api)

## Run the API using docker

```
docker-compose build
docker-compose up -d
```

Bash helpers can be found in `tools/docker-helpers.sh`

## travis

```
travis encrypt HEROKU_API_KEY=$(heroku auth:token) --add
travis encrypt APIARY_API_KEY=your-token --add
travis encrypt your-email-address
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