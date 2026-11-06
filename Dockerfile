FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN go build -o task-app ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/task-app .

CMD ["./task-app"]