import { Outlet, useLocation, useNavigate, useOutletContext } from "react-router-dom";
import { useEffect, useState } from "react";
import { getAdminSession, logoutAdmin } from "../lib/api";
import type { AdminSessionPayload } from "../types";

export function AdminShell() {
  const [session, setSession] = useState<AdminSessionPayload | null>(null);
  const [loading, setLoading] = useState(true);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const navItems = [
    { label: "文章列表", path: "/admin" },
    { label: "新建文章", path: "/admin/posts/new" },
    { label: "项目列表", path: "/admin/projects" },
    { label: "新建项目", path: "/admin/projects/new" },
    { label: "视频列表", path: "/admin/videos" },
    { label: "新建视频", path: "/admin/videos/new" },
    { label: "返回前台", path: "/" },
  ];

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

  useEffect(() => {
    setIsDrawerOpen(false);
  }, [location.pathname]);

  async function handleLogout() {
    await logoutAdmin();
    navigate("/admin/login", { replace: true });
  }

  function handleNavigate(path: string) {
    setIsDrawerOpen(false);
    navigate(path);
  }

  if (loading) {
    return <section className="state-panel">后台加载中...</section>;
  }

  if (!session?.authenticated) {
    return null;
  }

  return (
    <div className="admin-shell">
      <div
        className={isDrawerOpen ? "admin-sidebar-backdrop active" : "admin-sidebar-backdrop"}
        onClick={() => setIsDrawerOpen(false)}
        aria-hidden={!isDrawerOpen}
      />

      <aside className={isDrawerOpen ? "admin-sidebar active" : "admin-sidebar"}>
        <div>
          <div className="admin-sidebar-header">
            <div>
              <p className="eyebrow">Admin</p>
              <h1>内容后台</h1>
            </div>
            <button type="button" className="admin-sidebar-close" onClick={() => setIsDrawerOpen(false)}>
              关闭
            </button>
          </div>
          <p>维护文章与首页精选。</p>
        </div>
        <nav className="admin-nav">
          {navItems.map((item) => (
            <button key={item.path} type="button" onClick={() => handleNavigate(item.path)}>
              {item.label}
            </button>
          ))}
        </nav>
      </aside>

      <div className="admin-main">
        <header className="admin-topbar">
          <div className="admin-topbar-leading">
            <button type="button" className="admin-drawer-toggle" onClick={() => setIsDrawerOpen(true)}>
              菜单
            </button>
            <div>
              <span className="signal-label">已登录</span>
              <strong>{session.username}</strong>
            </div>
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
