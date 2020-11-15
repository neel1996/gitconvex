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
	ifeq(,$(wildcard ./ui)){
		git clone https://github.com/neel1996/gitconvex-ui.git ui/
	}
	cd ui
	npm install
	npm i -g create-react-app tailwindcss@1.6.0
	npm run build:tailwind
	npm run build
	mv ./build ../
	cd ..
	mkdir -p ./dist
	go build -o ./dist
test:
	go test -v ./...