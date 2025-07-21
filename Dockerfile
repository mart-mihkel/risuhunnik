FROM golang:alpine AS go

ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . /app

RUN go build -o main -ldflags='-s -w -extldflags "-static"' cmd/main.go

FROM node:alpine AS node

WORKDIR /app
COPY . /app

RUN npm install
RUN npm run tailwind

FROM scratch AS app

WORKDIR /app
COPY . /app
COPY --from=node /app/css /app/css
COPY --from=go /app/main /app/main

ENTRYPOINT [ "/app/main" ]
