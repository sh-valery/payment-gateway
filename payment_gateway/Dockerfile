# builder image
FROM golang:alpine as builder

WORKDIR /app

COPY .. .

RUN go build -o web_server ./cmd/web_server

# run image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/web_server .

EXPOSE 8080

# never separate the command and output redirection
CMD [ "/bin/sh", "-c", "/app/web_server > /var/log/gateway.log" ]


