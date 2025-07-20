FROM golang:alpine AS build

ENV CGO_ENABLED=1
RUN apk add --no-cache sqlite gcc musl-dev
WORKDIR /app/go
COPY ./go /app/go
RUN sqlite3 risuhunnik.db < dump.sql
RUN go build -o main -ldflags='-s -w -extldflags "-static"' main.go

FROM node:alpine AS tailwind

WORKDIR /app/web
COPY ./web /app/web
RUN npm install
RUN npx @tailwindcss/cli -i ./style.css -o ./static/tailwind.css --minify

FROM scratch

WORKDIR /app/go
COPY ./web/templates /app/web/templates
COPY --from=tailwind /app/web/static/ /app/web/static
COPY --from=build /app/go/risuhunnik.db /app/go/risuhunnik.db
COPY --from=build /app/go/main /app/go/main

ENTRYPOINT [ "/app/go/main" ]
