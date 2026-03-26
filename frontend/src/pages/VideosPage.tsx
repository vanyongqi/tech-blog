import { startTransition, useEffect, useState } from "react";
import { getHomeData } from "../lib/api";
import type { HomeResponse } from "../types";

function deriveVideoThumbnail(videoURL: string, thumbnailURL: string) {
  if (thumbnailURL.trim() !== "") {
    return thumbnailURL.trim();
  }

  try {
    const parsed = new URL(videoURL);
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

export function VideosPage() {
  const [homeData, setHomeData] = useState<HomeResponse | null>(null);
  const [error, setError] = useState("");
  const [brokenImages, setBrokenImages] = useState<Record<number, boolean>>({});

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
        <p>视频内容加载中...</p>
      </section>
    );
  }

  return (
    <div className="page">
      <section className="content-section">
        <div className="section-heading">
          <div>
            <h2>视频</h2>
          </div>
        </div>
        <div className="resource-grid">
          {homeData.videos.map((video) => (
            <a
              key={video.id}
              className="resource-card"
              href={video.url}
              target="_blank"
              rel="noreferrer"
            >
              {!brokenImages[video.id] && deriveVideoThumbnail(video.url, video.thumbnailUrl) ? (
                <div className="resource-thumb-wrap">
                  <span className="resource-badge">▶</span>
                  <img
                    className="resource-thumb"
                    src={deriveVideoThumbnail(video.url, video.thumbnailUrl)}
                    alt={video.title}
                    loading="lazy"
                    onError={() =>
                      setBrokenImages((currentImages) => ({
                        ...currentImages,
                        [video.id]: true,
                      }))
                    }
                  />
                </div>
              ) : (
                <div className="resource-thumb-placeholder">
                  <span className="resource-badge">▶</span>
                  VIDEO
                </div>
              )}
              <span className="resource-kicker">{video.publishedAt}</span>
              <h3>{video.title}</h3>
              <p>{video.description}</p>
            </a>
          ))}
        </div>
      </section>
    </div>
  );
}
