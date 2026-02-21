# ---- Build Stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/rate-limiter ./cmd/server

# ---- Runtime Stage ----
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /bin/rate-limiter /rate-limiter
COPY configs/config.yaml /configs/config.yaml

EXPOSE 8080

ENTRYPOINT ["/rate-limiter"]
