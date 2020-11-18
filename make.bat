@echo off
set ip=%1

set INSTALL="install"
set BUILD="build"
set RUN="run"
set TEST="test"
set START="start"

if "%ip%"==%INSTALL% (
    echo "Installing Go Dependencies..."
    go get
)

if "%ip%"==%BUILD% (
	echo "Initiating gitconvex build for windows"
	echo "Cleaning up unwanted folders"
	rd /s /q ui
	rd /s /q dist
	rd /s /q build
    echo "Cloning UI package from github gitconvex-ui/master"
    git clone https://github.com/neel1996/gitconvex-ui.git ui/
    cd ui
    echo "Installing UI dependencies"
	del package-lock.json
    npm install
    npm i -g tailwindcss@1.6.0
    echo "Building UI bundle"
    set NODE_ENV=production
    tailwindcss build -o src/index.css -c src/tailwind.config.js
    npm run build
    move .\build ..\
    cd ..
    mkdir .\dist
    move .\build .\dist\
	echo "Removing intermediary folder ui/"
	rd /s /q ui
    echo "Building gitconvex bundle"
    go build -o ./dist
	echo "Run ./dist/gitconvex-server.exe to start gitconvex on port 9001"
)

if "%ip%"==%TEST% (
    go test -v ./...
)

if "%ip%"==%RUN% (
    go run server.go
)

if "%ip%"==%START% (
	.\dist\gitconvex-server.exe
)
