FROM golang:alpine AS build-env

COPY . $GOPATH/src/nuqz/cfdyndns
RUN CGO_ENABLED=0 go build -o /tmp/cfdyndns nuqz/cfdyndns

FROM alpine
RUN apk update && apk add --no-cache ca-certificates
COPY --from=build-env /tmp/cfdyndns /usr/local/bin/cfdyndns
CMD ["/usr/local/bin/cfdyndns", "-config", "/etc/cfdyndns/cfdyndns.yml"]
