import { tokens as t } from "../../theme/tokens";

function formatTime(seconds) {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;

  if (h > 0) {
    return `${h}h ${m}m ${s}s`;
  } else if (m > 0) {
    return `${m}m ${s}s`;
  } else {
    return `${s}s`;
  }
}

export function TimeRangeSlider({ label, value, onChange, min, max, step = 1, unit = "s" }) {
  return (
    <div className="mb-3.5">
      <label
        className="block mb-1.5 text-xs font-medium"
        style={{
          color: t.colors.textMuted,
        }}
      >
        {label}
      </label>
      <input
        type="range"
        min={min}
        max={max}
        step={step}
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        className="w-full h-2 rounded-lg appearance-none cursor-pointer"
        style={{
          background: t.colors.surface,
        }}
      />
      <div className="flex justify-between mt-1">
        <span className="text-xs" style={{ color: t.colors.textMuted }}>
          {formatTime(min)}
        </span>
        <span className="text-sm font-medium" style={{ color: t.colors.primary }}>
          {formatTime(value)}
        </span>
        <span className="text-xs" style={{ color: t.colors.textMuted }}>
          {formatTime(max)}
        </span>
      </div>
    </div>
  );
}