FROM golang:alpine3.18 AS build

ENV CGO_ENABLED=1

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /workspace

COPY . /workspace/

RUN \
    cd _example/simple && \
    go mod init github.com/mattn/sample && \
    go mod edit -replace=github.com/mattn/go-sqlite3=../.. && \
    go mod tidy && \
    go install -ldflags='-s -w -extldflags "-static"' ./simple.go

FROM scratch

COPY --from=build /go/bin/simple /usr/local/bin/simple

ENTRYPOINT [ "/usr/local/bin/simple" ]
