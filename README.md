# 个人博客

一个本地优先的全栈个人博客，技术栈为 `Golang + React`。

## 当前实现

- 后端使用 Go，按 `api/controller/service/dao/model` 分层。
- 前端使用 React + Vite，包含首页、文章详情、项目区和时间线。
- 博客内容默认存入 SQLite，本地启动会自动建表并初始化示例数据。
- 游客可以直接评论和点赞。
- 游客身份由 `IP + User-Agent` 生成稳定匿名指纹，展示为固定匿名名，例如 `访客-AB12CD34`。
- 已包含后台管理页，可登录后新增、编辑、删除文章。

## 目录结构

```text
backend/
  api/
  controller/
  dao/
  middleware/
  model/
  router/
  service/
frontend/
  src/
```

## 本地启动

前提：

- 已安装 Go 1.22+
- 已安装 Node.js 18+

### 1. 启动后端

```bash
cd /Users/fitz/personal/blog/backend
go mod tidy
go run .
```

可选环境变量：

- `BLOG_ADDR`：监听地址，默认 `:8080`
- `BLOG_DB_PATH`：SQLite 文件路径，默认 `storage/blog.db`
- `FRONTEND_DIST`：前端构建产物目录；配置后后端可直接托管静态站点
- `BLOG_ADMIN_USER`：后台用户名，默认 `admin`
- `BLOG_ADMIN_PASSWORD`：后台密码，必填
- `BLOG_ADMIN_SECRET`：后台 Cookie 签名密钥，必填
- `BLOG_ADMIN_COOKIE_NAME`：后台 Cookie 名称，默认 `blog_admin_session`
- `BLOG_ADMIN_COOKIE_SECURE`：是否只在 HTTPS 下发送后台 Cookie，默认 `false`
- `BLOG_ADMIN_SESSION_HOURS`：后台登录时长，默认 `72`

建议先在仓库根目录创建 `.env`：

```bash
cd /Users/fitz/personal/blog
cp .env.example .env
```

### 2. 启动前端

```bash
cd /Users/fitz/personal/blog/frontend
npm install
npm run dev
```

如果后端不在 `http://localhost:8080`，可以在前端启动前指定：

```bash
VITE_API_BASE=http://localhost:8080 npm run dev
```

默认情况下，Vite 开发服务器已把 `/api` 代理到 `http://localhost:8080`，所以本地一般不需要额外配置 `VITE_API_BASE`。

## 后台入口

- 登录页：`/admin/login`
- 后台首页：`/admin`

本地默认管理员：

- 用户名：`admin`

后台密码和密钥不再写死在仓库中，必须从 `.env` 或运行环境变量注入。

## Docker

已提供：

- [Dockerfile](/Users/fitz/personal/blog/Dockerfile)
- [docker-compose.yml](/Users/fitz/personal/blog/docker-compose.yml)
- [.dockerignore](/Users/fitz/personal/blog/.dockerignore)

直接构建镜像：

```bash
cd /Users/fitz/personal/blog
docker build -t personal-blog:latest .
```

直接运行容器：

```bash
docker run -d \
  --name personal-blog \
  -p 8080:8080 \
  -e BLOG_ADMIN_PASSWORD=your-strong-password \
  -e BLOG_ADMIN_SECRET=your-strong-secret \
  -v personal_blog_storage:/app/storage \
  personal-blog:latest
```

使用 Compose：

```bash
cd /Users/fitz/personal/blog
cp .env.example .env
docker compose up -d --build
```

说明：

- 容器内会由 Go 服务直接托管前端静态文件。
- SQLite 数据库存放在 `/app/storage/blog.db`。
- 生产环境建议把 `BLOG_ADMIN_COOKIE_SECURE` 设为 `true`，并放在 HTTPS 之后。
- `.env` 已加入 `.gitignore`，不要把真实密码和密钥提交到仓库。

## 生产部署建议

- 本地开发先使用 SQLite。
- 部署到云服务器时，博客文章、评论、点赞建议迁移到 MySQL。
- 当前 DAO 已经抽出仓储接口，后续新增 MySQL Repository 即可平滑替换。
- 如果前端执行 `npm run build`，可以把 `frontend/dist` 作为 `FRONTEND_DIST` 交给后端统一托管。

## 当前限制

- 这台机器当前没有 `go`、`node`、`npm` 命令，当前提交只能做静态层面的代码检查，无法直接在本地执行构建和测试。
- 游客匿名名依赖请求链路中的 `IP + User-Agent`，当用户换网络、换浏览器或经过代理时，匿名名可能变化。
