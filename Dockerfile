FROM golang:1.23 AS builder

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o app ./cmd/app/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/app /app/db.sqlite ./

EXPOSE 8080

CMD ["./app"]

