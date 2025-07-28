GO_BUILD_FLAGS = -ldflags='-s -w -extldflags "-static"'
GO_OUT = build/risuhunnik

DB_FILE = build/risuhunnik.db

build: go db

dev: build
	$(GO_OUT)

go:
	mkdir -p build
	go build -o $(GO_OUT) $(GO_BUILD_FLAGS) cmd/main.go

db:
	mkdir -p build
	sqlite3 $(DB_FILE) < sql/schema.sql

clean:
	rm -rv $(GO_OUT) $(DB_FILE)
