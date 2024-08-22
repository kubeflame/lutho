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

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./lutho

# =========== FIN ===========
FROM alpine:3.18
WORKDIR /app

COPY --from=backend /app/lutho ./lutho

ENTRYPOINT ["/app/lutho"]