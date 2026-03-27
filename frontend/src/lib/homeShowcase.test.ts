import { describe, expect, it } from "vitest";
import { getHomeFeaturedShowcase } from "./homeShowcase";

describe("getHomeFeaturedShowcase", () => {
  it("uses featured posts as the homepage showcase", () => {
    const showcase = getHomeFeaturedShowcase([
      {
        slug: "featured-post",
        title: "精选文章",
        summary: "只在首页精选里展示",
        category: "Backend",
        readTime: "5 分钟",
        coverLabel: "精选",
        tags: [],
        featured: true,
        publishedAt: "2026-03-27",
        likeCount: 0,
        commentCount: 0,
      },
    ]);

    expect(showcase).toMatchObject({
      title: "首页精选",
      isEmpty: false,
    });
    expect(showcase.posts).toHaveLength(1);
    expect(showcase.posts[0].slug).toBe("featured-post");
  });

  it("returns an empty state when no featured posts are configured", () => {
    const showcase = getHomeFeaturedShowcase([]);

    expect(showcase).toMatchObject({
      title: "首页精选",
      isEmpty: true,
      emptyMessage: "当前还没有开启首页精选的文章，所以首页不会展示文章列表。",
    });
    expect(showcase.posts).toEqual([]);
  });
});
