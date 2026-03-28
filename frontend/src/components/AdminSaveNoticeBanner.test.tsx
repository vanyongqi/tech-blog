import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import { AdminSaveNoticeBanner } from "./AdminSaveNoticeBanner";

describe("AdminSaveNoticeBanner", () => {
  it("renders success notice as status banner", () => {
    const html = renderToStaticMarkup(
      <AdminSaveNoticeBanner notice={{ type: "success", message: "保存成功，文章内容已更新。" }} />,
    );

    expect(html).toContain("admin-save-notice");
    expect(html).toContain("admin-save-notice-success");
    expect(html).toContain('role="status"');
    expect(html).toContain("保存成功，文章内容已更新。");
  });

  it("renders error notice as alert banner", () => {
    const html = renderToStaticMarkup(
      <AdminSaveNoticeBanner notice={{ type: "error", message: "保存失败" }} />,
    );

    expect(html).toContain("admin-save-notice-error");
    expect(html).toContain('role="alert"');
    expect(html).toContain("保存失败");
  });
});
