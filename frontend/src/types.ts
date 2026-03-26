export interface SiteStat {
  label: string;
  value: string;
}

export interface SocialLink {
  label: string;
  url: string;
}

export interface SitePayload {
  name: string;
  headline: string;
  intro: string;
  location: string;
  domain: string;
  email: string;
  motto: string;
  techStack: string[];
  stats: SiteStat[];
  socialLinks: SocialLink[];
}

export interface PostSummaryPayload {
  slug: string;
  title: string;
  summary: string;
  category: string;
  readTime: string;
  heroNote: string;
  coverLabel: string;
  tags: string[];
  featured: boolean;
  publishedAt: string;
  likeCount: number;
  commentCount: number;
}

export interface ContentBlock {
  kind: "paragraph" | "heading" | "quote" | "list" | "video";
  title?: string;
  text?: string;
  url?: string;
  items?: string[];
}

export interface PostDetailPayload extends PostSummaryPayload {
  blocks: ContentBlock[];
  likedByVisitor: boolean;
  comments: CommentPayload[];
}

export interface CommentPayload {
  id: number;
  authorName: string;
  content: string;
  createdAt: string;
}

export interface VisitorPayload {
  displayName: string;
}

export interface AdminSessionPayload {
  authenticated: boolean;
  username: string;
}

export interface ProjectPayload {
  name: string;
  summary: string;
  status: string;
  link: string;
  imageUrl: string;
  accent: "ember" | "forest" | "ink" | string;
  techStack: string[];
}

export interface VideoPayload {
  id: number;
  title: string;
  description: string;
  url: string;
  thumbnailUrl: string;
  publishedAt: string;
}

export interface TimelineEntryPayload {
  period: string;
  title: string;
  description: string;
}

export interface HomeResponse {
  site: SitePayload;
  featuredPosts: PostSummaryPayload[];
  recentPosts: PostSummaryPayload[];
  projects: ProjectPayload[];
  videos: VideoPayload[];
  timeline: TimelineEntryPayload[];
}

export interface PostsResponse {
  posts: PostSummaryPayload[];
}

export interface PostResponse {
  post: PostDetailPayload;
  visitor: VisitorPayload;
}

export interface AdminContentBlock {
  kind: string;
  title?: string;
  text?: string;
  url?: string;
  items?: string[];
}

export interface AdminPostSummaryPayload {
  slug: string;
  title: string;
  summary: string;
  category: string;
  readTime: string;
  coverLabel: string;
  tags: string[];
  featured: boolean;
  publishedAt: string;
  likeCount: number;
  commentCount: number;
}

export interface AdminPostPayload {
  slug: string;
  title: string;
  summary: string;
  category: string;
  readTime: string;
  heroNote: string;
  coverLabel: string;
  tags: string[];
  featured: boolean;
  publishedAt: string;
  blocks: AdminContentBlock[];
  likeCount: number;
  commentCount: number;
}

export interface AdminSavePostRequest {
  slug: string;
  title: string;
  summary: string;
  category: string;
  readTime: string;
  heroNote: string;
  coverLabel: string;
  tags: string[];
  featured: boolean;
  publishedAt: string;
  blocks: AdminContentBlock[];
}

export interface CommentResponse {
  comment: CommentPayload;
  commentCount: number;
  visitor: VisitorPayload;
}

export interface LikeResponse {
  likeCount: number;
  liked: boolean;
  visitor: VisitorPayload;
}

export interface AdminSessionResponse {
  session: AdminSessionPayload;
}

export interface AdminPostsResponse {
  posts: AdminPostSummaryPayload[];
}

export interface AdminPostResponse {
  post: AdminPostPayload;
}

export interface AdminProjectPayload {
  id: number;
  name: string;
  summary: string;
  status: string;
  link: string;
  imageUrl: string;
  accent: string;
  techStack: string[];
}

export interface AdminSaveProjectRequest {
  name: string;
  summary: string;
  status: string;
  link: string;
  imageUrl: string;
  accent: string;
  techStack: string[];
}

export interface AdminProjectsResponse {
  projects: AdminProjectPayload[];
}

export interface AdminProjectResponse {
  project: AdminProjectPayload;
}

export interface AdminVideoPayload {
  id: number;
  title: string;
  description: string;
  url: string;
  thumbnailUrl: string;
  publishedAt: string;
}

export interface AdminSaveVideoRequest {
  title: string;
  description: string;
  url: string;
  thumbnailUrl: string;
  publishedAt: string;
}

export interface AdminVideosResponse {
  videos: AdminVideoPayload[];
}

export interface AdminVideoResponse {
  video: AdminVideoPayload;
}

export interface AdminSuggestThumbnailResponse {
  thumbnailUrl: string;
}
