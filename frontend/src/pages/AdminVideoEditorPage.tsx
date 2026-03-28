import type { ChangeEvent, FormEvent } from "react";
import { useEffect, useRef, useState } from "react";
import { Link, useLocation, useNavigate, useParams } from "react-router-dom";
import { AdminSaveNoticeBanner } from "../components/AdminSaveNoticeBanner";
import { useAdminSessionContext } from "../components/AdminShell";
import {
  createAdminVideo,
  getAdminVideo,
  suggestAdminVideoThumbnail,
  uploadAdminVideoAsset,
  updateAdminVideo,
} from "../lib/api";
import { getAdminEditorSaveSuccessMessage, type AdminEditorNotice } from "../lib/adminEditorNotice";
import type { AdminSaveVideoRequest } from "../types";

type FormState = {
  title: string;
  description: string;
  url: string;
  thumbnailUrl: string;
  publishedAt: string;
};

const emptyForm: FormState = {
  title: "",
  description: "",
  url: "",
  thumbnailUrl: "",
  publishedAt: new Date().toISOString().slice(0, 10),
};

function deriveYouTubeThumbnail(videoURL: string) {
  try {
    const parsed = new URL(videoURL.trim());
    if (parsed.hostname.includes("youtube.com")) {
      const videoId = parsed.searchParams.get("v");
      if (videoId) {
        return `https://i.ytimg.com/vi/${videoId}/hqdefault.jpg`;
      }
    }
    if (parsed.hostname.includes("youtu.be")) {
      const videoId = parsed.pathname.replace("/", "");
      if (videoId) {
        return `https://i.ytimg.com/vi/${videoId}/hqdefault.jpg`;
      }
    }
  } catch {
    return "";
  }
  return "";
}

function deriveThumbnailFromURL(videoURL: string) {
  return deriveYouTubeThumbnail(videoURL);
}

export function AdminVideoEditorPage() {
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
  const [generatingThumbnail, setGeneratingThumbnail] = useState(false);
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
        const response = await getAdminVideo(Number(id));
        if (!cancelled) {
          setForm({
            title: response.video.title,
            description: response.video.description,
            url: response.video.url,
            thumbnailUrl: response.video.thumbnailUrl,
            publishedAt: response.video.publishedAt,
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
      const response = await uploadAdminVideoAsset(file);
      setForm((currentForm) => ({
        ...currentForm,
        thumbnailUrl: response.url,
      }));
      setSaveNotice({
        type: "success",
        message: "封面预览已更新，保存修改后会同步到前台。",
      });
    } catch (uploadError) {
      setSaveNotice({
        type: "error",
        message: uploadError instanceof Error ? uploadError.message : "封面上传失败",
      });
    } finally {
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  }

  useEffect(() => {
    if (form.thumbnailUrl.trim() !== "") {
      return;
    }
    const derivedThumbnail = deriveYouTubeThumbnail(form.url);
    if (derivedThumbnail) {
      setForm((currentForm) => {
        if (currentForm.thumbnailUrl.trim() !== "" || currentForm.url !== form.url) {
          return currentForm;
        }
        return {
          ...currentForm,
          thumbnailUrl: derivedThumbnail,
        };
      });
    }
  }, [form.url, form.thumbnailUrl]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const payload: AdminSaveVideoRequest = {
      title: form.title.trim(),
      description: form.description.trim(),
      url: form.url.trim(),
      thumbnailUrl: form.thumbnailUrl.trim(),
      publishedAt: form.publishedAt,
    };

    try {
      setSaving(true);
      setSaveNotice(null);
      if (isEditMode && id) {
        await updateAdminVideo(Number(id), payload);
        setSaveNotice({
          type: "success",
          message: getAdminEditorSaveSuccessMessage("video", true),
        });
      } else {
        const response = await createAdminVideo(payload);
        navigate(`/admin/videos/${response.video.id}/edit`, {
          replace: true,
          state: { saveNotice: { type: "success", message: getAdminEditorSaveSuccessMessage("video", false) } },
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

  async function handleAutoGenerateThumbnail() {
    const localThumbnail = deriveThumbnailFromURL(form.url);
    if (localThumbnail) {
      setForm((currentForm) => ({
        ...currentForm,
        thumbnailUrl: localThumbnail,
      }));
      setSaveNotice({
        type: "success",
        message: "封面预览已更新，保存修改后会同步到前台。",
      });
      setError("");
      return;
    }

    try {
      setGeneratingThumbnail(true);
      const response = await suggestAdminVideoThumbnail(form.url.trim());
      if (!response.thumbnailUrl) {
        setError("当前链接未能自动解析出封面，建议手填封面图链接。");
        return;
      }
      setForm((currentForm) => ({
        ...currentForm,
        thumbnailUrl: response.thumbnailUrl,
      }));
      setSaveNotice({
        type: "success",
        message: "封面预览已更新，保存修改后会同步到前台。",
      });
      setError("");
    } catch (generateError) {
      setError(generateError instanceof Error ? generateError.message : "自动生成封面失败");
    } finally {
      setGeneratingThumbnail(false);
    }
  }

  if (loading) {
    return <section className="admin-panel">视频编辑器加载中...</section>;
  }

  return (
    <section className="admin-panel">
      {saveNotice ? <AdminSaveNoticeBanner notice={saveNotice} onDismiss={() => setSaveNotice(null)} /> : null}

      <div className="admin-panel-heading">
        <div>
          <p className="eyebrow">Video Editor</p>
          <h2>{isEditMode ? "编辑视频" : "新建视频"}</h2>
        </div>
        <Link className="ghost-link" to="/admin/videos">
          返回列表
        </Link>
      </div>

      <form className="admin-editor-form" onSubmit={handleSubmit}>
        <label>
          标题
          <input value={form.title} onChange={(event) => setForm({ ...form, title: event.target.value })} />
        </label>
        <label>
          视频链接
          <input value={form.url} onChange={(event) => setForm({ ...form, url: event.target.value })} />
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
            {uploading ? "上传中..." : form.thumbnailUrl ? "重新上传封面" : "上传封面"}
          </button>
          <button type="button" className="ghost-link" onClick={handleAutoGenerateThumbnail} disabled={generatingThumbnail}>
            {generatingThumbnail ? "生成中..." : "自动生成封面"}
          </button>
        </div>
        <p className="admin-field-hint">上传或自动生成封面后，还需要点击“保存修改”才会真正生效。</p>
        {form.thumbnailUrl ? (
          <div className="admin-video-preview">
            <img src={form.thumbnailUrl} alt="视频封面预览" />
          </div>
        ) : null}
        <label>
          发布时间
          <input type="date" value={form.publishedAt} onChange={(event) => setForm({ ...form, publishedAt: event.target.value })} />
        </label>
        <label>
          简介
          <textarea value={form.description} onChange={(event) => setForm({ ...form, description: event.target.value })} />
        </label>

        {error ? <p className="form-error">{error}</p> : null}

        <div className="admin-editor-actions">
          <button type="submit" className="primary-link" disabled={saving || uploading}>
            {saving ? "保存中..." : isEditMode ? "保存修改" : "创建视频"}
          </button>
        </div>
      </form>
    </section>
  );
}
