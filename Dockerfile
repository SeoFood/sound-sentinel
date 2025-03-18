FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sound-sentinel ./cmd

FROM alpine:3.17

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY --from=builder /app/sound-sentinel .

CMD ["./sound-sentinel"]
