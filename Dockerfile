FROM golang:1.23 AS builder

WORKDIR /auf-citaty

COPY ./ ./

RUN go mod download
RUN go build -o auf-citaty ./main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /auf-citaty

COPY --from=builder /auf-citaty/auf-citaty /auf-citaty/db.sqlite /auf-citaty/.env ./

EXPOSE 8080

CMD ["./auf-citaty"]

