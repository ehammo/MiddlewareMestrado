# How to run this project

## Pre requirements

Before running the docker file, remember to run
export DOCKER_BUILDKIT=1

## Build

docker build . -t middleware

## Run container

1. Naming
docker run --name naming -it middleware

### On container

cd main && go run NamingServer.go

2. Server
docker run --name server -it --link=naming middleware

### On container

cd main && go run runserver.go

3. Client(s)
docker run --name client -it --link=naming --link=server middleware

### On container

go run main/runclient.go