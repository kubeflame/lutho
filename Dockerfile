# syntax=docker/dockerfile:1

# =========== FE ===========
FROM node:alpine3.18 AS frontend
WORKDIR /app

COPY ./frontend/package.json ./
RUN yarn install

COPY ./frontend ./
RUN yarn run build

# =========== BE ===========
FROM golang:1.22 AS backend
WORKDIR /app

COPY --from=frontend /app/dist ./frontend/dist

COPY go.mod go.sum ./
RUN go mod download

ARG LUTHO_VERSION

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.AppVersion=$LUTHO_VERSION" -o ./lutho

# =========== FIN ===========
FROM alpine:3.18
WORKDIR /root

COPY --from=backend /app/lutho /usr/bin/lutho

ENTRYPOINT ["/usr/bin/lutho"]
