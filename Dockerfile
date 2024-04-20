FROM library/golang:1.22-alpine as build

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ARG version=unspecified
RUN go build -v -ldflags "-X 'github.com/guidewire/netwait/cmd.version=$version'"

FROM library/alpine:3.17

WORKDIR /app
ENV PATH=$PATH:/app

COPY --from=build /build/netwait .

ENTRYPOINT ["/app/netwait"]
