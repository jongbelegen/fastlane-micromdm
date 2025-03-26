FROM golang:1.20-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git
RUN git clone https://github.com/micromdm/micromdm.git .
RUN go build -o micromdm ./cmd/micromdm
RUN go build -o mdmctl ./cmd/mdmctl

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/micromdm /app/micromdm
COPY --from=builder /app/mdmctl /app/mdmctl

# Create directory for the database file
RUN mkdir -p /var/db/micromdm

EXPOSE 8080
EXPOSE 443

ENTRYPOINT ["/app/micromdm"]
CMD ["serve"] 