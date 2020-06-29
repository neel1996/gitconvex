## Dockerfile for gitconvex

# If you are running the container from windows, make sure you map the host volume 
# ... while running container using the docker run command

# docker container run -t -d -p 9001:9001 --name gitconvex -v HOST_VOLUME:CONTAINER_VOLUME itassistors/gitconvex:TAG_VERSION

FROM node:12.18.1-alpine3.12

WORKDIR /opt/workroom/gitconvex-package

COPY package*.json ./

RUN npm install 

COPY . .

EXPOSE 9001

RUN npm i -g pm2
RUN apk update && apk add git

CMD pm2-runtime ecosystem.config.js
