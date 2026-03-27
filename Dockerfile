ARG BUILDPLATFORM=linux/amd64
ARG TARGETOS=linux
ARG TARGETARCH=amd64

FROM --platform=$BUILDPLATFORM node:24-alpine AS frontend-builder
WORKDIR /src/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS backend-builder
WORKDIR /src/backend

ARG TARGETOS
ARG TARGETARCH

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/blog .

FROM alpine:3.22
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=backend-builder /out/blog /app/blog
COPY --from=frontend-builder /src/frontend/dist /app/frontend-dist

ENV BLOG_ADDR=:8080
ENV BLOG_DB_PATH=/app/storage/blog.db
ENV FRONTEND_DIST=/app/frontend-dist
ENV BLOG_ADMIN_USER=admin
ENV BLOG_ADMIN_COOKIE_NAME=blog_admin_session
ENV BLOG_ADMIN_COOKIE_SECURE=false
ENV BLOG_ADMIN_SESSION_HOURS=72

VOLUME ["/app/storage"]
EXPOSE 8080

ENTRYPOINT ["./blog"]
