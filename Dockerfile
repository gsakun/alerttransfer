FROM golang:1.13-alpine as builder

WORKDIR /go/src/github.com/gsakun/alerttransfer/

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux LINKFLAGS="-extldflags -static -s" go build -a -installsuffix cgo -o alerttransfer

FROM alpine:latest

WORKDIR /

RUN apk add --no-cache tzdata \
    && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

ENV TZ Asia/Shanghai

COPY --from=builder /go/src/github.com/gsakun/alerttransfer/alerttransfer /alerttransfer

ENTRYPOINT ["./alerttransfer"]