@echo off
set ip=%1

set INSTALL="install"
set BUILD="build"
set RUN="run"
set TEST="test"

if "%ip%"==%INSTALL% (
    echo "Installing Go Dependencies..."
    go get
)

if "%ip%"==%BUILD% (
    echo "Cloning UI package from github"
    git clone https://github.com/neel1996/gitconvex-ui.git ui/
    cd ui
    echo "Installing UI dependencies"
	npm install
	npm i -g create-react-app tailwindcss@1.6.0
	echo "Building UI bundle"
	npm run build:tailwind
	npm run build
	mv ./build ../
	cd ..
	mkdir -p ./dist
	echo "Building gitconvex bundle"
	go build -o ./dist
)

if "%ip%"==%TEST% (
    go test -v ./...
)