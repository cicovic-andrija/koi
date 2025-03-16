# koipond: personal inventory server

## Contents

- Concepts
- Customization
- Development Instructions
- Deployment Instructions

## Concepts

TODO

## Instructions

### Build

```bash
go build -o koipond main.go
```

### Run locally (development mode)

TODO: Mention store/

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

TODO: Mention store/

```bash
docker run \
    --name inventory-server \
    --publish 127.0.0.1:8072:8072 \
    --env 'KOIPOND_PORT=8072' \
    --volume $HOME/store:/srv/store \
    --restart on-failure:10 \
    --detach \
    acicovic/koipond:latest
```

### Kill

```bash
pkill -SIGINT koipond
```

> Or, CTRL-C if running in foreground.

or

```bash
docker kill -s SIGINT koipond-server
```

### Build Docker image

```bash
docker build -t acicovic/koipond:latest .
```

> Substitute image name.
