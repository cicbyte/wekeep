# ====== 构建前端 ======
FROM node:22-alpine AS frontend
WORKDIR /build/web
COPY web/package.json ./
RUN npm i
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
RUN go install github.com/gogf/gf/cmd/gf/v2@latest
RUN gf build -n wekeep -v $(cat VERSION) -p /build/bin -a amd64 -s linux \
    && mv /build/bin/$(cat VERSION)/linux_amd64/wekeep /build/wekeep

# ====== 运行镜像 ======
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=backend /build/wekeep /app/wekeep

ENV TZ=Asia/Shanghai
EXPOSE 8000

ENTRYPOINT ["/app/wekeep"]
