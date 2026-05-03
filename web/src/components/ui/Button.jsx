import { tokens as t } from "../../theme/tokens";

const variants = {
  primary: {
    background: t.colors.primary,
    color: "#fff",
    hover: t.colors.primaryHov,
  },
  danger: {
    background: t.colors.danger,
    color: "#fff",
    hover: "#b91c1c",
  },
  ghost: {
    background: "transparent",
    color: t.colors.textMuted,
    hover: t.colors.surfaceAlt,
  },
};

export function Button({
  children,
  variant = "primary",
  onClick,
  disabled,
  className = "",
}) {
  const v = variants[variant];
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`px-4 py-2 text-[13px] font-semibold rounded transition-colors duration-150 ${className}`}
      style={{
        background: v.background,
        color: v.color,
        borderRadius: t.radius.sm,
        opacity: disabled ? 0.5 : 1,
      }}
      onMouseEnter={(e) => (e.currentTarget.style.background = v.hover)}
      onMouseLeave={(e) => (e.currentTarget.style.background = v.background)}
    >
      {children}
    </button>
  );
}