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

git-testuser-setup:
	git config --global user.name "test"
	git config --global user.email "test@test.com"

test:
	export GITCONVEX_TEST_REPO="$(PWD)/gitconvex-test" && \
	export GITCONVEX_DEFAULT_PATH="$(PWD)/gitconvex-test" && \
	rm -rf $$GITCONVEX_TEST_REPO && \
 	go clean --cache && \
 	./build_scripts/remove_test_repo.sh && \
 	./build_scripts/clone_test_repo.sh && \
 	go test ./... -count=1 -cover -coverprofile=coverage.out && \
	rm -rf $$GITCONVEX_TEST_REPO

test-pretty:
	export GITCONVEX_TEST_REPO="$(PWD)/gitconvex-test" && \
	export GITCONVEX_DEFAULT_PATH="$(PWD)/gitconvex-test" && \
	rm -rf $$GITCONVEX_TEST_REPO && \
 	go clean --cache && \
 	./build_scripts/remove_test_repo.sh && \
 	./build_scripts/clone_test_repo.sh && \
 	gotestsum ./... -count=1 -p=1 -cover -coverprofile=coverage.out && \
	rm -rf $$GITCONVEX_TEST_REPO

test-ci: git-testuser-setup
	go generate && \
	go clean --cache && \
    sh ./build_scripts/clone_test_repo.sh && \
    go test ./... -count=1 -p=1 -cover -coverprofile=coverage.out && \
    rm -rf $$GITCONVEX_TEST_REPO

test-ci-pretty: git-testuser-setup
	go generate && \
	go clean --cache && \
    sh ./build_scripts/clone_test_repo.sh && \
    go get gotest.tools/gotestsum && \
    gotestsum ./... -count=1 -p=1 -cover -coverprofile=coverage.out && \
    rm -rf $$GITCONVEX_TEST_REPO

dockerise-test:
	rm -rf $$GITCONVEX_TEST_REPO && \
	docker-compose -f docker-compose.test.yaml build && \
	docker-compose -f docker-compose.test.yaml up

check-coverage:
	sh ./build_scripts/check_coverage.sh

show-coverage:
	go tool cover -html=coverage.out

start:
	./dist/gitconvex
