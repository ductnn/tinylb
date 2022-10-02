FROM golang:1.16 AS builder
LABEL maintainer="ductnn"

WORKDIR /app
COPY main.go go.mod ./

RUN CGO_ENABLED=0 GOOS=linux go build -o tinylb .


FROM alpine:latest
LABEL maintainer="ductnn"

RUN apk --no-cache add ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /root
COPY --from=builder /app/tinylb .

ENTRYPOINT [ "/root/tinylb" ]
