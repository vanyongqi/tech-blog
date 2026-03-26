import { startTransition, useEffect, useState } from "react";
import { getHomeData } from "../lib/api";
import type { HomeResponse } from "../types";

export function ProjectsPage() {
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
        <p>项目内容加载中...</p>
      </section>
    );
  }

  return (
    <div className="page">
      <section className="content-section">
        <div className="section-heading">
          <div>
            <h2>项目</h2>
          </div>
        </div>
        <div className="resource-grid">
          {homeData.projects.map((project) => (
            <a
              key={project.name}
              className="resource-card"
              href={project.link}
              target="_blank"
              rel="noreferrer"
            >
              {project.imageUrl ? (
                <div className="resource-thumb-wrap">
                  <img className="resource-thumb" src={project.imageUrl} alt={project.name} loading="lazy" />
                </div>
              ) : (
                <div className="resource-thumb-placeholder">PROJECT</div>
              )}
              <span className="resource-kicker">{project.status}</span>
              <h3>{project.name}</h3>
              <p>{project.summary}</p>
            </a>
          ))}
        </div>
      </section>
    </div>
  );
}
