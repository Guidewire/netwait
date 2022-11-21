# Netwaiter

Netwaiter is a CLI utility used to wait for a network resource (URL, host:port)
to become available.

The utility is a single binary with no dependencies so that it is easy to
include in various workflows.  It is meant to replace script-based logic, often
using different tools like curl or netcat, to test whether a network resource is
available.

# Installation

Netwaiter is provided as platform-specific binaries and a Docker image.

For platform binaries, see [releases](https://github.com/merusso/netwaiter/releases)

Docker image: [merusso/netwaiter](https://hub.docker.com/r/merusso/netwaiter)

# Usage

```bash
# Wait for resource to become available
$ netwaiter wait https://github.com
```

Netwaiter returns a non-zero exit code if it was unable to successfully connect
to the network resource within a certain time period.

```bash
$ if ./netwaiter wait bad-domain.fake:443; then echo 'SUCCESS'; else echo 'FAIL'; fi
Error: All attempts fail:
#1: dial tcp: lookup bad-domain.fake: no such host
#2: context deadline exceeded
FAIL
```

## HTTP resource

Netwaiter can wait for an HTTP URL. Success is defined as an HTTP response with
a 2xx status code. HTTP Redirects are followed.

```bash
$ netwaiter wait http://httpbin.org/status/200

$ netwaiter wait -t 10s http://httpbin.org/status/500
Error: All attempts fail:
#1: GET 'http://httpbin.org/status/500': returned status code 500
#2: GET 'http://httpbin.org/status/500': returned status code 500
#3: GET 'http://httpbin.org/status/500': returned status code 500
#4: context deadline exceeded
```

## TCP resource

```bash
$ netwaiter wait github.com:443
```

## Multiple resources

Netwaiter will attempt to connect to multiple resources in parallel. All
resources must successfully connect for the command to return success.

```bash
$ netwaiter wait https://github.com https://go.dev/
```

## Timeout

```bash
# stop waiting after up to 10 seconds
$ netwaiter wait --timeout 10s github.com:443

# stop waiting after up to 2 minutes
$ netwaiter wait --timeout 2m github.com:443
```

## Docker

Netwaiter is available as a Docker image. The image is configured with an
`ENTRYPOINT` pointing to `netwaiter`, so `CMD` arguments are treated as
arguments passed to `netwaiter`

```bash
$ docker run --rm merusso/netwaiter wait https://github.com
```
