import type { ChangeEvent, ClipboardEvent, DragEvent, FormEvent } from "react";
import { useEffect, useRef, useState } from "react";
import { Link, useLocation, useNavigate, useParams } from "react-router-dom";
import { useAdminSessionContext } from "../components/AdminShell";
import { MarkdownContent } from "../components/MarkdownContent";
import { createAdminPost, getAdminPost, updateAdminPost, uploadAdminPostAsset } from "../lib/api";
import { getAdminEditorSaveSuccessMessage, type AdminEditorNotice } from "../lib/adminEditorNotice";
import { buildMarkdownImageSnippet } from "../lib/markdownImages";
import type { AdminSavePostRequest } from "../types";

type FormState = {
  slug: string;
  title: string;
  summary: string;
  category: string;
  readTime: string;
  coverLabel: string;
  featured: boolean;
  publishedAt: string;
  contentMarkdown: string;
};

const emptyForm: FormState = {
  slug: "",
  title: "",
  summary: "",
  category: "Engineering",
  readTime: "5 分钟",
  coverLabel: "新文章",
  featured: false,
  publishedAt: new Date().toISOString().slice(0, 10),
  contentMarkdown: "",
};

export function AdminPostEditorPage() {
  useAdminSessionContext();
  const navigate = useNavigate();
  const location = useLocation();
  const { slug } = useParams();
  const isEditMode = Boolean(slug);
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const markdownTextareaRef = useRef<HTMLTextAreaElement | null>(null);

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
    if (!isEditMode || !slug) {
      return;
    }
    const currentSlug = slug;
    let cancelled = false;

    async function load() {
      try {
        const response = await getAdminPost(currentSlug);
        if (cancelled) {
          return;
        }

        setForm({
          slug: response.post.slug,
          title: response.post.title,
          summary: response.post.summary,
          category: response.post.category,
          readTime: response.post.readTime,
          coverLabel: response.post.coverLabel,
          featured: response.post.featured,
          publishedAt: response.post.publishedAt,
          contentMarkdown: response.post.contentMarkdown,
        });
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
  }, [isEditMode, slug]);

  function insertMarkdownSnippet(snippet: string) {
    const textarea = markdownTextareaRef.current;
    const nextSnippet = snippet.trim();
    if (!textarea) {
      setForm((currentForm) => ({
        ...currentForm,
        contentMarkdown: currentForm.contentMarkdown.trimEnd()
          ? `${currentForm.contentMarkdown.trimEnd()}\n\n${nextSnippet}`
          : nextSnippet,
      }));
      return;
    }

    const { selectionStart, selectionEnd } = textarea;
    const currentValue = form.contentMarkdown;
    const prefix = currentValue.slice(0, selectionStart);
    const suffix = currentValue.slice(selectionEnd);
    const needsLeadingBreak = prefix.length > 0 && !prefix.endsWith("\n");
    const needsTrailingBreak = suffix.length > 0 && !suffix.startsWith("\n");
    const insertion = `${needsLeadingBreak ? "\n\n" : ""}${nextSnippet}${needsTrailingBreak ? "\n\n" : ""}`;
    const nextValue = `${prefix}${insertion}${suffix}`;

    setForm((currentForm) => ({
      ...currentForm,
      contentMarkdown: nextValue,
    }));

    const nextCursor = prefix.length + insertion.length;
    requestAnimationFrame(() => {
      const nextTextarea = markdownTextareaRef.current;
      if (!nextTextarea) {
        return;
      }
      nextTextarea.focus();
      nextTextarea.setSelectionRange(nextCursor, nextCursor);
    });
  }

  async function uploadFiles(files: File[]) {
    if (files.length === 0) {
      return;
    }

    try {
      setUploading(true);
      const snippets: string[] = [];

      for (const file of files) {
        const response = await uploadAdminPostAsset(file);
        snippets.push(buildMarkdownImageSnippet(file.name, response.url));
      }

      insertMarkdownSnippet(snippets.join("\n\n"));
      setSaveNotice({
        type: "success",
        message: "图片已插入正文预览，保存修改后会同步到前台。",
      });
      setError("");
    } catch (uploadError) {
      setError(uploadError instanceof Error ? uploadError.message : "图片上传失败");
    } finally {
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  }

  async function handleFileSelect(event: ChangeEvent<HTMLInputElement>) {
    const files = Array.from(event.target.files ?? []);
    await uploadFiles(files);
  }

  async function handleEditorPaste(event: ClipboardEvent<HTMLTextAreaElement>) {
    const imageFiles = Array.from(event.clipboardData.items)
      .filter((item) => item.type.startsWith("image/"))
      .map((item) => item.getAsFile())
      .filter((file): file is File => file !== null);

    if (imageFiles.length === 0) {
      return;
    }

    event.preventDefault();
    await uploadFiles(imageFiles);
  }

  async function handleEditorDrop(event: DragEvent<HTMLTextAreaElement>) {
    event.preventDefault();
    const files = Array.from(event.dataTransfer.files).filter((file) => file.type.startsWith("image/"));
    await uploadFiles(files);
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const payload: AdminSavePostRequest = {
      slug: form.slug.trim(),
      title: form.title.trim(),
      summary: form.summary.trim(),
      category: form.category.trim(),
      readTime: form.readTime.trim(),
      coverLabel: form.coverLabel.trim(),
      contentMarkdown: form.contentMarkdown.trim(),
      tags: [],
      featured: form.featured,
      publishedAt: form.publishedAt,
    };

    try {
      setSaving(true);
      setSaveNotice(null);
      if (isEditMode && slug) {
        await updateAdminPost(slug, payload);
        setSaveNotice({
          type: "success",
          message: getAdminEditorSaveSuccessMessage("post", true),
        });
      } else {
        const response = await createAdminPost(payload);
        navigate(`/admin/posts/${response.post.slug}/edit`, {
          replace: true,
          state: { saveNotice: { type: "success", message: getAdminEditorSaveSuccessMessage("post", false) } },
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
    return <section className="admin-panel">编辑器加载中...</section>;
  }

  return (
    <section className="admin-panel">
      <div className="admin-panel-heading">
        <div>
          <p className="eyebrow">Editor</p>
          <h2>{isEditMode ? "编辑文章" : "新建文章"}</h2>
        </div>
        <Link className="ghost-link" to="/admin">
          返回列表
        </Link>
      </div>

      <form className="admin-editor-form" onSubmit={handleSubmit}>
        <div className="admin-form-grid">
          <label>
            标题
            <input value={form.title} onChange={(event) => setForm({ ...form, title: event.target.value })} />
          </label>
          <label>
            分类
            <input value={form.category} onChange={(event) => setForm({ ...form, category: event.target.value })} />
          </label>
          <label>
            阅读时长
            <input value={form.readTime} onChange={(event) => setForm({ ...form, readTime: event.target.value })} />
          </label>
          <label>
            封面标签
            <input value={form.coverLabel} onChange={(event) => setForm({ ...form, coverLabel: event.target.value })} />
          </label>
          <label>
            发布日期
            <input type="date" value={form.publishedAt} onChange={(event) => setForm({ ...form, publishedAt: event.target.value })} />
          </label>
        </div>

        <label>
          摘要
          <textarea value={form.summary} onChange={(event) => setForm({ ...form, summary: event.target.value })} />
        </label>

        <section className={form.featured ? "admin-feature-card active" : "admin-feature-card"}>
          <div className="admin-feature-copy">
            <p className="eyebrow">Homepage</p>
            <h3>首页精选</h3>
          </div>
          <label className="admin-switch">
            <input
              type="checkbox"
              checked={form.featured}
              onChange={(event) => setForm({ ...form, featured: event.target.checked })}
            />
            <span className="admin-switch-track" aria-hidden="true">
              <span className="admin-switch-thumb" />
            </span>
            <span className="admin-switch-label">{form.featured ? "已开启" : "未开启"}</span>
          </label>
        </section>

        <section className="admin-markdown-editor">
          <div className="admin-panel-heading tight-heading">
            <div>
              <p className="eyebrow">Markdown</p>
              <h2>正文编辑</h2>
            </div>
            <div className="admin-inline-actions">
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                multiple
                hidden
                onChange={handleFileSelect}
              />
              <button type="button" className="ghost-link" onClick={() => fileInputRef.current?.click()} disabled={uploading}>
                {uploading ? "上传中..." : "上传图片"}
              </button>
            </div>
          </div>

          <p className="admin-field-hint">
            支持直接写 Markdown，也支持按钮上传、粘贴截图、拖拽图片到编辑区，上传成功后会自动插入统一图片语法。
            图片格式统一为：<code>![说明](/api/assets/资源ID "size=lg align=center")</code>。
            可用尺寸：<code>sm / md / lg / full</code>；可用对齐：<code>left / center / right</code>；也支持 <code>width=320</code> 或 <code>width=60%</code>。
          </p>

          <textarea
            ref={markdownTextareaRef}
            className="admin-markdown-textarea"
            value={form.contentMarkdown}
            onChange={(event) => setForm({ ...form, contentMarkdown: event.target.value })}
            onPaste={handleEditorPaste}
            onDrop={handleEditorDrop}
            onDragOver={(event) => event.preventDefault()}
            placeholder={"# 文章标题\n\n从这里开始写正文...\n\n## 小节标题\n\n支持 Markdown 列表、引用、代码块和图片。"}
          />

          <div className="admin-markdown-preview-card">
            <div className="admin-panel-heading tight-heading">
              <div>
                <p className="eyebrow">Preview</p>
                <h2>正文预览</h2>
              </div>
            </div>
            <div className="post-body post-body-preview">
              {form.contentMarkdown.trim() ? (
                <MarkdownContent content={form.contentMarkdown} />
              ) : (
                <p className="admin-field-hint">正文为空，输入 Markdown 后这里会实时预览。</p>
              )}
            </div>
          </div>
        </section>

        {error ? <p className="form-error">{error}</p> : null}
        {saveNotice ? <p className={saveNotice.type === "success" ? "form-success" : "form-error"}>{saveNotice.message}</p> : null}

        <div className="admin-editor-actions">
          <button type="submit" className="primary-link" disabled={saving || uploading}>
            {saving ? "保存中..." : isEditMode ? "保存修改" : "创建文章"}
          </button>
        </div>
      </form>
    </section>
  );
}
