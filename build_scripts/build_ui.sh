#!/bin/bash

echo "⏬ Cloning gitconvex react repo"
git clone https://github.com/neel1996/gitconvex-ui.git ui/
cd ui

echo "⏳ Installing UI dependencies..."
npm install
export NODE_ENV=production
npm i -g create-react-app tailwindcss@1.6.0
npm run build:tailwind

echo "🔧 Building react UI bundle"
npm run build
mv ./build ../build
cd ..