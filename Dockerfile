FROM golang:1.14.3-alpine
#server port
EXPOSE 1111
#Naming port
EXPOSE 1234
RUN apk update
RUN apk add vim
#No need for copy if we are mirroring volumes
#COPY . .
