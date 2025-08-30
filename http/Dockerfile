FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go mod init interceptor && go build -o interceptor main.go

FROM alpine:latest
COPY --from=builder /app/interceptor /interceptor
CMD ["/interceptor"]
