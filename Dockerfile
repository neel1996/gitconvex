FROM golang:1.15.5-alpine

WORKDIR /go/src/github.com/neel1996/gitconvex-server

COPY . .

RUN apk update && apk upgrade
RUN apk add --update nodejs nodejs-npm openssl libgit2 libssh2 pkgconfig

RUN cd ui/ && \
    npm install && \
    export NODE_ENV=production && \
    npm install tailwindcss postcss autoprefixer && \
    npx tailwindcss build -o src/index.css -c src/tailwind.config.js && \
    npm run build && \
    mv build/ gitconvex-ui/ && \
    mv gitconvex-ui/ ../ && \
    cd .. && \
    rm -rf ui/

RUN apk del nodejs nodejs-npm

EXPOSE 9001

CMD go run /go/src/github.com/neel1996/gitconvex-server/server.go
