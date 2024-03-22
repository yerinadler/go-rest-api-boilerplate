FROM golang:1.21-alpine AS build_stage

ARG TARGET=api

WORKDIR /usr/src/app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o bin/app -v -ldflags="-s -w" cmd/$TARGET/main.go

FROM alpine:3

WORKDIR /app

COPY --from=build_stage /usr/src/app/config.yaml .
COPY --from=build_stage /usr/src/app/bin/app .

CMD ["/app/app"]
