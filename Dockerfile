FROM golang:1.14 as goBuilder
WORKDIR /go_build
COPY . /go_build
ENV GOPROXY=https://goproxy.cn,direct
RUN go build -o /go/bin/app /go_build/main.go

# debian release as the same as golang image
# set TimeZone as Asia/Shanghai
# set Local as zh-hans
FROM hjin/app:stretch
ENV LOG_MODE production
ENV DICT_DIR /var/assets/dict
EXPOSE 8080
COPY --from=goBuilder /go/bin/app /usr/local/bin/app
COPY ./assets /var/assets
CMD ["app"]
