import type { FormEvent } from "react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { loginAdmin } from "../lib/api";

export function AdminLoginPage() {
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      setLoading(true);
      const response = await loginAdmin(username, password);
      if (response.session.authenticated) {
        navigate("/admin", { replace: true });
      }
    } catch (submitError) {
      setError(submitError instanceof Error ? submitError.message : "登录失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="admin-login-page">
      <div className="admin-login-card">
        <p className="eyebrow">Content Admin</p>
        <h1>登录博客后台</h1>
        <p>使用管理员账号进入内容管理后台，维护文章与展示结构。</p>

        <form className="admin-login-form" onSubmit={handleSubmit}>
          <label>
            用户名
            <input value={username} onChange={(event) => setUsername(event.target.value)} />
          </label>
          <label>
            密码
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
            />
          </label>
          {error ? <p className="form-error">{error}</p> : null}
          <button type="submit" className="primary-link" disabled={loading}>
            {loading ? "登录中..." : "进入后台"}
          </button>
        </form>
      </div>
    </section>
  );
}
