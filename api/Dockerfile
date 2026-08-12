FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy the shared package and the api service source code
COPY shared/ ./shared/
COPY api/ ./api/
COPY frontend/ ./frontend/

# Build the API binary
WORKDIR /build/api
RUN go build -o /app .


# Run Stage
FROM alpine:3.21 AS final

COPY --from=builder /app /bin/app

COPY --from=builder /build/frontend /frontend

EXPOSE 8000

CMD ["/bin/app"]