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
    echo "Building UI bundle"
    set NODE_ENV=production
    npm install tailwindcss postcss autoprefixer
    npx tailwindcss build -o src/index.css -c src/tailwind.config.js
    npm run build
	echo "Moving react bundle to gitconvex-ui"
    move .\build gitconvex-ui
    move .\gitconvex-ui ..\
    cd ..
    mkdir .\dist
	echo "Moving UI artifacts to dist folder"
    move .\gitconvex-ui .\dist\
    echo "Moving etc content to dist"
    move .\etc .\dist\
	echo "Removing intermediary folder ui/"
	rd /s /q ui
    echo "Building gitconvex bundle"
    go build -a -o ./dist
	cd .\dist
    rename gitconvex-server.exe gitconvex.exe
	echo "Run ./dist/gitconvex.exe to start gitconvex on port 9001"
	echo "Try ./dist/gitconvex.exe --port PORT_NUMBER to run gitconvex on the desired port"
	cd ..
)

if "%ip%"==%TEST% (
    go test -tags static -v ./...
)

if "%ip%"==%RUN% (
    go run server.go
)

if "%ip%"==%START% (
	.\dist\gitconvex.exe
)
