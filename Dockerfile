## docker build -t registry.cn-hangzhou.aliyuncs.com/generals-kuber/kube-scheduler-extender:v0.0.1 .

FROM golang:1.13 as builder
WORKDIR /app
COPY . .
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn
RUN go build -o kube-scheduler-extender main.go

########################################################################

FROM registry.cn-hangzhou.aliyuncs.com/generals-space/centos7

COPY --from=builder /app/kube-scheduler-extender /
RUN chmod 755 /kube-scheduler-extender
