BUILD_DIR = build

GO_BUILD_FLAGS = -ldflags='-s -w -extldflags "-static"'
GO_OUT = $(BUILD_DIR)/risuhunnik

build: go db

dev: db
	go run cmd/main.go

go:
	mkdir -p $(BUILD_DIR)
	go build -o $(GO_OUT) $(GO_BUILD_FLAGS) cmd/main.go

db:
	mkdir -p $(BUILD_DIR)
	./sql/run-migrations.sh

docker:
	docker rm -f risuhunnik 
	docker build -t risuhunnik .
	docker run -d \
		-v ./build:/app/build \
		-p 8080:8080 \
		--name risuhunnik \
		--restart always \
		risuhunnik

format:
	go fmt ./...

test:
	go test ./...
