FROM node:12.18.1

WORKDIR /opt/gitconvex-docker

COPY package*.json ./

RUN npm install 

COPY . .

EXPOSE 9001

RUN npm i -g pm2

CMD pm2-runtime ecosystem.config.js
