FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY kuksa ./kuksa
COPY *.go ./

RUN go build -o /rest-feeder

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /rest-feeder /rest-feeder

USER nonroot:nonroot

ENTRYPOINT ["/rest-feeder"]