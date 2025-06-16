# Container Compose CLI

This project provides developer ergonomics to start using Apple's [Container](https://github.com/apple/container). Inspired by [Docker Compose](https://github.com/docker/compose), this CLI offers a similar experience for managing containerised applications - without the need for Docker Desktop.

## Installation

Download and install Container [from the release page](https://github.com/apple/container/releases).

Install this CLI:

```bash
go install github.com/container-compose/cli@latest
```

## Building from source

```bash
go build -o container-compose ./cmd
./container-compose -h
```

## Usage

Make sure that you have started the container system:

```bash
container system start
```

Write your first compose spec:

```yaml
version: "0.1"
services:
  web:
    image: nginx:latest
    ports:
      - "8080:80"
```

### `container-compose up`

```bash
container-compose up -f compose.yaml
```
