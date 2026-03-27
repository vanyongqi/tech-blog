# 自动部署说明

当前线上实例使用的是“GitHub Actions 发布 GHCR 镜像 + 云服务器定时轮询 GHCR”的方式。

实际链路：

1. 推送到 `main`
2. GitHub Actions 执行后端测试和前端构建
3. 构建并推送镜像到 `ghcr.io/vanyongqi/tech-blog:latest`
4. 云服务器定时执行部署脚本，拉取 GHCR 最新镜像
5. 如果镜像 digest 发生变化，则重建 `tech-blog` 容器

## 线上当前约定

- 部署目录：`/opt/tech-blog`
- 数据目录：`/opt/tech-blog/storage`
- 运行时环境文件：`/etc/tech-blog/blog.env`
- 线上容器名：`tech-blog`
- 监听方式：容器 `8080` 绑定到宿主机 `127.0.0.1:18080`
- Nginx 再把公网 `80/443` 转发到 `127.0.0.1:18080`

## 服务器部署脚本

服务器当前使用的脚本可以参考 [deploy/blog/deploy-ghcr.sh](/Users/fitz/personal/blog/deploy/blog/deploy-ghcr.sh)。

建议把脚本放到：

```bash
/opt/tech-blog/deploy-ghcr.sh
```

## 定时任务

服务器当前采用 root crontab 定时拉取：

```cron
* * * * * flock -n /tmp/tech-blog-deploy.lock /opt/tech-blog/deploy-ghcr.sh >> /var/log/tech-blog-deploy.log 2>&1
```

这样可以保证：

- 最多 1 分钟内感知到 GHCR 新镜像
- 上一轮部署未结束时，不会并发执行下一轮

## 运行时环境

运行时环境变量不放仓库，统一放在：

```bash
/etc/tech-blog/blog.env
```

至少应包含：

- `BLOG_ADMIN_PASSWORD`
- `BLOG_ADMIN_SECRET`
- `BLOG_ADMIN_COOKIE_SECURE=true`

## 注意事项

- 这套方案依赖服务器可以直接拉取 `ghcr.io/vanyongqi/tech-blog:latest`
- 如果 GHCR 镜像是私有的，需要先在服务器执行 `docker login ghcr.io`
- 如果以后不想轮询，可以再改成 GitHub Actions 主动 SSH 触发部署
