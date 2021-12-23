FROM golang:1.17-alpine AS builder

RUN apk --no-cache add ca-certificates
WORKDIR /build
COPY go.* ./
RUN go mod tidy
COPY ./ ./
RUN CGO_ENABLED=0 go build -o zippo -a -ldflags '-w' main.go

FROM scratch
COPY --from=builder /build/zippo /bin/zippo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/bin/zippo"]
CMD ["/bin/zippo"]
