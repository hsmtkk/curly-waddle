FROM golang:1.18 AS builder
WORKDIR /opt
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build

FROM gcr.io/distroless/cc-debian11 AS runtime
COPY --from=builder /opt/curly-waddle /usr/local/bin/curly-waddle
ENTRYPOINT ["/usr/local/bin/curly-waddle"]