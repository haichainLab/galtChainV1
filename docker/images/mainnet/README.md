# Supported tags

## Simple Tags

- latest
- latest-arm32v5
- latest-arm32v6
- latest-arm32v7
- latest-arm64v8
- develop
- develop-arm32v5
- develop-arm32v6
- develop-arm32v7
- develop-arm64v8
- release-v0.23.0
- release-v0.23.0-arm32v5
- release-v0.23.0-arm32v6
- release-v0.23.0-arm32v7
- release-v0.23.0-arm64v8
- release-v0.22.0

## Building your own images

This Dockerfile build your working copy by default, but if you pass the
GALTCOIN_VERSION build argument to the `docker build` command, it will checkout
to the branch, a tag or a commit you specify on that variable.

Example

```sh
$ git clone https://github.com/samoslab/galtcoin
$ cd galtcoin
$ GALTCOIN_VERSION=v1.24.1
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=GALTCOIN_VERSION=$GALTCOIN_VERSION \
  -t galtcoin:$GALTCOIN_VERSION .
```

or just

```sh
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=GALTCOIN_VERSION=v1.24.1 \
  -t galtcoin:v1.24.1
```

## ARM Architecture

Build arguments are provided to make it easy if you want to build for the ARM
architecture.

Example for ARMv5.

```sh
$ git clone https://github.com/samoslab/galtcoin
$ cd galtcoin
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=ARCH=arm \
  --build-arg=GOARM=5 \
  --build-arg=IMAGE_FROM="arm32v5/alpine" \
  -t galtcoin:$GALTCOIN_VERSION-arm32v5 .
```

## How to use this images

### Run a galtcoin node

This command pulls latest stable image from Docker Hub, and launches a node inside a Docker container that runs as a service daemon in the background. It is possible to use the tags listed above to run another version of the node

```sh
$ docker volume create galtcoin-data
$ docker volume create galtcoin-wallet
$ docker run -d -v galtcoin-data:/data/.galtcoin \
  -v galtcoin-wallet:/wallet \
  -p 6868:6868 -p 6869:6869 \
  --name galtcoin-node-stable samoslab/galtcoin
```

When invoking the container this way the options of the galtcoin command are set to their respective default values , except the following

| Parameter  | Value |
| ------------- | ------------- |
| web-interface-addr | 0.0.0.0  |
| gui-dir | /usr/local/galtcoin/src/gui/static |

In order to stop the container , just run

```sh
$ docker stop galtcoin-node-stable
```

Restart it once again by executing

```sh
$ docker start galtcoin-node-stable
```
