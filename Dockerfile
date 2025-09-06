# --- Build stage ---
FROM golang:1.24.6-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# --- Runtime stage ---
FROM debian:stable-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
      ffmpeg yt-dlp ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY web /app/web
RUN mkdir -p /app/data/temp
ENV PORT=8080 DATABASE_PATH=/app/data/songs.db TEMP_DIR=/app/data/temp
VOLUME ["/app/data"]
EXPOSE 8080
CMD ["/app/server"]
