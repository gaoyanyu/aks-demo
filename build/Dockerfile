#FROM golang:alpine3.14 as builder
FROM registry.sensetime.com/sensecore-lepton/golang:1.20-alpine as builder

WORKDIR /root/
COPY . .

ENV GO111MODULE=on
ENV GOSUMDB=off

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "aks" "main.go"

FROM registry.sensetime.com/lepton/ubuntu:22.04  as prod

WORKDIR /root/

RUN apt update && apt -y install sshpass

COPY --from=0 /root/aks /root/
ENTRYPOINT ["/root/aks"]