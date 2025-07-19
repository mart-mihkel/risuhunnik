FROM golang:alpine AS build

ENV CGO_ENABLED=1
RUN apk add --no-cache sqlite gcc musl-dev
WORKDIR /workspace
COPY . /workspace/
RUN go build -o main -ldflags='-s -w -extldflags "-static"' main.go
RUN sqlite3 risuhunnik.db < dump.sql

FROM scratch

WORKDIR /app
COPY ./static /app/static
COPY ./templates /app/templates
COPY --from=build /workspace/risuhunnik.db /app/risuhunnik.db
COPY --from=build /workspace/main /app/main

ENTRYPOINT [ "/app/main" ]
