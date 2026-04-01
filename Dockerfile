FROM golang:1.20 AS builder

WORKDIR /go/src/github.com/micromdm/micromdm/
COPY micromdm/ .
RUN make deps
RUN make

FROM alpine:latest
RUN apk --update add ca-certificates
COPY --from=builder /go/src/github.com/micromdm/micromdm/build/linux/micromdm /usr/bin/
COPY --from=builder /go/src/github.com/micromdm/micromdm/build/linux/mdmctl /usr/bin/

EXPOSE 8080
EXPOSE 443

ENTRYPOINT ["micromdm"]
CMD ["serve"]
