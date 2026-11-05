build_dir := "build"
go_out := build_dir / "risuhunnik"
go_build_flags := "-ldflags='-s -w -extldflags \"-static\"'"

# show commands
help:
    @just --list

# build release binary
build:
    @mkdir -p {{ build_dir }}
    go build -o {{ go_out }} {{ go_build_flags }} cmd/main.go

# run migrations
migrate:
    @mkdir -p {{ build_dir }}
    ./sql/run-migrations.sh

# run dev server
dev: migrate
    go run cmd/main.go

# deploy docker
docker: migrate
    docker rm -f risuhunnik
    docker build -t risuhunnik .
    docker run -d \
        -v ./build:/app/build \
        -p 8080:8080 \
        --name risuhunnik \
        --restart always \
        risuhunnik

# format and lint
check:
    go fmt ./...
    go vet ./...

# run tests
test:
    go test ./...
