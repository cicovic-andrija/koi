# Koipond

## Concepts

TODO

## Instructions

### Build

```bash
go build -o koipond main.go
```

### Run locally (development mode)

```bash
KOIPOND_MODE=dev go run main.go
```

> TCP port in dev mode is hard-coded to 8072.

### Run in production

```bash
KOIPOND_MODE=prod-local-listener KOIPOND_PORT=52000 ./koipond
```

> Requires a service manager to handle crashes and log redirection. See `systemd.service` for an example.

> For encrypted traffic, configure a reverse HTTPS proxy, e.g. `nginx`.

> For authentication, configure a stanalone authentication service.

### Run in production (Docker)

```bash
docker run --name koipond-server --publish 127.0.0.1:8072:8072 --env 'KOIPOND_PORT=8072' --detach acicovic/koipond:latest
```

### Kill

```bash
pkill -SIGINT koipond
```

> Or, CTRL-C if running in foreground.
