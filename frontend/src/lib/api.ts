import type {
  AdminPostResponse,
  AdminPostsResponse,
  AdminProjectResponse,
  AdminProjectsResponse,
  AdminSaveProjectRequest,
  AdminSavePostRequest,
  AdminSaveVideoRequest,
  AdminSessionResponse,
  AdminSuggestThumbnailResponse,
  AdminVideoResponse,
  AdminVideosResponse,
  CommentResponse,
  HomeResponse,
  LikeResponse,
  PostResponse,
  PostsResponse,
} from "../types";

const API_BASE = import.meta.env.VITE_API_BASE ?? "";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers);
  if (init?.body !== undefined && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: "include",
    headers,
  });
  if (!response.ok) {
    let message = `request failed with status ${response.status}`;
    try {
      const errorBody = (await response.json()) as { message?: string };
      if (errorBody.message) {
        message = errorBody.message;
      }
    } catch {
      // Keep fallback message.
    }
    throw new Error(message);
  }
  if (response.status === 204) {
    return undefined as T;
  }
  return response.json() as Promise<T>;
}

export function getHomeData() {
  return request<HomeResponse>("/api/home");
}

export function getAllPosts() {
  return request<PostsResponse>("/api/posts");
}

export function getPost(slug: string) {
  return request<PostResponse>(`/api/posts/${slug}`);
}

export function createComment(slug: string, content: string) {
  return request<CommentResponse>(`/api/posts/${slug}/comments`, {
    method: "POST",
    body: JSON.stringify({ content }),
  });
}

export function toggleLike(slug: string) {
  return request<LikeResponse>(`/api/posts/${slug}/likes`, {
    method: "POST",
    body: JSON.stringify({}),
  });
}

export function loginAdmin(username: string, password: string) {
  return request<AdminSessionResponse>("/api/admin/login", {
    method: "POST",
    body: JSON.stringify({ username, password }),
  });
}

export function logoutAdmin() {
  return request<AdminSessionResponse>("/api/admin/logout", {
    method: "POST",
    body: JSON.stringify({}),
  });
}

export function getAdminSession() {
  return request<AdminSessionResponse>("/api/admin/session");
}

export function getAdminPosts() {
  return request<AdminPostsResponse>("/api/admin/posts");
}

export function getAdminPost(slug: string) {
  return request<AdminPostResponse>(`/api/admin/posts/${slug}`);
}

export function createAdminPost(payload: AdminSavePostRequest) {
  return request<AdminPostResponse>("/api/admin/posts", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateAdminPost(slug: string, payload: AdminSavePostRequest) {
  return request<AdminPostResponse>(`/api/admin/posts/${slug}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export function deleteAdminPost(slug: string) {
  return request<void>(`/api/admin/posts/${slug}`, {
    method: "DELETE",
  });
}

export function getAdminProjects() {
  return request<AdminProjectsResponse>("/api/admin/projects");
}

export function getAdminProject(id: number) {
  return request<AdminProjectResponse>(`/api/admin/projects/${id}`);
}

export function createAdminProject(payload: AdminSaveProjectRequest) {
  return request<AdminProjectResponse>("/api/admin/projects", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateAdminProject(id: number, payload: AdminSaveProjectRequest) {
  return request<AdminProjectResponse>(`/api/admin/projects/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export function deleteAdminProject(id: number) {
  return request<void>(`/api/admin/projects/${id}`, {
    method: "DELETE",
  });
}

export function getAdminVideos() {
  return request<AdminVideosResponse>("/api/admin/videos");
}

export function getAdminVideo(id: number) {
  return request<AdminVideoResponse>(`/api/admin/videos/${id}`);
}

export function createAdminVideo(payload: AdminSaveVideoRequest) {
  return request<AdminVideoResponse>("/api/admin/videos", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateAdminVideo(id: number, payload: AdminSaveVideoRequest) {
  return request<AdminVideoResponse>(`/api/admin/videos/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export function deleteAdminVideo(id: number) {
  return request<void>(`/api/admin/videos/${id}`, {
    method: "DELETE",
  });
}

export function suggestAdminVideoThumbnail(url: string) {
  return request<AdminSuggestThumbnailResponse>("/api/admin/videos/thumbnail", {
    method: "POST",
    body: JSON.stringify({ url }),
  });
}

export function formatDate(date: string) {
  return new Intl.DateTimeFormat("zh-CN", {
    year: "numeric",
    month: "short",
    day: "numeric",
  }).format(new Date(date));
}

export function formatDateTime(date: string) {
  return new Intl.DateTimeFormat("zh-CN", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(date));
}
