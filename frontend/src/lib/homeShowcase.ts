import type { PostSummaryPayload } from "../types";

export type HomeFeaturedShowcase = {
  title: string;
  posts: PostSummaryPayload[];
  isEmpty: boolean;
  emptyMessage?: string;
};

export function getHomeFeaturedShowcase(featuredPosts: PostSummaryPayload[]): HomeFeaturedShowcase {
  if (featuredPosts.length === 0) {
    return {
      title: "首页精选",
      posts: [],
      isEmpty: true,
      emptyMessage: "当前还没有开启首页精选的文章，所以首页不会展示文章列表。",
    };
  }

  return {
    title: "首页精选",
    posts: featuredPosts,
    isEmpty: false,
  };
}
