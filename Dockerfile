FROM golang:1.14.3-alpine
#server port
EXPOSE 1111
#Naming port
EXPOSE 1234

COPY . .
