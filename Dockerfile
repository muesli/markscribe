FROM docker.io/golang:1.25.5-alpine3.23 AS build
LABEL org.opencontainers.image.source=https://gitlab.hufschlaeger.net/
LABEL org.opencontainers.image.description="My container image"
LABEL org.opencontainers.image.licenses=MIT
WORKDIR /go/markscribe
COPY . .
RUN apk update --no-cache \
    && apk add --no-cache make zip golangci-lint \
    && make build-linux

FROM alpine:3.23 AS prod

COPY --from=build /go/markscribe/bin/markscribe_unix /go/bin/markscribe
ENTRYPOINT ["/go/bin/markscribe"]