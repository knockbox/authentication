# knockbox/authentication

The User & Authentication service.

# Docker

This assumes you have the `local/docker-compose` found in [knockbox/architecture](https://github.com/knockbox/architecture)
up and running to ensure the required database, cache, etc. are all existing for you
to work with. If you are using other hosted services feel free to stand this up however you see
fit.

## Build

You can build the current environment using:

```shell
docker build -t knockbox/auth .
```

## Run

The container can be run locally using:

```shell
docker run -p 9090:9090 --env-file .env.docker-local knockbox/auth
```