FROM golang:1.12.1-alpine as builder
COPY . /go/src/toml-dev/graphite-exporter
WORKDIR /go/src/toml-dev/graphite-exporter 
RUN apk update && apk add --no-cache git
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /go/src/toml-dev/graphite-exporter/main /app/
WORKDIR /app
CMD ["./main"]
