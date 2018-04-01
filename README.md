# cfdyndns

CloudFlare dynamic DNS client

## Installing

```
$ go get github.com/nuqz/cfdyndns
$ go install github.com/nuqz/cfdyndns
```

## Config

Example:

```toml
"example.com" = ["example.com", "subdomain.example.com", "mail.example.com"]
"my.website" = ["my.website"]
```

## Run

```
$ CF_API_EMAIL=account@gmail.com CF_API_KEY=secret-key ./cfdyndns -interval 30s -config ~/dns.toml
```

or

```
$ make local
```

## Install as systemd service

```
$ make install_service
```

