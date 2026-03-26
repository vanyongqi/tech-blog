import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useAdminSessionContext } from "../components/AdminShell";
import { deleteAdminVideo, formatDate, getAdminVideos } from "../lib/api";
import type { AdminVideoPayload } from "../types";

export function AdminVideosPage() {
  useAdminSessionContext();
  const [videos, setVideos] = useState<AdminVideoPayload[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const response = await getAdminVideos();
        if (!cancelled) {
          setVideos(response.videos);
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

  async function handleDelete(id: number) {
    if (!window.confirm("确认删除这个视频吗？")) {
      return;
    }
    try {
      await deleteAdminVideo(id);
      setVideos((currentVideos) => currentVideos.filter((video) => video.id !== id));
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "删除失败");
    }
  }

  if (loading) {
    return <section className="admin-panel">视频列表加载中...</section>;
  }

  return (
    <section className="admin-panel">
      <div className="admin-panel-heading">
        <div>
          <p className="eyebrow">Videos</p>
          <h2>视频管理</h2>
        </div>
        <Link className="primary-link" to="/admin/videos/new">
          新建视频
        </Link>
      </div>

      {error ? <p className="form-error">{error}</p> : null}

      <div className="admin-post-list">
        {videos.map((video) => (
          <article key={video.id} className="admin-post-card">
            <div className="admin-post-card-top">
              <div>
                <p className="eyebrow">Video</p>
                <h3>{video.title}</h3>
              </div>
            </div>
            <p>{video.description}</p>
            <div className="admin-post-meta">
              <span>{formatDate(video.publishedAt)}</span>
              <span>{video.thumbnailUrl ? "有封面" : "无封面"}</span>
              <a href={video.url} target="_blank" rel="noreferrer">
                打开链接
              </a>
            </div>
            <div className="admin-post-actions">
              <Link className="ghost-link" to={`/admin/videos/${video.id}/edit`}>
                编辑
              </Link>
              <button type="button" className="danger-link" onClick={() => handleDelete(video.id)}>
                删除
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
