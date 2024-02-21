FROM golang:1.16.5-alpine3.13 AS builder
WORKDIR /go/src/github.com/alexellis/href-counter/
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.13
WORKDIR /root/
COPY --from=builder /go/src/github.com/alexellis/href-counter/app .
CMD ["./app"]