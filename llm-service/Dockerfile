FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependency download unless go.mod changes
COPY go.mod go.sum ./
RUN go mod download

# Rebuild only app layer when source changes
COPY . .
RUN go build -o /main .

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /main ./main

CMD ["./main"]