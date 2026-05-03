import { useEffect } from "react";
import { tokens as t } from "../../theme/tokens";

export function Modal({ title, onClose, children }) {
  useEffect(() => {
    const handler = (e) => e.key === "Escape" && onClose();
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onClose]);

  return (
    <div
      className="fixed inset-0 bg-black/60 flex items-center justify-center z-[1000]"
      onClick={onClose}
    >
      <div
        className="bg-[var(--bg)] border border-[var(--border)] rounded-lg p-7 w-full max-w-[520px] max-h-[90vh] overflow-y-auto"
        style={{
          background: t.colors.surface,
          borderColor: t.colors.border,
          borderRadius: t.radius.lg,
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex justify-between items-center mb-5">
          <h2 className="text-base font-semibold text-[var(--text-h)]">{title}</h2>
          <button
            onClick={onClose}
            className="text-[var(--text)] text-xl leading-none hover:text-[var(--text-h)]"
            style={{
              color: t.colors.textMuted,
            }}
          >
            ×
          </button>
        </div>
        {children}
      </div>
    </div>
  );
}