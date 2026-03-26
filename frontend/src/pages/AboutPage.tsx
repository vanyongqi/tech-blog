import { startTransition, useEffect, useState } from "react";
import { getHomeData } from "../lib/api";
import type { HomeResponse } from "../types";

export function AboutPage() {
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
        <p>About 内容加载中...</p>
      </section>
    );
  }

  return (
    <div className="page">
      <section className="content-section">
        <div className="section-heading">
          <div>
            <h2>About Me</h2>
          </div>
        </div>

        <div className="about-layout">
          <div className="about-block">
            <p className="about-lead">{homeData.site.intro}</p>
            <p>{homeData.site.motto}</p>
          </div>

          <div className="about-block">
            <h3>联系方式</h3>
            <div className="about-links">
              {homeData.site.socialLinks.map((link) => (
                <a key={link.label} href={link.url} target="_blank" rel="noreferrer">
                  {link.label}
                </a>
              ))}
            </div>
            <p className="about-meta">邮箱：{homeData.site.email}</p>
          </div>

          <div className="about-block">
            <h3>关注方向</h3>
            <div className="about-tech">
              {homeData.site.techStack.map((item) => (
                <span key={item}>{item}</span>
              ))}
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
