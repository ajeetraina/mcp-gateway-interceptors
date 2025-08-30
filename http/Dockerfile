FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init interceptor 2>/dev/null || true
RUN go build -o interceptor main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/interceptor /interceptor
EXPOSE 8080
CMD ["/interceptor"]
