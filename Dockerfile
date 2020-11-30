FROM golang:1.15.5-alpine

WORKDIR /opt/workroom/gitconvex

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

RUN cd /opt/workroom/gitconvex && \
    mkdir dist/ && \
    go build -o ./dist/ && \
    mv gitconvex-ui/ dist/ 

EXPOSE 9001

CMD go run /opt/workroom/gitconvex/dist/gitconvex-server
