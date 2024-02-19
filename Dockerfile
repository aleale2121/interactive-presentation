# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN  CGO_ENABLED=0 go build -o main cmd/rest/main.go && chmod +x /app/main

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .

COPY db/migration ./db/migration
EXPOSE 8081
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
