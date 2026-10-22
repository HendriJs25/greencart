FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o user-service .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/user-service ./
#COPY --from=builder /app/migrations ./migrations

EXPOSE 8001

ENTRYPOINT ["./user-service"]
CMD ["serve"]