BUILD_DIR = build

GO_BUILD_FLAGS = -ldflags='-s -w -extldflags "-static"'
GO_OUT = $(BUILD_DIR)/risuhunnik

.PHONY: help
# show avalable commands
help:
	@printf "\nrisuhunnik development commands:\n\n"
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{printf "  \033[36m%-20s\033[0m %s\n", substr($$1,1,index($$1,":")-1),c}1{c=0}' $(MAKEFILE_LIST)
	@printf "\n"

.PHONY: build
# build release binary
build: database
	mkdir -p $(BUILD_DIR)
	go build -o $(GO_OUT) $(GO_BUILD_FLAGS) cmd/main.go

.PHONY: database
# run database migrations
database:
	mkdir -p $(BUILD_DIR)
	./sql/run-migrations.sh

.PHONY: dev
# start development server
dev: database
	go run cmd/main.go

.PHONY: docker
# start dockerized release server
docker:
	docker rm -f risuhunnik 
	docker build -t risuhunnik .
	docker run -d \
		-v ./build:/app/build \
		-p 8080:8080 \
		--name risuhunnik \
		--restart always \
		risuhunnik

.PHONY: goformat
# format go files
goformat:
	go fmt ./...

.PHONY: test
# run go tests
test:
	go test ./...
