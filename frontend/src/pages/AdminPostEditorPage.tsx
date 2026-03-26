import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useAdminSessionContext } from "../components/AdminShell";
import { createAdminPost, getAdminPost, updateAdminPost } from "../lib/api";
import type { AdminContentBlock, AdminSavePostRequest } from "../types";

type BlockDraft = {
  kind: string;
  title: string;
  text: string;
  url: string;
  itemsText: string;
};

type FormState = {
  slug: string;
  title: string;
  summary: string;
  category: string;
  readTime: string;
  heroNote: string;
  coverLabel: string;
  featured: boolean;
  publishedAt: string;
  blocks: BlockDraft[];
};

const emptyForm: FormState = {
  slug: "",
  title: "",
  summary: "",
  category: "Engineering",
  readTime: "5 分钟",
  heroNote: "",
  coverLabel: "新文章",
  featured: false,
  publishedAt: new Date().toISOString().slice(0, 10),
  blocks: [{ kind: "paragraph", title: "", text: "", url: "", itemsText: "" }],
};

export function AdminPostEditorPage() {
  useAdminSessionContext();
  const navigate = useNavigate();
  const { slug } = useParams();
  const isEditMode = Boolean(slug);

  const [form, setForm] = useState<FormState>(emptyForm);
  const [loading, setLoading] = useState(isEditMode);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

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
          heroNote: response.post.heroNote,
          coverLabel: response.post.coverLabel,
          featured: response.post.featured,
          publishedAt: response.post.publishedAt,
          blocks: response.post.blocks.map((block) => ({
            kind: block.kind,
            title: block.title ?? "",
            text: block.text ?? "",
            url: block.url ?? "",
            itemsText: (block.items ?? []).join("\n"),
          })),
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

  function updateBlock(index: number, nextBlock: Partial<BlockDraft>) {
    setForm((currentForm) => ({
      ...currentForm,
      blocks: currentForm.blocks.map((block, blockIndex) =>
        blockIndex === index ? { ...block, ...nextBlock } : block,
      ),
    }));
  }

  function addBlock() {
    setForm((currentForm) => ({
      ...currentForm,
      blocks: [...currentForm.blocks, { kind: "paragraph", title: "", text: "", url: "", itemsText: "" }],
    }));
  }

  function removeBlock(index: number) {
    setForm((currentForm) => ({
      ...currentForm,
      blocks: currentForm.blocks.filter((_, blockIndex) => blockIndex !== index),
    }));
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const payload: AdminSavePostRequest = {
      slug: form.slug.trim(),
      title: form.title.trim(),
      summary: form.summary.trim(),
      category: form.category.trim(),
      readTime: form.readTime.trim(),
      heroNote: form.heroNote.trim(),
      coverLabel: form.coverLabel.trim(),
      tags: [],
      featured: form.featured,
      publishedAt: form.publishedAt,
      blocks: form.blocks.map<AdminContentBlock>((block) => ({
        kind: block.kind,
        title: block.title.trim(),
        text: block.text.trim(),
        url: block.url.trim(),
        items: block.itemsText
          .split("\n")
          .map((item) => item.trim())
          .filter(Boolean),
      })),
    };

    try {
      setSaving(true);
      if (isEditMode && slug) {
        const response = await updateAdminPost(slug, payload);
        navigate(`/admin/posts/${response.post.slug}/edit`, { replace: true });
      } else {
        const response = await createAdminPost(payload);
        navigate(`/admin/posts/${response.post.slug}/edit`, { replace: true });
      }
      setError("");
    } catch (submitError) {
      setError(submitError instanceof Error ? submitError.message : "保存失败");
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
            Slug
            <input value={form.slug} onChange={(event) => setForm({ ...form, slug: event.target.value })} />
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

        <label>
          Hero Note
          <textarea value={form.heroNote} onChange={(event) => setForm({ ...form, heroNote: event.target.value })} />
        </label>

        <label className="checkbox-row">
          <input
            type="checkbox"
            checked={form.featured}
            onChange={(event) => setForm({ ...form, featured: event.target.checked })}
          />
          设为精选文章
        </label>

        <section className="admin-blocks-section">
          <div className="admin-panel-heading tight-heading">
            <div>
              <p className="eyebrow">Blocks</p>
              <h2>内容编排</h2>
            </div>
            <button type="button" className="ghost-link" onClick={addBlock}>
              添加区块
            </button>
          </div>

          <div className="admin-block-list">
            {form.blocks.map((block, index) => (
              <article key={`${block.kind}-${index}`} className="admin-block-card">
                <div className="admin-block-header">
                  <strong>区块 {index + 1}</strong>
                  <button type="button" className="danger-link" onClick={() => removeBlock(index)}>
                    删除
                  </button>
                </div>
                <label>
                  类型
                  <select value={block.kind} onChange={(event) => updateBlock(index, { kind: event.target.value })}>
                    <option value="paragraph">paragraph</option>
                    <option value="heading">heading</option>
                    <option value="quote">quote</option>
                    <option value="list">list</option>
                    <option value="video">video</option>
                  </select>
                </label>
                <label>
                  标题
                  <input value={block.title} onChange={(event) => updateBlock(index, { title: event.target.value })} />
                </label>
                {block.kind === "video" ? (
                  <label>
                    视频地址
                    <input
                      value={block.url}
                      onChange={(event) => updateBlock(index, { url: event.target.value })}
                      placeholder="支持 mp4/webm/ogg，或可嵌入的 YouTube / Bilibili 链接"
                    />
                  </label>
                ) : null}
                <label>
                  {block.kind === "video" ? "说明文字" : "文本"}
                  <textarea value={block.text} onChange={(event) => updateBlock(index, { text: event.target.value })} />
                </label>
                {block.kind === "list" ? (
                  <label>
                    列表项
                    <textarea
                      value={block.itemsText}
                      onChange={(event) => updateBlock(index, { itemsText: event.target.value })}
                      placeholder="每行一个条目"
                    />
                  </label>
                ) : null}
              </article>
            ))}
          </div>
        </section>

        {error ? <p className="form-error">{error}</p> : null}

        <div className="admin-editor-actions">
          <button type="submit" className="primary-link" disabled={saving}>
            {saving ? "保存中..." : isEditMode ? "保存修改" : "创建文章"}
          </button>
        </div>
      </form>
    </section>
  );
}
