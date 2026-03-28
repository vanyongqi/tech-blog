import { useEffect } from "react";
import type { AdminEditorNotice } from "../lib/adminEditorNotice";

type AdminSaveNoticeBannerProps = {
  notice: AdminEditorNotice;
  onDismiss?: () => void;
};

export function AdminSaveNoticeBanner({ notice, onDismiss }: AdminSaveNoticeBannerProps) {
  useEffect(() => {
    if (!onDismiss) {
      return;
    }
    const timer = window.setTimeout(() => {
      onDismiss();
    }, 3200);

    return () => {
      window.clearTimeout(timer);
    };
  }, [notice, onDismiss]);

  return (
    <div
      className={`admin-save-notice admin-save-notice-${notice.type}`}
      role={notice.type === "error" ? "alert" : "status"}
      aria-live="polite"
    >
      <p>{notice.message}</p>
      {onDismiss ? (
        <button type="button" className="admin-save-notice-close" onClick={onDismiss} aria-label="关闭提示">
          关闭
        </button>
      ) : null}
    </div>
  );
}
