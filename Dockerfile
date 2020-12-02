FROM golang:1.15.5-alpine

WORKDIR /go/src/github.com/neel1996/gitconvex-server

COPY . .

RUN apk update && apk upgrade
RUN apk add --update nodejs nodejs-npm
RUN apk add git

RUN cd ui/ && \
    npm install && \
    export NODE_ENV=production && \
    npx tailwindcss build -o src/index.css -c src/tailwind.config.js && \
    npm run build && \
    mv build/ gitconvex-ui/ && \
    mv gitconvex-ui/ ../  

EXPOSE 9001

CMD go run /go/src/github.com/neel1996/gitconvex-server/server.go
