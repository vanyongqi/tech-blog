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
            <p className="about-lead">
              目前在某网络空间测绘搜索引擎担任后端研发工程师，主要负责后端服务、数据链路和分析系统相关工作。
            </p>
            <p>
              我更关注的是系统怎么在真实流量下稳定运行，链路怎么在复杂数据条件下保持清晰，以及一套工程能力如何从“能跑”
              走到“可观测、可维护、可扩展”。
            </p>
          </div>

          <div className="about-block">
            <h3>工作经历</h3>
            <div className="about-meta-grid">
              <div className="about-meta-card">
                <strong>某网络空间测绘搜索引擎</strong>
                <span>后端研发工程师</span>
              </div>
              <div className="about-meta-card">
                <strong>时间区间</strong>
                <span>2024.06 - 至今</span>
              </div>
            </div>
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
