get:
	go get
run:
	go run server.go
build-ui:
	git clone https://github.com/neel1996/gitconvex-ui.git ui/
	cd ui
	npm install
	export NODE_ENV=production
	npm i -g create-react-app tailwindcss@1.6.0
	npm run build:tailwind
	npm run build
	mv ./build ../
build-server:
	mkdir -p ./dist
	go build -o ./dist

install_deps:
	if [[ "$(shell uname)" == "Darwin" ]]; then cd ./lib/macos_amd64 && ./install_deps && cd .. ; else cd ./lib/linux_x86_64 && ./install_deps && cd ..; fi

build:
	echo "⚒️ Initiating gitconvex build"
	echo "🗑️ Cleaning up old directories"
	rm -rf ui/ dist/ build/
	echo "⏬ Cloning gitconvex react repo"
	git clone -q https://github.com/neel1996/gitconvex-ui.git ui/ && \
	cd ui && \
	echo "⏳ Installing UI dependencies..." && \
	npm install --silent && \
	export NODE_ENV=production && \
	npm install tailwindcss postcss autoprefixer && \
	npx tailwindcss build -o ./src/index.css -c ./src/tailwind.config.js && \
	rm package-*.json && \
	rm -rf .git/ && \
	echo "🔧 Building react UI bundle" && \
	npm run build && \
	mv ./build ../ && \
	cd .. && \
	mkdir -p ./dist && \
	mv build/ ./dist/ && \
	mv ./dist/build ./dist/gitconvex-ui
	echo "🚀 Building final go source with UI bundle" && \
	go build -v -a -o ./dist && \
	echo "Gitconvex build completed!" && \
	mv ./dist/gitconvex-server ./dist/gitconvex 
	echo "Installing libs"
	$(MAKE) install_deps
    echo "✅ Gitconvex Build Completed successfully!"
	echo "📬 Use ./dist/gitconvex to start Gitconvex on port 9001"
	echo "📬 Try ./dist/gitconvex --port PORT_NUMBER to run gitconvex on the desired port"
test:
	go test -tags static -v ./...
start:
	./dist/gitconvex
