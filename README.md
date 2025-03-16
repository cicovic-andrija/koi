# koipond: personal inventory server

> v1.x | self-hosted | do-it-yourself | minimal | web-server | docker-containerized

## About

`koipond` is a software system for personal inventory management. The north star of the project is
to provide a simple and minimal web interface that can be used to browse and update small-scale
(thousands) collections of items, with raw data persisted in textual, human-readable, standard
format. It is designed to be easily extensible, customizable and self-hosted, for users that are
comfortable with technology and prefer minimal systems.

In its current state, the server offers only browsing (read-only) functionality, meaning that the
XML file used for persistent storage must be populated/updated through a text editor program. For
the one person currently working on this project, that was a high enough bar for version `1.x` of
the system (demo link below). This decision also aligns with the idea of target audience for this
system. Version `2.x` will offer read-write capabilities (to be implemented).

#### Version `1.x` Demo: [https://inventory.acicovic.me](https://inventory.acicovic.me)

## Concepts

For the rest of this document, notes related to source code pointers and implementation details
will be visually separated by using Markdown footnote formatting, for example:

> The text you are currently reading can be found in the `README.md` at the root of the source tree.

Concepts:


## Instructions

## How To: Programming

### Build

```bash
$ go build -o koipond main.go
```

### Run locally (development mode)

TODO: Mention store/

```bash
$ KOIPOND_MODE=dev go run main.go
```

## How To: Deployment

> TCP port in dev mode is hard-coded to 8072.

### Run in production

```bash
$ KOIPOND_MODE=prod-local-listener KOIPOND_PORT=52000 ./koipond
```

> Requires a service manager to handle crashes and log redirection. See `systemd.service` for an example.

> For encrypted traffic, configure a reverse HTTPS proxy, e.g. `nginx`.

> For authentication, configure a stanalone authentication service.

### Run in production (Docker)

TODO: Mention store/

```bash
$ docker run \
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
$ pkill -SIGINT koipond
```

> Or, CTRL-C if running in foreground.

or

```bash
$ docker kill -s SIGINT koipond-server
```

### Build Docker image

```bash
$ docker build -t acicovic/koipond:latest .
```

> Substitute image name.
