FROM registry.sensetime.com/lepton/ubuntu:22.04  as prod

# GO111MODULE=on GOSUMDB=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "./_output/bin/aks" "main.go"
WORKDIR /root/

RUN apt update && apt -y install sshpass

COPY ./_output/bin/aks /root/
ENTRYPOINT ["/root/aks"]