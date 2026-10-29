import { useState } from "react";
import { Download, Lock } from "lucide-react";
import { settingsApi } from "../../api/settings";
import { backupFilename, formatBytes } from "../../lib/format";
import { Card } from "./Card";

// O tamanho é o do banco: o backup (VACUUM INTO) sai um pouco menor, daí o "≈".
export function BackupCard({ databaseSize }) {
  const [downloading, setDownloading] = useState(false);
  const [error, setError] = useState(null);
  const filename = backupFilename();

  const download = async () => {
    setDownloading(true);
    setError(null);
    try {
      const blob = await settingsApi.export();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
    } catch (e) {
      setError(e.message);
    } finally {
      setDownloading(false);
    }
  };

  return (
    <Card
      title="Backup"
      description="Baixa uma cópia completa do banco SQLite: monitores, checks, incidentes e configurações."
    >
      <div className="flex gap-2.5 items-start mb-4 px-4 py-3 rounded-pg-control bg-pg-ciano-nevoa text-pg-caption text-pg-ardosia">
        <Lock size={16} aria-hidden="true" className="shrink-0 mt-px text-pg-ciano-fundo" />
        <p className="m-0">
          O arquivo inclui a URL do webhook e os tokens de ping dos monitores Heartbeat. Guarde com o mesmo cuidado que
          a API Key.
        </p>
      </div>
      <div className="flex flex-wrap items-center gap-3">
        <button
          type="button"
          onClick={download}
          disabled={downloading}
          className="inline-flex items-center gap-2 px-4 py-2.5 rounded-pg-control border border-pg-ardosia bg-pg-branco text-pg-ui font-medium text-pg-ardosia hover:bg-pg-nuvem disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
        >
          <Download size={16} aria-hidden="true" />
          {downloading ? "Gerando…" : "Baixar backup (.db)"}
        </button>
        <span className="font-pg-mono text-pg-caption text-pg-neblina break-all">
          {filename}
          {databaseSize != null && ` · ≈ ${formatBytes(databaseSize)}`}
        </span>
      </div>
      {error && (
        <p role="alert" className="mt-2 mb-0 text-pg-ui text-pg-down">
          {error}
        </p>
      )}
    </Card>
  );
}
