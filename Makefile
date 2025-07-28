GO_BUILD_FLAGS = -ldflags='-s -w -extldflags "-static"'
GO_OUT = build/main

DB_FILE = risuhunnik.db

all: build

build: go database

go:
	mkdir -p build
	go build -o $(GO_OUT) $(GO_BUILD_FLAGS) cmd/main.go

database:
	for file in sql/*; do sqlite3 $(DB_FILE) < $$file; done

clean:
	rm -rv $(GO_OUT) $(DB_FILE)
