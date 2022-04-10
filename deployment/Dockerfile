# 启动编译环境
FROM golang:1.17
# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
# 拷贝源代码到镜像中
COPY . /go/src/coolcar/server
# 编译
WORKDIR /go/src/coolcar/server
RUN go install ./gateway/..
# 设置服务入口
ENTRYPOINT ["go/bin/gateway"]