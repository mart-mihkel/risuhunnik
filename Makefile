GO_BUILD_FLAGS = -ldflags='-s -w -extldflags "-static"'
GO_OUT = build/main

DB_FILE = risuhunnik.db

all: build

build: go tailwind database

go:
	mkdir -p build
	go build -o $(GO_OUT) $(GO_BUILD_FLAGS) cmd/main.go

tailwind:
	npm clean-install
	npm run tailwind

database:
	for file in sql/*; do sqlite3 $(DB_FILE) < $$file; done

clean:
	rm -rv $(GO_OUT) $(DB_FILE) node_modules
