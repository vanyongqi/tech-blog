import type { FormEvent } from "react";
import { startTransition, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import {
  createComment,
  formatDate,
  formatDateTime,
  getPost,
  toggleLike,
} from "../lib/api";
import type { PostDetailPayload } from "../types";

function isDirectVideo(url: string) {
  return /\.(mp4|webm|ogg)(\?.*)?$/i.test(url);
}

function toEmbeddedVideoURL(url: string) {
  const trimmed = url.trim();
  if (trimmed === "") {
    return "";
  }

  if (trimmed.includes("youtube.com/watch")) {
    try {
      const parsed = new URL(trimmed);
      const videoId = parsed.searchParams.get("v");
      if (videoId) {
        return `https://www.youtube.com/embed/${videoId}`;
      }
    } catch {
      return "";
    }
  }

  if (trimmed.includes("youtu.be/")) {
    try {
      const parsed = new URL(trimmed);
      const videoId = parsed.pathname.replace("/", "");
      if (videoId) {
        return `https://www.youtube.com/embed/${videoId}`;
      }
    } catch {
      return "";
    }
  }

  if (trimmed.includes("bilibili.com/video/")) {
    try {
      const parsed = new URL(trimmed);
      const parts = parsed.pathname.split("/").filter(Boolean);
      const bvid = parts[parts.length - 1];
      if (bvid) {
        return `https://player.bilibili.com/player.html?bvid=${bvid}&page=1`;
      }
    } catch {
      return "";
    }
  }

  if (trimmed.includes("/embed/") || trimmed.includes("player.bilibili.com")) {
    return trimmed;
  }

  return "";
}

export function PostDetailPage() {
  const { slug = "" } = useParams();
  const [post, setPost] = useState<PostDetailPayload | null>(null);
  const [visitorName, setVisitorName] = useState("");
  const [error, setError] = useState("");
  const [draft, setDraft] = useState("");
  const [actionError, setActionError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLiking, setIsLiking] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const response = await getPost(slug);
        if (cancelled) {
          return;
        }
        startTransition(() => {
          setPost(response.post);
          setVisitorName(response.visitor.displayName);
          setError("");
        });
      } catch (loadError) {
        if (cancelled) {
          return;
        }
        setError(loadError instanceof Error ? loadError.message : "加载失败");
      }
    }

    void load();
    return () => {
      cancelled = true;
    };
  }, [slug]);

  async function handleCommentSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const content = draft.trim();
    if (!content) {
      setActionError("请输入评论内容");
      return;
    }

    try {
      setIsSubmitting(true);
      const response = await createComment(slug, content);
      startTransition(() => {
        setPost((currentPost) =>
          currentPost
            ? {
                ...currentPost,
                commentCount: response.commentCount,
                comments: [response.comment, ...currentPost.comments],
              }
            : currentPost,
        );
        setDraft("");
        setVisitorName(response.visitor.displayName);
        setActionError("");
      });
    } catch (submitError) {
      setActionError(submitError instanceof Error ? submitError.message : "评论提交失败");
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleToggleLike() {
    if (!post || isLiking) {
      return;
    }

    try {
      setIsLiking(true);
      const response = await toggleLike(slug);
      startTransition(() => {
        setPost((currentPost) =>
          currentPost
            ? {
                ...currentPost,
                likeCount: response.likeCount,
                likedByVisitor: response.liked,
              }
            : currentPost,
        );
        setVisitorName(response.visitor.displayName);
        setActionError("");
      });
    } catch (likeError) {
      setActionError(likeError instanceof Error ? likeError.message : "点赞失败");
    } finally {
      setIsLiking(false);
    }
  }

  if (error) {
    return (
      <section className="state-panel">
        <p>文章加载失败：{error}</p>
        <Link to="/" className="ghost-link">
          返回首页
        </Link>
      </section>
    );
  }

  if (!post) {
    return (
      <section className="state-panel">
        <p>文章加载中...</p>
      </section>
    );
  }

  return (
    <article className="page post-page">
      <div className="post-hero">
        <Link to="/" className="back-link">
          返回首页
        </Link>
        <p className="eyebrow">{post.coverLabel}</p>
        <h1>{post.title}</h1>
        <p className="post-summary">{post.summary}</p>
        <div className="post-meta-bar post-meta-bar-inline">
          <span>{post.category}</span>
          <span>{formatDate(post.publishedAt)}</span>
          <span>{post.readTime}</span>
          <button
            type="button"
            className={post.likedByVisitor ? "action-button action-button-inline active" : "action-button action-button-inline"}
            onClick={handleToggleLike}
            disabled={isLiking}
          >
            {isLiking ? "处理中..." : post.likedByVisitor ? "已点赞" : "点赞"}
            <strong>{post.likeCount}</strong>
          </button>
          <span>评论 {post.commentCount}</span>
          <span>匿名身份 {visitorName}</span>
        </div>
        <blockquote>{post.heroNote}</blockquote>
      </div>

      <div className="post-body">
        {post.blocks.map((block, index) => {
          if (block.kind === "heading") {
            return <h2 key={`${block.kind}-${index}`}>{block.title}</h2>;
          }

          if (block.kind === "quote") {
            return <blockquote key={`${block.kind}-${index}`}>{block.text}</blockquote>;
          }

          if (block.kind === "list") {
            return (
              <ul key={`${block.kind}-${index}`}>
                {block.items?.map((item) => (
                  <li key={item}>{item}</li>
                ))}
              </ul>
            );
          }

          if (block.kind === "video") {
            const videoURL = block.url ?? "";
            const embeddedURL = toEmbeddedVideoURL(videoURL);

            return (
              <figure key={`${block.kind}-${index}`} className="video-block">
                {isDirectVideo(videoURL) ? (
                  <video controls preload="metadata" className="video-player">
                    <source src={videoURL} />
                  </video>
                ) : embeddedURL ? (
                  <div className="video-frame-wrap">
                    <iframe
                      className="video-frame"
                      src={embeddedURL}
                      title={block.title || `video-${index + 1}`}
                      loading="lazy"
                      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                      allowFullScreen
                    />
                  </div>
                ) : (
                  <a className="inline-link" href={videoURL} target="_blank" rel="noreferrer">
                    打开视频链接
                  </a>
                )}
                {block.title ? <figcaption className="video-caption">{block.title}</figcaption> : null}
                {block.text ? <p className="video-description">{block.text}</p> : null}
              </figure>
            );
          }

          return <p key={`${block.kind}-${index}`}>{block.text}</p>;
        })}
      </div>

      <section className="comment-section">
        <div className="section-heading compact-heading">
          <div>
            <p className="eyebrow">Comments</p>
            <h2>游客评论</h2>
          </div>
          <p>将以 {visitorName} 的固定匿名名发布，身份由 IP 与 User-Agent 生成。</p>
        </div>

        <form className="comment-form" onSubmit={handleCommentSubmit}>
          <textarea
            value={draft}
            onChange={(event) => setDraft(event.target.value)}
            placeholder="写点你的看法。支持游客评论，不需要注册。"
            maxLength={500}
          />
          <div className="comment-form-footer">
            <span>{draft.length}/500</span>
            <button type="submit" className="primary-link" disabled={isSubmitting}>
              {isSubmitting ? "发布中..." : "发布评论"}
            </button>
          </div>
        </form>

        {actionError ? <p className="form-error">{actionError}</p> : null}

        <div className="comment-list">
          {post.comments.length > 0 ? (
            post.comments.map((comment) => (
              <article key={comment.id} className="comment-card">
                <div className="comment-header">
                  <strong>{comment.authorName}</strong>
                  <span>{formatDateTime(comment.createdAt)}</span>
                </div>
                <p>{comment.content}</p>
              </article>
            ))
          ) : (
            <div className="comment-empty">还没有评论，来做第一个留下痕迹的人。</div>
          )}
        </div>
      </section>
    </article>
  );
}
