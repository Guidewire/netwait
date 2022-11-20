FROM library/golang:1.18-alpine as build

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -v

FROM library/alpine:3.17

WORKDIR /app
ENV PATH=$PATH:/app

COPY --from=build /build/netwaiter .

ENTRYPOINT ["/app/netwaiter"]
