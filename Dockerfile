FROM golang:latest AS builder

COPY . /app

WORKDIR /app

RUN go build -o main .

FROM debian:latest

COPY --from=builder /app/main /usr/local/bin/main
COPY --from=builder /app/internal/config/config.yaml /internal/config/config.yaml

ENV PORT 11011
EXPOSE 11011

CMD ["main"]