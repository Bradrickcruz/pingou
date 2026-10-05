import { useId } from "react";

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

export function TimeRangeSlider({ label, value, onChange, min, max, step = 1 }) {
  const id = useId();
  return (
    <div className="mb-3.5">
      <label htmlFor={id} className="block mb-1.5 text-pg-caption font-medium text-pg-neblina">
        {label}
      </label>
      <input
        id={id}
        type="range"
        min={min}
        max={max}
        step={step}
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        className="w-full cursor-pointer accent-pg-ciano-fundo"
      />
      <div className="flex justify-between mt-1">
        <span className="text-pg-caption text-pg-neblina">{formatTime(min)}</span>
        <span className="text-pg-ui font-medium text-pg-ciano-fundo">{formatTime(value)}</span>
        <span className="text-pg-caption text-pg-neblina">{formatTime(max)}</span>
      </div>
    </div>
  );
}
