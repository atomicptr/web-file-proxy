FROM alpine:3

ENV DATABASE_DRIVER "sqlite3"
ENV DATABASE_URL "/data/proxy.db"
ENV SERVICE_ADDR ":8081"

RUN apk add build-base go git

WORKDIR /app
COPY . /app

RUN go build

RUN mkdir -p /data

CMD ["./web-file-proxy"]