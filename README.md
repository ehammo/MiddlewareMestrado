# How to run this project

Before running the docker file, remember to run

export DOCKER_BUILDKIT=1

docker build --target bin --output bin/ .

To run execute:
docker run -p 1111:1111 server_image -it