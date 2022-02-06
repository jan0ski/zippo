FROM golang:1.17-alpine AS builder
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY zippo /
ENTRYPOINT ["/zippo"]
CMD ["/zippo"]
