import { startTransition, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getHomeData } from "../lib/api";
import { getHomeFeaturedShowcase } from "../lib/homeShowcase";
import type { HomeResponse } from "../types";

export function HomePage() {
  const [homeData, setHomeData] = useState<HomeResponse | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const home = await getHomeData();
        if (cancelled) {
          return;
        }
        startTransition(() => {
          setHomeData(home);
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

  if (!homeData) {
    return (
      <section className="state-panel">
        <p>博客内容加载中...</p>
      </section>
    );
  }

  const showcase = getHomeFeaturedShowcase(homeData.featuredPosts);

  return (
    <div className="page home-page">
      <section className="hero-panel">
        <p className="eyebrow">718614413.xyz</p>
        <p className="hero-intro">{homeData.site.intro}</p>
      </section>

      <section className="content-section">
        <div className="section-heading">
          <div>
            <h2>{showcase.title}</h2>
          </div>
          <Link className="inline-link-text" to="/articles">
            查看全部文章
          </Link>
        </div>
        <div className="archive-grid archive-grid-simple">
          {showcase.isEmpty ? (
            <div className="content-empty-card">
              <p>{showcase.emptyMessage}</p>
            </div>
          ) : (
            showcase.posts.map((post) => (
              <Link key={post.slug} className="post-card post-card-compact post-card-link" to={`/posts/${post.slug}`}>
                <h3>
                  <span className="post-title-link">{post.title}</span>
                  <span className="post-summary-inline">：{post.summary}</span>
                </h3>
              </Link>
            ))
          )}
        </div>
      </section>
    </div>
  );
}
