FROM golang:alpine AS build

ENV CGO_ENABLED=1
RUN apk add --no-cache sqlite gcc musl-dev
WORKDIR /workspace
COPY . /workspace/
RUN go build -o main -ldflags='-s -w -extldflags "-static"' cmd/main.go
RUN sqlite3 risuhunnik.db < sql/dump.sql

FROM scratch

WORKDIR /workspace
COPY ./static /workspace/static
COPY ./templates /workspace/templates
COPY --from=build /workspace/risuhunnik.db /workspace/sql/risuhunnik.db
COPY --from=build /workspace/main /workspace/main

ENTRYPOINT [ "/workspace/main" ]
