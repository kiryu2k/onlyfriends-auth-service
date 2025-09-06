FROM golang:alpine AS builder

RUN go version
ENV GOPATH=/

WORKDIR /build

ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.40/grpc_health_probe-linux-amd64 /build/grpc_health_probe
RUN chmod +x /build/grpc_health_probe

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine AS runner

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

WORKDIR /app

USER $USERNAME

COPY --from=builder /build/main /app/main
COPY --from=builder /build/migrations/*.sql /app/migrations
COPY --from=builder /build/grpc_health_probe /app/grpc_health_probe

CMD ["./main"]