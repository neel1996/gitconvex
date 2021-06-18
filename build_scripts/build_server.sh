#!/bin/bash

if [[ "$GITCONVEX_DEFAULT_PATH" == "" ]]; then
  echo "Default path missing";
  echo "Please set the default path where gitconvex can maintain its data";
  echo "========================================="
  echo "export GITCONVEX_DEFAULT_PATH=/some/path"
  echo "========================================="
  exit 1;
fi

echo "🚀 Building server modules"
go generate
go build -v -o ./dist