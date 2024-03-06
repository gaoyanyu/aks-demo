FROM golang:alpine3.14 as builder

WORKDIR /root/
COPY . .

ENV GO111MODULE=on
ENV GOSUMDB=off

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "aks" "main.go"

FROM registry.sensetime.com/lepton/ubuntu:22.04  as prod

#默认日志收集路径
ENV LOG_PATH="/root/logs"

WORKDIR /root/

COPY --from=0 /root/aks /root/
ENTRYPOINT ["/root/aks"]