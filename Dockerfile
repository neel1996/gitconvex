FROM golang:1.16.0-alpine

WORKDIR /go/src/github.com/neel1996/gitconvex-server

COPY . .

# Install required packages
RUN apk update && \
    apk add --update libgit2-dev libssh2-dev gcc make nodejs npm musl-dev

# Building React UI bundle
RUN cd ui/ && \
    npm install && \
    export NODE_ENV=production && \
    npm i -g npm@6 && \
    npm install tailwindcss postcss autoprefixer && \
    npx tailwindcss build -o src/index.css -c src/tailwind.config.js && \
    npm run build && \
    mv build/ gitconvex-ui/ && \
    mv gitconvex-ui/ ../ && \
    cd .. && \
    rm -rf ui/

RUN go get -v

EXPOSE 9001

CMD go run /go/src/github.com/neel1996/gitconvex-server/server.go