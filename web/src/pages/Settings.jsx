import { useState } from "react";
import { Check, LogOut } from "lucide-react";
import { useSettings } from "../hooks/useSettings";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import { WebhookCard } from "../components/settings/WebhookCard";
import { RetentionCard } from "../components/settings/RetentionCard";
import { BackupCard } from "../components/settings/BackupCard";
import { UnsavedBar } from "../components/settings/UnsavedBar";
import { MobileHeader } from "../components/layout/MobileHeader";

const REPO = "https://github.com/Bradrickcruz/pingou";

export function Settings({ onLogout }) {
  const {
    settings,
    info,
    loading,
    error: loadError,
    update,
    reload,
  } = useSettings({ withInfo: true });
  const [form, setForm] = useState(null);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState(null);

  if (loading)
    return (
      <>
        <MobileHeader title="Configurações" />
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      </>
    );

  if (!settings)
    return (
      <div className="max-w-[760px]">
        <MobileHeader title="Configurações" />
        <h1 className="hidden md:block text-pg-h1 my-0">Configurações</h1>
        <p role="alert" className="mt-4 text-pg-ui text-pg-down">
          Não foi possível carregar as configurações
          {loadError ? `: ${loadError}` : "."}
        </p>
        <Button variant="ghost" onClick={reload}>
          Tentar de novo
        </Button>
      </div>
    );

  const current = form ?? settings;
  const webhookURL = current.webhook_url ?? "";
  const dirty =
    webhookURL !== (settings.webhook_url ?? "") ||
    current.retention_days !== settings.retention_days;
  const size = info?.database_size_bytes;

  const set = (k, v) => {
    setSaved(false);
    setForm((f) => ({ ...(f ?? settings), [k]: v }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!dirty) return;
    setSaving(true);
    setError(null);
    try {
      await update({
        webhook_url: webhookURL,
        retention_days: current.retention_days,
      });
      setForm(null);
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (e) {
      setError(e.message);
    } finally {
      setSaving(false);
    }
  };

  const discard = () => {
    setForm(null);
    setError(null);
  };

  return (
    <div className={`max-w-[760px] ${dirty ? "pb-20 md:pb-0" : ""}`}>
      <MobileHeader title="Configurações" />
      <h1 className="hidden md:block text-pg-h1 my-0">Configurações</h1>
      <p className="hidden md:block mt-1 mb-6 text-pg-ui text-pg-neblina">
        Valem para todos os monitores.
      </p>

      <form
        id="settings-form"
        onSubmit={handleSubmit}
        className="flex flex-col gap-4 md:gap-6"
      >
        <WebhookCard
          value={webhookURL}
          onChange={(v) => set("webhook_url", v)}
        />
        <RetentionCard
          value={current.retention_days}
          savedValue={settings.retention_days}
          onChange={(v) => set("retention_days", v)}
          databaseSize={size}
        />

        <div className="flex flex-wrap items-center justify-end gap-3">
          {error && (
            <p
              role="alert"
              className="hidden md:block m-0 mr-auto text-pg-ui text-pg-down"
            >
              {error}
            </p>
          )}
          <span
            role="status"
            className="text-pg-ui inline-flex items-center gap-1 text-pg-up"
          >
            {saved && (
              <>
                <Check size={14} aria-hidden="true" />
                Salvo
              </>
            )}
          </span>
          <Button
            type="submit"
            disabled={!dirty || saving}
            className="hidden md:inline-flex"
          >
            {saving ? "Salvando…" : "Salvar alterações"}
          </Button>
        </div>
      </form>

      <div className="mt-4 md:mt-6">
        <BackupCard databaseSize={size} />
      </div>

      <Button
        variant="ghost"
        onClick={onLogout}
        className="md:hidden w-full min-h-[44px] mt-6 gap-2 bg-pg-branco"
      >
        <LogOut size={16} aria-hidden="true" />
        Sair
      </Button>
      <p className="md:hidden mt-3 text-center text-pg-caption text-pg-neblina">
        pingou {info?.version ?? "v1.0"} ·{" "}
        <a
          href={REPO}
          target="_blank"
          rel="noreferrer"
          className="text-pg-ciano-fundo underline"
        >
          GitHub
        </a>{" "}
        ·{" "}
        <a
          href={`${REPO}/issues/new`}
          target="_blank"
          rel="noreferrer"
          className="text-pg-ciano-fundo underline"
        >
          Relatar problema
        </a>
      </p>

      <footer className="hidden md:flex mt-6 pt-3 border-t border-pg-borda flex-wrap justify-between gap-x-5 gap-y-2 text-pg-caption text-pg-neblina">
        <span>
          {info && `pingou ${info.version} · `}
          {info &&
            info.commit &&
            info.commit !== "none" &&
            `commit ${info.commit} · `}
          Apache 2.0
        </span>
        <span className="flex gap-5">
          <a
            href={REPO}
            target="_blank"
            rel="noreferrer"
            className="text-pg-ciano-fundo underline"
          >
            Código no GitHub
          </a>
          <a
            href={`${REPO}/issues/new`}
            target="_blank"
            rel="noreferrer"
            className="text-pg-ciano-fundo underline"
          >
            Relatar um problema
          </a>
        </span>
      </footer>

      {dirty && (
        <UnsavedBar
          formId="settings-form"
          onDiscard={discard}
          saving={saving}
          error={error}
        />
      )}
    </div>
  );
}
