get:
	go get
run:
	go run server.go
build-ui:
	git clone https://github.com/neel1996/gitconvex-ui.git ui/
	cd ui
	npm install
	npm i -g create-react-app tailwindcss@1.6.0
	npm run build:tailwind
	npm run build
	mv ./build ../
build-server:
	mkdir -p ./dist
	go build -o ./dist
build:
	echo "Initiating gitconvex build"
	echo "Cloning gitconvex react repo"
	git clone https://github.com/neel1996/gitconvex-ui.git ui/
	cd ui && \
	npm i -g tailwindcss@1.6.0 && \
	echo "Installing UI dependencies..."
	npm install && \
	export NODE_ENV=production && \
	echo "Generating production ready css"
	tailwindcss build -o ./src/index.css -c ./src/tailwind.config.js && \
	rm package-*.json && \
    	rm -rf .git/ && \
    	echo "Building react UI bundle"
	npm run build && \
	mv ./build ../ && \
	cd ..
	mkdir -p ./dist
	echo "Building final go source with UI bundle"
	go build -o ./dist
	echo "Gitconvex build completed!"
	echo "Use ./dist/gitconvex-server to start Gitconvex on port 9001"
test:
	go test -v ./...
start:
	./dist/gitconvex-server
