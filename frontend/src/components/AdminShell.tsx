import { Outlet, useNavigate, useOutletContext } from "react-router-dom";
import { useEffect, useState } from "react";
import { getAdminSession, logoutAdmin } from "../lib/api";
import type { AdminSessionPayload } from "../types";

export function AdminShell() {
  const [session, setSession] = useState<AdminSessionPayload | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    let cancelled = false;

    async function loadSession() {
      try {
        const response = await getAdminSession();
        if (!cancelled) {
          setSession(response.session);
        }
      } catch {
        if (!cancelled) {
          navigate("/admin/login", { replace: true });
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadSession();
    return () => {
      cancelled = true;
    };
  }, [navigate]);

  async function handleLogout() {
    await logoutAdmin();
    navigate("/admin/login", { replace: true });
  }

  if (loading) {
    return <section className="state-panel">后台加载中...</section>;
  }

  if (!session?.authenticated) {
    return null;
  }

  return (
    <div className="admin-shell">
      <aside className="admin-sidebar">
        <div>
          <p className="eyebrow">Admin</p>
          <h1>内容后台</h1>
          <p>维护文章、组织内容结构、管理你的长期输出。</p>
        </div>
        <nav className="admin-nav">
          <button type="button" onClick={() => navigate("/admin")}>
            文章列表
          </button>
          <button type="button" onClick={() => navigate("/admin/posts/new")}>
            新建文章
          </button>
          <button type="button" onClick={() => navigate("/admin/projects")}>
            项目列表
          </button>
          <button type="button" onClick={() => navigate("/admin/projects/new")}>
            新建项目
          </button>
          <button type="button" onClick={() => navigate("/admin/videos")}>
            视频列表
          </button>
          <button type="button" onClick={() => navigate("/admin/videos/new")}>
            新建视频
          </button>
          <button type="button" onClick={() => navigate("/")}>
            返回前台
          </button>
        </nav>
      </aside>

      <div className="admin-main">
        <header className="admin-topbar">
          <div>
            <span className="signal-label">已登录</span>
            <strong>{session.username}</strong>
          </div>
          <button type="button" className="ghost-link" onClick={handleLogout}>
            退出登录
          </button>
        </header>
        <Outlet context={session} />
      </div>
    </div>
  );
}

export function useAdminSessionContext() {
  return useOutletContext<AdminSessionPayload>();
}
