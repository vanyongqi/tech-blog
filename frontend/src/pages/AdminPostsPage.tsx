import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useAdminSessionContext } from "../components/AdminShell";
import { deleteAdminPost, formatDate, getAdminPosts } from "../lib/api";
import type { AdminPostSummaryPayload } from "../types";

export function AdminPostsPage() {
  useAdminSessionContext();
  const [posts, setPosts] = useState<AdminPostSummaryPayload[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const response = await getAdminPosts();
        if (!cancelled) {
          setPosts(response.posts);
          setError("");
        }
      } catch (loadError) {
        if (!cancelled) {
          setError(loadError instanceof Error ? loadError.message : "加载失败");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void load();
    return () => {
      cancelled = true;
    };
  }, []);

  async function handleDelete(slug: string) {
    const confirmed = window.confirm("确认删除这篇文章吗？删除后评论和点赞会一并删除。");
    if (!confirmed) {
      return;
    }

    try {
      await deleteAdminPost(slug);
      setPosts((currentPosts) => currentPosts.filter((post) => post.slug !== slug));
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "删除失败");
    }
  }

  if (loading) {
    return <section className="admin-panel">文章列表加载中...</section>;
  }

  return (
    <section className="admin-panel">
      <div className="admin-panel-heading">
        <div>
          <p className="eyebrow">Posts</p>
          <h2>文章管理</h2>
        </div>
        <Link className="primary-link" to="/admin/posts/new">
          新建文章
        </Link>
      </div>

      {error ? <p className="form-error">{error}</p> : null}

      <div className="admin-stats-grid">
        <article className="admin-stat-card">
          <span>文章总数</span>
          <strong>{posts.length}</strong>
        </article>
        <article className="admin-stat-card">
          <span>精选文章</span>
          <strong>{posts.filter((post) => post.featured).length}</strong>
        </article>
        <article className="admin-stat-card">
          <span>累计评论</span>
          <strong>{posts.reduce((sum, post) => sum + post.commentCount, 0)}</strong>
        </article>
      </div>

      <div className="admin-post-list">
        {posts.map((post) => (
          <article key={post.slug} className="admin-post-card">
            <div className="admin-post-card-top">
              <div>
                <p className="eyebrow">{post.coverLabel}</p>
                <h3>{post.title}</h3>
              </div>
              <span className={post.featured ? "status-pill active" : "status-pill"}>
                {post.featured ? "Featured" : "Normal"}
              </span>
            </div>
            <p>{post.summary}</p>
            <div className="admin-post-meta">
              <span>{post.category}</span>
              <span>{formatDate(post.publishedAt)}</span>
              <span>点赞 {post.likeCount}</span>
              <span>评论 {post.commentCount}</span>
            </div>
            <div className="admin-post-actions">
              <Link className="ghost-link" to={`/admin/posts/${post.slug}/edit`}>
                编辑
              </Link>
              <button type="button" className="danger-link" onClick={() => handleDelete(post.slug)}>
                删除
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
