# Netwait

Netwait is a CLI utility used to wait for a network resource (URL, host:port)
to become available.

The utility is a single binary with no dependencies so that it is easy to
include in various workflows.  It is meant to replace script-based logic, often
using different tools like curl or netcat, to test whether a network resource is
available.

# Installation

Netwait is provided as platform-specific binaries and a Docker image.

For platform binaries, see [releases](https://github.com/guidewire/netwait/releases)

Docker image: [gwre/netwait](https://hub.docker.com/r/gwre/netwait)

# Usage

```bash
# Wait for resource to become available
$ netwait https://github.com
available: https://github.com
```

Netwait returns a non-zero exit code if it was unable to successfully connect
to the network resource within a certain time period.

```bash
$ netwait http://httpbin.org/status/200 && echo 'SUCCESS'
available: http://httpbin.org/status/200
SUCCESS

$ if netwait bad-domain.fake:443; then echo 'SUCCESS'; else echo 'FAIL'; fi
unavailable: bad-domain.fake:443
Error: All attempts fail:
#1: dial tcp: lookup bad-domain.fake: no such host
#2: dial tcp: lookup bad-domain.fake: no such host
#3: context deadline exceeded
FAIL
```

## HTTP resource

Netwait can wait for an HTTP URL. Success is defined as an HTTP response with
a 2xx status code. HTTP Redirects are followed.

```bash
$ netwait http://httpbin.org/status/200
available: http://httpbin.org/status/200

$ netwait -t 10s http://httpbin.org/status/500
unavailable: http://httpbin.org/status/500
Error: All attempts fail:
#1: GET 'http://httpbin.org/status/500': returned status code 500
#2: GET 'http://httpbin.org/status/500': returned status code 500
#3: GET 'http://httpbin.org/status/500': returned status code 500
#4: context deadline exceeded
```

## TCP resource

```bash
$ netwait github.com:443
available: github.com:443
```

## Multiple resources

Netwait will attempt to connect to multiple resources in parallel. All
resources must successfully connect for the command to return success.

```bash
$ netwait https://github.com https://go.dev/
available: https://github.com
available: https://go.dev/
```

## Timeout

```bash
# stop waiting after up to 10 seconds
$ netwait --timeout 10s github.com:443
available: github.com:443

# stop waiting after up to 2 minutes
$ netwait --timeout 2m github.com:443
available: github.com:443
```

## Docker

Netwait is available as a Docker image.

```bash
$ docker run gwre/netwait https://github.com
available: https://github.com
```

## Kubernetes

One common use for this utility is to add a wait step when starting a service
that has a startup dependency on another service.

For example, service A depends on service B and needs to wait to start until
service B is available. By using netwait in service A before running the normal
startup process, the startup will be paused until service B is available.

The simplest way to achieve this is by adding an initContainer for netwait. This
is run before a Pod's normal containers are started.

For example, here service "alpha" is configured to wait for services "bravo" and
"charlie" using netwait in an initContainer.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: alpha
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: alpha
  template:
    metadata:
      labels:
        app: alpha
    spec:
      initContainers:
        - name: netwait
          image: gwre/netwait:0.2
          args:
            # wait for TCP service bravo and HTTP service charlie
            - bravo.some-namespace.svc:1234
            - http://charlie.some-namespace.svc:8080
      containers:
        # Main container for alpha service, using podinfo as an example.
        # Assume this service has a startup requirement for bravo and charlie.
        - name: podinfo
          image: ghcr.io/stefanprodan/podinfo:6.1.6
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          command:
            - ./podinfo
            - --port=8080
            - --level=debug
```

Alternatively, netwait can be installed into service A's image and the startup
process can be adjusted to call netwait before the primary startup.
