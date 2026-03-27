export type AdminEditorNotice = {
  type: "success" | "error";
  message: string;
};

export function getAdminEditorSaveSuccessMessage(entity: "post" | "project" | "video", isEditMode: boolean) {
  const action = isEditMode ? "保存成功" : "创建成功";

  switch (entity) {
    case "post":
      return `${action}，文章内容已更新。`;
    case "project":
      return `${action}，项目内容已更新。`;
    case "video":
      return `${action}，视频内容已更新。`;
  }
}
