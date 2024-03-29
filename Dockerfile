FROM golang:1.17-alpine3.15 as builder
RUN apk add --no-cache git
RUN go get -u -v github.com/kyxap1/pwgen && \
  rm -rf /go/src
RUN apk del git

FROM alpine:3.15 as base
RUN adduser -S -h /app pwgen
COPY --from=builder /go/bin/pwgen /app/pwgen
RUN chown -R pwgen: /app
USER pwgen
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/app/pwgen"]
