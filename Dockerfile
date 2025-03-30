FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o restorio-auth-service ./cmd/restorio-auth


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/restorio-auth-service /app/

CMD ["./restorio-auth-service"]