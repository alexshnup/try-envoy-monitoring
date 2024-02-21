FROM klakegg/hugo:latest AS builder

COPY quickstart/. /src
RUN hugo

FROM alpine:latest
RUN apk add --no-cache darkhttpd
COPY --from=builder /src/public /var/www
CMD ["darkhttpd", "/var/www", "--port", "10000"]
