export function Spinner() {
  return (
    <div className="w-5 h-5 border-2 border-[#2e3248] border-t-[#6c63ff] rounded-full animate-spin">
      <style>{`@keyframes spin { to { transform: rotate(360deg); } }`}</style>
    </div>
  );
}