FROM golang:1.26-alpine AS build

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o /tmp/courier-api ./cmd/api

FROM alpine:3.22
WORKDIR /app
COPY --from=build /tmp/courier-api /app/courier-api
EXPOSE 8080
CMD ["/app/courier-api"]
