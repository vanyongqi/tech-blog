import { startTransition, useEffect, useState } from "react";
import { getAllPosts } from "../lib/api";
import { PostCard } from "../components/PostCard";
import type { PostSummaryPayload } from "../types";

function sortPostsByPublishedAt(posts: PostSummaryPayload[]) {
  return [...posts].sort((left, right) => {
    return new Date(right.publishedAt).getTime() - new Date(left.publishedAt).getTime();
  });
}

export function ArticlesPage() {
  const [posts, setPosts] = useState<PostSummaryPayload[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const response = await getAllPosts();
        if (cancelled) {
          return;
        }
        startTransition(() => {
          setPosts(sortPostsByPublishedAt(response.posts));
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
  }, []);

  if (error) {
    return (
      <section className="state-panel">
        <p>数据加载失败：{error}</p>
      </section>
    );
  }

  return (
    <div className="page">
      <section className="content-section">
        <div className="section-heading">
          <div>
            <h2>文章</h2>
          </div>
        </div>
        <div className="archive-grid archive-grid-simple">
          {posts.map((post) => (
            <PostCard key={post.slug} post={post} />
          ))}
        </div>
      </section>
    </div>
  );
}
