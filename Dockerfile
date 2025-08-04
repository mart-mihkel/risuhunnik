FROM golang:alpine AS go

ENV CGO_ENABLED=1
RUN apk add --no-cache gcc make musl-dev sqlite bash

WORKDIR /app
COPY . /app

RUN --mount=type=cache,target=/root/.cache/go-build make build



FROM scratch AS risuhunnik

WORKDIR /app
COPY . /app
COPY --from=go /app/build/risuhunnik /app/risuhunnik

ENTRYPOINT [ "/app/risuhunnik" ]
