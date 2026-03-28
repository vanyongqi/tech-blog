import type { ChangeEvent, FormEvent } from "react";
import { useEffect, useRef, useState } from "react";
import { Link, useLocation, useNavigate, useParams } from "react-router-dom";
import { AdminSaveNoticeBanner } from "../components/AdminSaveNoticeBanner";
import { useAdminSessionContext } from "../components/AdminShell";
import { createAdminProject, getAdminProject, updateAdminProject, uploadAdminProjectAsset } from "../lib/api";
import { getAdminEditorSaveSuccessMessage, type AdminEditorNotice } from "../lib/adminEditorNotice";
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
  const location = useLocation();
  const { id } = useParams();
  const isEditMode = Boolean(id);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const [form, setForm] = useState<FormState>(emptyForm);
  const [loading, setLoading] = useState(isEditMode);
  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState("");
  const [saveNotice, setSaveNotice] = useState<AdminEditorNotice | null>(null);

  useEffect(() => {
    const nextNotice = (location.state as { saveNotice?: AdminEditorNotice } | null)?.saveNotice;
    if (!nextNotice) {
      return;
    }
    setSaveNotice(nextNotice);
    navigate(location.pathname, { replace: true, state: null });
  }, [location.pathname, location.state, navigate]);

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

  async function handleFileSelect(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    try {
      setUploading(true);
      const response = await uploadAdminProjectAsset(file);
      setForm((currentForm) => ({
        ...currentForm,
        imageUrl: response.url,
      }));
      setSaveNotice({
        type: "success",
        message: "项目图片预览已更新，保存修改后会同步到前台。",
      });
    } catch (uploadError) {
      setSaveNotice({
        type: "error",
        message: uploadError instanceof Error ? uploadError.message : "图片上传失败",
      });
    } finally {
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  }

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
      setSaveNotice(null);
      if (isEditMode && id) {
        await updateAdminProject(Number(id), payload);
        setSaveNotice({
          type: "success",
          message: getAdminEditorSaveSuccessMessage("project", true),
        });
      } else {
        const response = await createAdminProject(payload);
        navigate(`/admin/projects/${response.project.id}/edit`, {
          replace: true,
          state: { saveNotice: { type: "success", message: getAdminEditorSaveSuccessMessage("project", false) } },
        });
      }
      setError("");
    } catch (submitError) {
      setSaveNotice({
        type: "error",
        message: submitError instanceof Error ? submitError.message : "保存失败",
      });
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return <section className="admin-panel">项目编辑器加载中...</section>;
  }

  return (
    <section className="admin-panel">
      {saveNotice ? <AdminSaveNoticeBanner notice={saveNotice} onDismiss={() => setSaveNotice(null)} /> : null}

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
        <div className="admin-inline-actions">
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            hidden
            onChange={handleFileSelect}
          />
          <button type="button" className="ghost-link" onClick={() => fileInputRef.current?.click()} disabled={uploading}>
            {uploading ? "上传中..." : form.imageUrl ? "重新上传图片" : "上传图片"}
          </button>
        </div>
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
          <button type="submit" className="primary-link" disabled={saving || uploading}>
            {saving ? "保存中..." : isEditMode ? "保存修改" : "创建项目"}
          </button>
        </div>
      </form>
    </section>
  );
}
