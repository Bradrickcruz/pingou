import { useState } from "react";
import { Check, LogOut } from "lucide-react";
import { useSettings } from "../hooks/useSettings";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import { WebhookCard } from "../components/settings/WebhookCard";
import { EmailCard } from "../components/settings/EmailCard";
import { NtfyCard } from "../components/settings/NtfyCard";
import { MaintenanceNoticesCard } from "../components/settings/MaintenanceNoticesCard";
import { RetentionCard } from "../components/settings/RetentionCard";
import { TimezoneCard } from "../components/settings/TimezoneCard";
import { useTimeZone } from "../hooks/timeZoneContext";
import { BackupCard } from "../components/settings/BackupCard";
import { ThemeCard } from "../components/settings/ThemeCard";
import { useTheme } from "../hooks/themeContext";
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
  const { setTimeZone } = useTimeZone();
  const { preference: themePreference, setPreference: setThemePreference } =
    useTheme();
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
        <p role="alert" className="mt-4 text-pg-ui text-pg-status-down">
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
  const emailRecipients = current.email_recipients ?? "";
  const ntfyTopic = current.ntfy_topic ?? "";
  const timezone = current.timezone ?? "UTC";
  const maintenanceNotices = current.maintenance_notices ?? false;
  const dirty =
    webhookURL !== (settings.webhook_url ?? "") ||
    emailRecipients !== (settings.email_recipients ?? "") ||
    ntfyTopic !== (settings.ntfy_topic ?? "") ||
    timezone !== (settings.timezone ?? "UTC") ||
    maintenanceNotices !== (settings.maintenance_notices ?? false) ||
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
      // o valor exibido depois vem da resposta, já normalizado pelo backend
      const updated = await update({
        webhook_url: webhookURL,
        retention_days: current.retention_days,
        timezone: timezone.trim(),
        maintenance_notices: maintenanceNotices,
        // com o e-mail bloqueado, mandar o campo daria 400
        ...(settings.email_available && { email_recipients: emailRecipients }),
        // idem para o ntfy sem PINGOU_NTFY_URL
        ...(settings.ntfy_available && { ntfy_topic: ntfyTopic }),
      });
      // o painel inteiro passa a usar o fuso salvo, sem recarregar
      setTimeZone(updated.timezone);
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
      <p className="hidden md:block mt-1 mb-6 text-pg-ui text-pg-texto-2">
        Valem para todos os monitores.
      </p>

      <form
        id="settings-form"
        onSubmit={handleSubmit}
        className="flex flex-col gap-4 md:gap-6"
      >
        <div className="flex flex-col gap-3 md:gap-4">
          <p className="m-0 text-pg-ui text-pg-texto-2">
            <strong className="font-medium text-pg-texto">Notificações:</strong>{" "}
            os canais ligados recebem os mesmos eventos.
          </p>
          <WebhookCard
            value={webhookURL}
            onChange={(v) => set("webhook_url", v)}
          />
          <EmailCard
            value={emailRecipients}
            onChange={(v) => set("email_recipients", v)}
            available={settings.email_available}
          />
          <NtfyCard
            value={ntfyTopic}
            onChange={(v) => set("ntfy_topic", v)}
            available={settings.ntfy_available}
          />
          <MaintenanceNoticesCard
            checked={maintenanceNotices}
            onChange={(v) => set("maintenance_notices", v)}
          />
        </div>
        <TimezoneCard value={timezone} onChange={(v) => set("timezone", v)} />
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
              className="hidden md:block m-0 mr-auto text-pg-ui text-pg-status-down"
            >
              {error}
            </p>
          )}
          <span
            role="status"
            className="text-pg-ui inline-flex items-center gap-1 text-pg-status-up"
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

      {/* fora do formulário: o tema é por navegador e aplica na hora */}
      <div className="mt-4 md:mt-6">
        <ThemeCard preference={themePreference} onChange={setThemePreference} />
      </div>

      <div className="mt-4 md:mt-6">
        <BackupCard databaseSize={size} />
      </div>

      <Button
        variant="ghost"
        onClick={onLogout}
        className="md:hidden w-full min-h-[44px] mt-6 gap-2 bg-pg-superficie"
      >
        <LogOut size={16} aria-hidden="true" />
        Sair
      </Button>
      <p className="md:hidden mt-3 text-center text-pg-caption text-pg-texto-2">
        pingou {info?.version ?? "v1.0"} ·{" "}
        <a
          href={REPO}
          target="_blank"
          rel="noreferrer"
          className="text-pg-link underline"
        >
          GitHub
        </a>{" "}
        ·{" "}
        <a
          href={`${REPO}/issues/new`}
          target="_blank"
          rel="noreferrer"
          className="text-pg-link underline"
        >
          Relatar problema
        </a>
      </p>

      <footer className="hidden md:flex mt-6 pt-3 border-t border-pg-linha flex-wrap justify-between gap-x-5 gap-y-2 text-pg-caption text-pg-texto-2">
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
            className="text-pg-link underline"
          >
            Código no GitHub
          </a>
          <a
            href={`${REPO}/issues/new`}
            target="_blank"
            rel="noreferrer"
            className="text-pg-link underline"
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
