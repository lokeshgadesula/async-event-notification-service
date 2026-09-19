FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o /worker ./cmd/worker
FROM alpine:3.20
COPY --from=build /worker /worker
CMD ["/worker"]
