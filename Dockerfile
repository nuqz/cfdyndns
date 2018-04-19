FROM golang:alpine AS build-env

COPY . $GOPATH/src/nuqz/cfdyndns
RUN go build -o /tmp/cfdyndns nuqz/cfdyndns

FROM scratch
COPY --from=build-env /tmp/cfdyndns /usr/local/bin/cfdyndns
CMD ["/usr/local/bin/cfdyndns", "-config", "/etc/cfdyndns.yml"]
