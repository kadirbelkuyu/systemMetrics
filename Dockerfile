FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main cmd/systemmetrics/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/cmd/systemmetrics/main.go /app/main.go
EXPOSE 3000
CMD ["./main"]