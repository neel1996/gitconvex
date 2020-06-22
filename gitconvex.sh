#!/bin/bash

pm2 start serve.js --name gitconvex-ui

pm2 start server.js --name gitconvex-server
