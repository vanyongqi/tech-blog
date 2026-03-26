import { Route, Routes } from "react-router-dom";
import { AdminShell } from "./components/AdminShell";
import { PageShell } from "./components/PageShell";
import { AdminLoginPage } from "./pages/AdminLoginPage";
import { AdminPostEditorPage } from "./pages/AdminPostEditorPage";
import { AdminPostsPage } from "./pages/AdminPostsPage";
import { AdminProjectEditorPage } from "./pages/AdminProjectEditorPage";
import { AdminProjectsPage } from "./pages/AdminProjectsPage";
import { AdminVideoEditorPage } from "./pages/AdminVideoEditorPage";
import { AdminVideosPage } from "./pages/AdminVideosPage";
import { AboutPage } from "./pages/AboutPage";
import { ArticlesPage } from "./pages/ArticlesPage";
import { HomePage } from "./pages/HomePage";
import { PostDetailPage } from "./pages/PostDetailPage";
import { ProjectsPage } from "./pages/ProjectsPage";
import { VideosPage } from "./pages/VideosPage";

export default function App() {
  return (
    <Routes>
      <Route path="/admin/login" element={<AdminLoginPage />} />
      <Route path="/admin" element={<AdminShell />}>
        <Route index element={<AdminPostsPage />} />
        <Route path="posts/new" element={<AdminPostEditorPage />} />
        <Route path="posts/:slug/edit" element={<AdminPostEditorPage />} />
        <Route path="projects" element={<AdminProjectsPage />} />
        <Route path="projects/new" element={<AdminProjectEditorPage />} />
        <Route path="projects/:id/edit" element={<AdminProjectEditorPage />} />
        <Route path="videos" element={<AdminVideosPage />} />
        <Route path="videos/new" element={<AdminVideoEditorPage />} />
        <Route path="videos/:id/edit" element={<AdminVideoEditorPage />} />
      </Route>
      <Route element={<PageShell />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/about" element={<AboutPage />} />
        <Route path="/articles" element={<ArticlesPage />} />
        <Route path="/projects" element={<ProjectsPage />} />
        <Route path="/videos" element={<VideosPage />} />
        <Route path="/posts/:slug" element={<PostDetailPage />} />
      </Route>
    </Routes>
  );
}
