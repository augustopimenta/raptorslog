FROM golang:1.12.9-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o truck-service


FROM scratch

COPY --from=builder /app/truck-service /app/truck-service

EXPOSE 8000

ENTRYPOINT [ "/app/truck-service" ]