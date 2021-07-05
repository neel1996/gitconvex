FROM golang:1.16.0-alpine

WORKDIR /home/gitconvex

COPY . .

# Install required packages
RUN apk update && \
    apk add --update libgit2-dev libssh2-dev gcc make musl-dev git


