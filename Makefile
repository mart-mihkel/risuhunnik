BUILD_DIR = build

GO_BUILD_FLAGS = -ldflags='-s -w -extldflags "-static"'
GO_OUT = $(BUILD_DIR)/risuhunnik

DB_FILE = $(BUILD_DIR)/risuhunnik.db

build: go db

dev: db
	go run cmd/main.go

go:
	mkdir -p build
	go build -o $(GO_OUT) $(GO_BUILD_FLAGS) cmd/main.go

db:
	mkdir -p build
	sqlite3 $(DB_FILE) < sql/schema.sql

clean:
	rm -rv $(BUILD_DIR)
