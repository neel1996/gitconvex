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
	git clone https://github.com/neel1996/gitconvex-ui.git ui/
	cd ui && \
	npm i -g tailwindcss@1.6.0 && \
	npm install && \
	export NODE_ENV=production && \
	tailwindcss build -o ./src/index.css -c ./src/tailwind.config.js && \
	rm package-*.json && \
    rm -rf .git/ && \
	npm run build && \
	mv ./build ../ && \
	cd ..
	mkdir -p ./dist
	go build -o ./dist
test:
	go test -v ./...