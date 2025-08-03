BUILD_DIR = build

GO_BUILD_FLAGS = -ldflags='-s -w -extldflags "-static"'
GO_OUT = $(BUILD_DIR)/risuhunnik

build: go db

dev: db
	go run cmd/main.go

go:
	mkdir -p build
	go build -o $(GO_OUT) $(GO_BUILD_FLAGS) cmd/main.go

db:
	mkdir -p build
	./sql/run-migrations.sh
