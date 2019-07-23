FROM golang:alpine AS builder
RUN echo $GOPATH
RUN apk add build-base --no-cache --update-cache --repository https://mirrors.aliyun.com/alpine/latest-stable/main/ --allow-untrusted
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
RUN adduser -D -g 'root' appuser
WORKDIR /go/src/github.com/flnow/plugin-http
COPY . .
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o /bin/plugin-http

FROM alpine
# RUN apk add tzdata --no-cache --update-cache --repository https://mirrors.aliyun.com/alpine/latest-stable/main/ --allow-untrusted && \
#     ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
#     echo "Asia/Shanghai" > /etc/timezone
WORKDIR /data/plugin-http
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/plugin-http .
RUN chown -R appuser:root /data/plugin-http
USER appuser
ENTRYPOINT [ "./plugin-http" ]
