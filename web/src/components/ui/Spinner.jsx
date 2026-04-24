export function Spinner() {
  return (
    <div
      style={{
        width: "20px",
        height: "20px",
        border: "2px solid #2e3248",
        borderTop: "2px solid #6c63ff",
        borderRadius: "50%",
        animation: "spin 0.7s linear infinite",
      }}
    >
      <style>{`@keyframes spin { to { transform: rotate(360deg); } }`}</style>
    </div>
  );
}
