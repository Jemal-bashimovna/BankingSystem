FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/myapp /myapp

ENV PORT=8080

EXPOSE 8080

CMD ["/myapp"]
