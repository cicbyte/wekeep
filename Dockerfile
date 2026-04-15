# ====== 构建前端 ======
FROM node:22-alpine AS frontend
WORKDIR /build/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# ====== 构建后端 ======
FROM golang:1.24-alpine AS backend

ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /build/web/dist/ resource/public/html/
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/ciclebyte/wekeep/internal/logic/health.Version=$(cat VERSION)" -o wekeep .

# ====== 运行镜像 ======
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend /build/wekeep /app/wekeep
COPY --from=backend /build/manifest /app/manifest

ENV TZ=Asia/Shanghai
EXPOSE 8000
VOLUME /data

ENTRYPOINT ["/app/wekeep"]
