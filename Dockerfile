FROM golang:1.21-alpine AS build_stage

ARG TARGET
ARG BUILD_OS=linux
ARG BUILD_ARCHITECTURE=amd64
ARG USER_ID=1000
ARG GROUP_ID=1000

RUN apk update \
  && apk add --no-cache ca-certificates tzdata \
  && update-ca-certificates \
  && addgroup -g $GROUP_ID appuser \
  && adduser -u $USER_ID -G appuser -s /bin/sh -D appuser 

WORKDIR /usr/src/app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN GOOS=$BUILD_OS GOARCH=$BUILD_ARCHITECTURE go build -o bin/app -v -ldflags="-s -w" cmd/$TARGET/main.go

FROM alpine:3

WORKDIR /app

COPY --from=build_stage /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build_stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build_stage /etc/passwd /etc/passwd
COPY --from=build_stage /etc/group /etc/group

COPY --from=build_stage /usr/src/app/config.yaml .
COPY --from=build_stage /usr/src/app/bin/app .

ENV TZ=UTC

CMD ["/app/app"]
