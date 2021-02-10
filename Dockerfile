FROM golang:1.15

RUN apt update && apt upgrade -y
RUN apt install -y \
git \
curl \
unzip

COPY --from=hashicorp/terraform:0.13.6 /bin/terraform /usr/local/bin/

# RUN curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
# RUN apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
# RUN apt-get update && sudo apt-get install terraform

# Install terraform
# RUN curl -sL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip | unzip -d /usr/local/bin -
RUN adduser \
--disabled-password \
--gecos "" \
json

WORKDIR /app
RUN chown -R json /app

COPY . .

USER json

RUN go get -d -v ./...

RUN go install -v ./...
