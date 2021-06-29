get:
	go get

run:
	go run server.go

build-ui:
	@./build_scripts/build_ui.sh

build-server:
	@./build_scripts/build_server.sh

clean:
	@echo "💨 Cleaning up previous build and bundles..."
	@rm -rf ui/ dist/ gitconvex build/

bundle:
	@mkdir -p ./dist && \
    mv ./build/ ./dist/ && \
    mv ./dist/build ./dist/gitconvex-ui;

pack:
	cd ./dist && \
	tar -cvzf gitconvex-$$GITCONVEX_VERSION.tar.gz .

build: clean build-ui bundle build-server
	@echo "✅ Gitconvex Build Completed successfully!"
	@echo "📬 Use ./dist/gitconvex to start Gitconvex on port 9001"
	@echo "📬 Try ./dist/gitconvex --port PORT_NUMBER to run gitconvex on the desired port"

test:
	export GITCONVEX_TEST_REPO="$(PWD)/gitconvex-test" && \
	export GITCONVEX_DEFAULT_PATH="$(PWD)/gitconvex-test" && \
	rm -rf $$GITCONVEX_TEST_REPO && \
 	go clean --cache && \
 	./build_scripts/clone_test_repo.sh && \
 	go test ./... -count=1 && \
	rm -rf $$GITCONVEX_TEST_REPO

start:
	./dist/gitconvex
