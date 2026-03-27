import { describe, expect, it } from "vitest";
import { getAdminEditorSaveSuccessMessage } from "./adminEditorNotice";

describe("getAdminEditorSaveSuccessMessage", () => {
  it("returns post save notice", () => {
    expect(getAdminEditorSaveSuccessMessage("post", true)).toBe("保存成功，文章内容已更新。");
  });

  it("returns project create notice", () => {
    expect(getAdminEditorSaveSuccessMessage("project", false)).toBe("创建成功，项目内容已更新。");
  });

  it("returns video save notice", () => {
    expect(getAdminEditorSaveSuccessMessage("video", true)).toBe("保存成功，视频内容已更新。");
  });
});
