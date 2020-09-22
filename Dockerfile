FROM golang:1.14.3-alpine AS build
#server port
EXPOSE 1111
#Naming port
EXPOSE 1234

WORKDIR /src
COPY . .
RUN go build -o /out/example .

FROM scratch AS bin
COPY --from=build /out/example /

