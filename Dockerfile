FROM golang:alpine AS builder

RUN go version
ENV GOPATH=/

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

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

COPY --from=builder /build/app /app/app
COPY --from=builder /build/migrations/*.sql /app/migrations

CMD ["./app"]