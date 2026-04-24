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
  style = {},
}) {
  const v = variants[variant];
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      style={{
        background: v.background,
        color: v.color,
        borderRadius: t.radius.sm,
        padding: "8px 16px",
        fontWeight: 600,
        fontSize: "13px",
        opacity: disabled ? 0.5 : 1,
        transition: "background 0.15s",
        ...style,
      }}
      onMouseEnter={(e) => (e.currentTarget.style.background = v.hover)}
      onMouseLeave={(e) => (e.currentTarget.style.background = v.background)}
    >
      {children}
    </button>
  );
}
