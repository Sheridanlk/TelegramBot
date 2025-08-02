FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o main ./cmd/main.go


FROM gcr.io/distroless/base-debian12

WORKDIR /app
WORKDIR /app/data/sqlite
WORKDIR /app

COPY --from=build /src/main ./main
COPY --from=build /src/config ./config


ENV CONFIG_PATH=/app/config/cfg.yaml
ENV DATABASE_PATH=/app/data/sqlite/storage.db

CMD ["./main"]


