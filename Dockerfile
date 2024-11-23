FROM golang:1.23 AS builder

WORKDIR /auf-citaty

COPY ./ ./

RUN go mod download
RUN go build -o auf-citaty ./cmd/app/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /auf-citaty/auf-citaty /auf-citaty/internal/database/db.sqlite ./

EXPOSE 8080

CMD ["./auf-citaty"]

