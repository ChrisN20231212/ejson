FROM golang:1.22 AS builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/ejson .

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /out/ejson /app/ejson

EXPOSE 8080

CMD ["/app/ejson"]
