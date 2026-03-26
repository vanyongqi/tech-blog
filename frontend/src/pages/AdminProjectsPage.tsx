import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useAdminSessionContext } from "../components/AdminShell";
import { deleteAdminProject, getAdminProjects } from "../lib/api";
import type { AdminProjectPayload } from "../types";

export function AdminProjectsPage() {
  useAdminSessionContext();
  const [projects, setProjects] = useState<AdminProjectPayload[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const response = await getAdminProjects();
        if (!cancelled) {
          setProjects(response.projects);
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
    if (!window.confirm("确认删除这个项目吗？")) {
      return;
    }
    try {
      await deleteAdminProject(id);
      setProjects((currentProjects) => currentProjects.filter((project) => project.id !== id));
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "删除失败");
    }
  }

  if (loading) {
    return <section className="admin-panel">项目列表加载中...</section>;
  }

  return (
    <section className="admin-panel">
      <div className="admin-panel-heading">
        <div>
          <p className="eyebrow">Projects</p>
          <h2>项目管理</h2>
        </div>
        <Link className="primary-link" to="/admin/projects/new">
          新建项目
        </Link>
      </div>

      {error ? <p className="form-error">{error}</p> : null}

      <div className="admin-post-list">
        {projects.map((project) => (
          <article key={project.id} className="admin-post-card">
            <div className="admin-post-card-top">
              <div>
                <p className="eyebrow">{project.status}</p>
                <h3>{project.name}</h3>
              </div>
            </div>
            <p>{project.summary}</p>
            <div className="admin-post-meta">
              <span>{project.imageUrl ? "有图片" : "无图片"}</span>
              <span>{project.techStack.slice(0, 3).join(" / ")}</span>
              <a href={project.link} target="_blank" rel="noreferrer">
                打开链接
              </a>
            </div>
            <div className="admin-post-actions">
              <Link className="ghost-link" to={`/admin/projects/${project.id}/edit`}>
                编辑
              </Link>
              <button type="button" className="danger-link" onClick={() => handleDelete(project.id)}>
                删除
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
