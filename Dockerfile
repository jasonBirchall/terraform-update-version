FROM golang:1.15
RUN apt install git

COPY . .
