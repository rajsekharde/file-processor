FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy the shared package and the worker service source code
COPY shared/ ./shared/
COPY worker/ ./worker/

# Build the Worker binary
WORKDIR /build/worker
RUN go build -o /app .


# Run Stage
FROM alpine:3.21 AS final

COPY --from=builder /app /bin/app

EXPOSE 8001

CMD ["/bin/app"]