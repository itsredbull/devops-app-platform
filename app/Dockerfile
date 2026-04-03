FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/uptime-api ./cmd/uptime-api

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/uptime-api /uptime-api
COPY --from=builder /src/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/uptime-api"]
