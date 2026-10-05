const variants = {
  primary: "bg-pg-ciano-fundo text-pg-branco hover:brightness-90",
  ghost: "bg-transparent border border-pg-borda text-pg-ardosia hover:bg-pg-ciano-nevoa",
  destructive: "bg-transparent border border-pg-down text-pg-down hover:bg-pg-down-bg",
  danger: "bg-pg-down text-pg-branco hover:brightness-90",
};

export function Button({
  children,
  variant = "primary",
  onClick,
  disabled,
  type = "button",
  className = "",
}) {
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={`inline-flex items-center justify-center px-4 py-2 rounded-pg-control text-pg-ui font-medium transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed ${variants[variant]} ${className}`}
    >
      {children}
    </button>
  );
}
