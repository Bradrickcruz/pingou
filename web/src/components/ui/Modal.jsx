import { useEffect } from "react";
import { X } from "lucide-react";

export function Modal({ title, onClose, children }) {
  useEffect(() => {
    const handler = (e) => e.key === "Escape" && onClose();
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onClose]);

  return (
    <div
      className="fixed inset-0 bg-[rgba(30,43,51,0.5)] flex items-center justify-center z-[1000]"
      onClick={onClose}
    >
      <div
        className="w-full max-w-[520px] max-h-[90vh] overflow-y-auto p-7 rounded-pg-card border border-pg-borda bg-pg-branco"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex justify-between items-center mb-5">
          <h2 className="text-pg-h3 font-medium text-pg-ardosia my-0">{title}</h2>
          <button
            onClick={onClose}
            aria-label="Fechar"
            className="text-pg-neblina hover:text-pg-ardosia leading-none rounded-pg-control focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
          >
            <X size={18} aria-hidden="true" />
          </button>
        </div>
        {children}
      </div>
    </div>
  );
}
