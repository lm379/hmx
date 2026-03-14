# 构建前端
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

RUN npm install -g pnpm

COPY frontend/package.json frontend/pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY frontend/ ./

RUN pnpm run build

# 构建后端
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git sed

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY --from=frontend-builder /app/cmd/server/static ./cmd/server/static

RUN sed -i '/Logger:.*logger.Default.LogMode(logger.Info)/d' database/database.go && \
    sed -i '/\t"gorm.io\/gorm\/logger"/d' database/database.go

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o hmx \
    ./cmd/server/main.go

FROM golang:1.24-alpine

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=backend-builder /app/hmx .

COPY --from=backend-builder /app/cmd/server/static ./cmd/server/static

ENV TZ=Asia/Shanghai

EXPOSE 15379

ENTRYPOINT ["./hmx"]