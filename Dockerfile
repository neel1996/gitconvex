FROM golang:1.15.5-alpine

WORKDIR /opt/workroom/gitconvex

COPY . .

RUN apk update && apk upgrade
RUN apk add --update nodejs nodejs-npm
RUN cd ui/ && \
    npm install && \
    npm i -g tailwindcss@1.6.0 && \
    export NODE_ENV=production && \
    tailwindcss build -o src/index.css -c src/tailwind.config.js && \
    npm run build && \
    cp ./build ../ 

RUN mkdir dist/ 
RUN go build -o ./dist/

EXPOSE 9001

CMD ./dist/gitconvex-server
