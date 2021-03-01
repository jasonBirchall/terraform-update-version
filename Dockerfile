FROM alpine:latest as gh
RUN apk --no-cache add wget tar
RUN wget https://github.com/cli/cli/releases/download/v0.5.7/gh_0.5.7_linux_amd64.tar.gz
RUN tar -zxvf gh_0.5.7_linux_amd64.tar.gz
RUN chmod a+x gh_0.5.7_linux_amd64/bin/gh

FROM golang:1.15

RUN apt update && apt upgrade -y
RUN apt install -y \
git \
curl \
unzip \
openssh-client

COPY --from=hashicorp/terraform:0.13.6 /bin/terraform /usr/local/bin/terraform0.13
COPY --from=hashicorp/terraform:0.14.7 /bin/terraform /usr/local/bin/terraform0.14
COPY --from=gh gh_0.5.7_linux_amd64/bin/gh /usr/bin/gh

RUN adduser \
--disabled-password \
--gecos "" \
json

WORKDIR /app
RUN chown -R json /app

COPY . .
RUN eval "$(ssh-agent -s)"

USER json

RUN go get -d -v ./...

RUN go install -v ./...

