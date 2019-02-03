FROM golang:1.11.5-alpine as builder
RUN mkdir /build 
COPY . /build/
COPY . $GOPATH/src/toml-dev/graphite-exporter
WORKDIR $GOPATH/src/toml-dev/graphite-exporter 
RUN apk update && apk add --no-cache git
RUN go build -o main .

FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
