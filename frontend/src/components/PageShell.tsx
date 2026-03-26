import type { ReactNode } from "react";
import { Link, NavLink, Outlet } from "react-router-dom";

export function PageShell({ children }: { children?: ReactNode }) {
  return (
    <div className="app-shell">
      <div className="ambient ambient-left" />
      <div className="ambient ambient-right" />
      <header className="site-header">
        <Link className="brand" to="/">
          <span className="brand-mark">F</span>
          <span>
            Fite
            <small>Backend Systems Notebook</small>
          </span>
        </Link>
        <nav className="site-nav">
          <NavLink to="/">首页</NavLink>
          <NavLink to="/articles">文章</NavLink>
          <NavLink to="/projects">项目</NavLink>
          <NavLink to="/videos">视频</NavLink>
          <NavLink to="/about">About</NavLink>
        </nav>
      </header>
      <main>{children ?? <Outlet />}</main>
    </div>
  );
}
