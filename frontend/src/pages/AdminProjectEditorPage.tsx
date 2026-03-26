import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useAdminSessionContext } from "../components/AdminShell";
import { createAdminProject, getAdminProject, updateAdminProject } from "../lib/api";
import type { AdminSaveProjectRequest } from "../types";

type FormState = {
  name: string;
  summary: string;
  status: string;
  link: string;
  imageUrl: string;
  accent: string;
  techStackText: string;
};

const emptyForm: FormState = {
  name: "",
  summary: "",
  status: "Ongoing",
  link: "",
  imageUrl: "",
  accent: "ink",
  techStackText: "",
};

export function AdminProjectEditorPage() {
  useAdminSessionContext();
  const navigate = useNavigate();
  const { id } = useParams();
  const isEditMode = Boolean(id);

  const [form, setForm] = useState<FormState>(emptyForm);
  const [loading, setLoading] = useState(isEditMode);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!isEditMode || !id) {
      return;
    }

    let cancelled = false;
    async function load() {
      try {
        const response = await getAdminProject(Number(id));
        if (!cancelled) {
          setForm({
            name: response.project.name,
            summary: response.project.summary,
            status: response.project.status,
            link: response.project.link,
            imageUrl: response.project.imageUrl,
            accent: response.project.accent,
            techStackText: response.project.techStack.join(", "),
          });
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
  }, [id, isEditMode]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const payload: AdminSaveProjectRequest = {
      name: form.name.trim(),
      summary: form.summary.trim(),
      status: form.status.trim(),
      link: form.link.trim(),
      imageUrl: form.imageUrl.trim(),
      accent: form.accent.trim(),
      techStack: form.techStackText.split(",").map((item) => item.trim()).filter(Boolean),
    };

    try {
      setSaving(true);
      if (isEditMode && id) {
        const response = await updateAdminProject(Number(id), payload);
        navigate(`/admin/projects/${response.project.id}/edit`, { replace: true });
      } else {
        const response = await createAdminProject(payload);
        navigate(`/admin/projects/${response.project.id}/edit`, { replace: true });
      }
      setError("");
    } catch (submitError) {
      setError(submitError instanceof Error ? submitError.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return <section className="admin-panel">项目编辑器加载中...</section>;
  }

  return (
    <section className="admin-panel">
      <div className="admin-panel-heading">
        <div>
          <p className="eyebrow">Project Editor</p>
          <h2>{isEditMode ? "编辑项目" : "新建项目"}</h2>
        </div>
        <Link className="ghost-link" to="/admin/projects">
          返回列表
        </Link>
      </div>

      <form className="admin-editor-form" onSubmit={handleSubmit}>
        <label>
          名称
          <input value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} />
        </label>
        <label>
          链接
          <input value={form.link} onChange={(event) => setForm({ ...form, link: event.target.value })} />
        </label>
        <label>
          图片链接
          <input value={form.imageUrl} onChange={(event) => setForm({ ...form, imageUrl: event.target.value })} />
        </label>
        {form.imageUrl ? (
          <div className="admin-video-preview">
            <img src={form.imageUrl} alt="项目图片预览" />
          </div>
        ) : null}
        <label>
          状态
          <input value={form.status} onChange={(event) => setForm({ ...form, status: event.target.value })} />
        </label>
        <label>
          Accent
          <input value={form.accent} onChange={(event) => setForm({ ...form, accent: event.target.value })} />
        </label>
        <label>
          技术栈
          <input value={form.techStackText} onChange={(event) => setForm({ ...form, techStackText: event.target.value })} />
        </label>
        <label>
          简介
          <textarea value={form.summary} onChange={(event) => setForm({ ...form, summary: event.target.value })} />
        </label>

        {error ? <p className="form-error">{error}</p> : null}

        <div className="admin-editor-actions">
          <button type="submit" className="primary-link" disabled={saving}>
            {saving ? "保存中..." : isEditMode ? "保存修改" : "创建项目"}
          </button>
        </div>
      </form>
    </section>
  );
}
