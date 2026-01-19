# Multi-stage Dockerfile
# Build frontend
FROM node:18-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci --silent
COPY web/ ./
RUN npm run build

# Build backend
FROM golang:1.24-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy built frontend into the repo so go:embed can see it
COPY --from=frontend /app/web/dist ./web/dist
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /tamalabs ./cmd/rest

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder /tamalabs /tamalabs
ENV PORT=8322
EXPOSE 8322
CMD ["/tamalabs"]
