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
    mv ./build ../ 

RUN cd /opt/workroom/gitconvex
RUN mkdir dist/ 
RUN go build -o ./dist/
RUN mv build/ dist/

EXPOSE 9001

# CMD cd /opt/workroom/gitconvex/dist && ./gitconvex-server
CMD go run server.go
