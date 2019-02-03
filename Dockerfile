FROM golang:1.11.5-alpine as builder
COPY . /go/src/toml-dev/graphite-exporter
WORKDIR /go/src/toml-dev/graphite-exporter 
RUN apk update && apk add --no-cache git
RUN go build -o main .

FROM scratch
COPY --from=builder /go/src/toml-dev/graphite-exporter/main /app/
WORKDIR /app
CMD ["./main"]
