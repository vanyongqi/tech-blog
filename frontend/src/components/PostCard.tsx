import { Link } from "react-router-dom";
import { formatDate } from "../lib/api";
import type { PostSummaryPayload } from "../types";

export function PostCard({
  post,
}: {
  post: PostSummaryPayload;
}) {
  return (
    <Link className="post-card post-card-compact post-card-link" to={`/posts/${post.slug}`}>
      <h3>
        <span className="post-title-link">{post.title}</span>
        <span className="post-summary-inline">：{post.summary}</span>
      </h3>
      <div className="post-card-stats post-card-stats-inline">
        <span>{formatDate(post.publishedAt)}</span>
        <span>{post.readTime}</span>
        <span>点赞 {post.likeCount}</span>
        <span>评论 {post.commentCount}</span>
        <span className="inline-link-text">阅读全文</span>
      </div>
    </Link>
  );
}
