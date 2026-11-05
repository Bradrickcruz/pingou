import { useEffect, useId, useMemo, useRef, useState } from "react";
import { Search } from "lucide-react";
import { toast } from "react-toastify";
import { maintenanceApi } from "../../api/maintenance";
import { monitorsApi } from "../../api/monitors";
import { Modal } from "../ui/Modal";
import { Button } from "../ui/Button";
import { Spinner } from "../ui/Spinner";
import { ReasonField } from "./MonitorModals";
import { matchesQuery } from "../../lib/monitors";
import {
  defaultRange,
  toApiDate,
  validateWindow,
  windowState,
} from "../../lib/maintenance";
import { getTimeZone, utcToZonedInput } from "../../lib/timezone";

const LOCKED_NOTE = "A manutenção já começou: só o motivo pode mudar.";

// Dois modos (D15):
// - criar: agenda uma janela para um ou mais monitores. `monitors` é a lista
//   para escolher (sem ela, o modal busca); `initialIds` vem marcado (detalhe).
// - editar (`window`): uma janela existente, com o monitor fixo (D13). Janela já
//   iniciada tem início e fim desabilitados: só o motivo muda (D12).
// onStale: a janela mudou no servidor (começou ou sumiu); quem abriu recarrega.
export function MaintenanceModal({
  monitors,
  initialIds = [],
  window: editing,
  onClose,
  onCreated,
  onSaved,
  onStale,
}) {
  const edit = editing != null;
  const [initial] = useState(() =>
    edit
      ? {
          starts: utcToZonedInput(new Date(editing.starts_at)),
          ends: utcToZonedInput(new Date(editing.ends_at)),
          reason: editing.reason ?? "",
        }
      : { ...defaultRange(), reason: "" },
  );
  const [locked, setLocked] = useState(
    () => edit && windowState(editing) === "active",
  );
  const [list, setList] = useState(edit ? [] : (monitors ?? null));
  const [ids, setIds] = useState(initialIds);
  const [query, setQuery] = useState("");
  const [starts, setStarts] = useState(initial.starts);
  const [ends, setEnds] = useState(initial.ends);
  const [reason, setReason] = useState(initial.reason);
  const [errors, setErrors] = useState({});
  const [submitError, setSubmitError] = useState(null);
  const [saving, setSaving] = useState(false);
  const formRef = useRef(null);
  const uid = useId();
  const fieldId = (name) => `${uid}-${name}`;
  const errorId = (name) => `${uid}-${name}-error`;
  const zone = getTimeZone();
  const lockedNoteId = `${uid}-locked`;

  useEffect(() => {
    if (monitors || edit) return;
    let alive = true;
    monitorsApi
      .list({ limit: 100 })
      .then((r) => alive && setList(r.data))
      .catch((e) => alive && setSubmitError(e.message));
    return () => {
      alive = false;
    };
  }, [monitors, edit]);

  const byName = useMemo(
    () =>
      [...(list ?? [])].sort((a, b) => a.name.localeCompare(b.name, "pt-BR")),
    [list],
  );
  const visible = byName.filter((m) => matchesQuery(m, query));
  const toggle = (id) =>
    setIds((cur) =>
      cur.includes(id) ? cur.filter((x) => x !== id) : [...cur, id],
    );

  const focusFirstError = (errs) => {
    const first = ["monitors", "starts", "ends"].find((k) => errs[k]);
    if (!first) return;
    const el =
      first === "monitors"
        ? formRef.current?.querySelector(
            `#${CSS.escape(fieldId("monitors"))} input[type=checkbox]`,
          )
        : formRef.current?.querySelector(`#${CSS.escape(fieldId(first))}`);
    el?.focus();
  };

  // edição: só os campos que mudaram; nada mudou → fecha sem chamar a API
  const saveEdit = async () => {
    const patch = {};
    if (!locked && starts !== initial.starts)
      patch.starts_at = toApiDate(starts);
    if (!locked && ends !== initial.ends) patch.ends_at = toApiDate(ends);
    if (reason.trim() !== initial.reason.trim()) patch.reason = reason.trim();
    if (Object.keys(patch).length === 0) {
      onClose();
      return;
    }
    setSaving(true);
    try {
      await maintenanceApi.update(editing.id, patch);
      toast.success("Manutenção atualizada");
      await onSaved?.();
      onClose();
    } catch (err) {
      if (err.status === 409 && err.data?.code === "MAINTENANCE_LOCKED") {
        setLocked(true);
        setStarts(initial.starts);
        setEnds(initial.ends);
        setErrors({});
        setSubmitError(
          "Esta manutenção começou enquanto você editava, então o horário não pode mais mudar. O motivo ainda pode.",
        );
        void onStale?.();
      } else if (err.status === 404) {
        setSubmitError("Esta manutenção não existe mais.");
        void onStale?.();
      } else {
        setSubmitError(err.message);
      }
    } finally {
      setSaving(false);
    }
  };

  const submit = async (e) => {
    e.preventDefault();
    if (saving) return;
    setSubmitError(null);
    const errs = validateWindow({
      ids,
      starts,
      ends,
      edit,
      initialStarts: initial.starts,
      locked,
    });
    setErrors(errs);
    if (Object.keys(errs).length > 0) {
      focusFirstError(errs);
      return;
    }
    if (edit) {
      await saveEdit();
      return;
    }
    setSaving(true);
    try {
      await maintenanceApi.create({
        monitor_ids: ids,
        starts_at: toApiDate(starts),
        ends_at: toApiDate(ends),
        reason: reason.trim(),
      });
      toast.success(
        ids.length === 1
          ? "Manutenção agendada"
          : `Manutenção agendada para ${ids.length} monitores`,
      );
      await onCreated?.();
      onClose();
    } catch (err) {
      if (err.status === 409 && err.data?.monitor_ids) {
        const names = err.data.monitor_ids.map(
          (id) => byName.find((m) => m.id === id)?.name ?? id,
        );
        setSubmitError(
          `Já ${names.length === 1 ? "tem" : "têm"} manutenção agendada ou em curso: ${names.join(", ")}. Cancele a atual antes de agendar outra.`,
        );
      } else {
        setSubmitError(err.message);
      }
    } finally {
      setSaving(false);
    }
  };

  const inputCls = (invalid) =>
    `w-full px-3 py-2.5 rounded-pg-control text-pg-ui text-pg-texto bg-pg-superficie border ${
      invalid ? "border-pg-status-down" : "border-pg-linha"
    }`;
  const fieldError = (name) =>
    errors[name] && (
      <p
        id={errorId(name)}
        className="mt-1 mb-0 text-pg-caption text-pg-status-down"
      >
        {errors[name]}
      </p>
    );
  const dateField = (name, label, value, set) => (
    <div className="min-w-0">
      <label
        htmlFor={fieldId(name)}
        className="block mb-2 text-pg-ui font-medium text-pg-texto"
      >
        {label} <span className="font-normal text-pg-texto-2">({zone})</span>
      </label>
      <input
        id={fieldId(name)}
        type="datetime-local"
        value={value}
        onChange={(e) => set(e.target.value)}
        disabled={locked}
        aria-invalid={errors[name] ? true : undefined}
        aria-describedby={
          [errors[name] && errorId(name), locked && lockedNoteId]
            .filter(Boolean)
            .join(" ") || undefined
        }
        className={`${inputCls(errors[name])} disabled:bg-pg-fundo disabled:text-pg-texto-2 disabled:cursor-not-allowed`}
      />
      {fieldError(name)}
    </div>
  );

  return (
    <Modal
      title={edit ? "Editar manutenção" : "Agendar manutenção"}
      subtitle="Durante a janela, nenhum check roda, nenhum incidente abre e ninguém é avisado."
      onClose={onClose}
    >
      <form ref={formRef} onSubmit={submit} noValidate>
        {edit ? (
          <p className="mb-5 text-pg-ui text-pg-texto">
            <span className="text-pg-texto-2">Monitor:</span>{" "}
            <strong className="font-medium">
              {editing.monitor_name ?? "—"}
            </strong>
          </p>
        ) : (
          <fieldset
            id={fieldId("monitors")}
            className="mb-5 border-0 p-0 m-0 min-w-0"
            aria-invalid={errors.monitors ? true : undefined}
            aria-describedby={errors.monitors ? errorId("monitors") : undefined}
          >
            <legend className="mb-2 text-pg-ui font-medium text-pg-texto">
              Monitores{" "}
              <span className="font-normal text-pg-texto-2">
                ({ids.length} {ids.length === 1 ? "escolhido" : "escolhidos"})
              </span>
            </legend>
            {list == null ? (
              <div className="flex justify-center py-6">
                <Spinner />
              </div>
            ) : (
              <>
                {byName.length > 6 && (
                  <div className="relative mb-2">
                    <Search
                      size={16}
                      aria-hidden="true"
                      className="absolute left-3 top-1/2 -translate-y-1/2 text-pg-texto-2"
                    />
                    <input
                      type="search"
                      aria-label="Filtrar monitores"
                      placeholder="Filtrar por nome ou endereço"
                      value={query}
                      onChange={(e) => setQuery(e.target.value)}
                      className="w-full pl-9 pr-3 py-2 rounded-pg-control text-pg-ui text-pg-texto bg-pg-superficie border border-pg-linha placeholder:text-pg-texto-2"
                    />
                  </div>
                )}
                <ul
                  className={`list-none m-0 p-1 max-h-52 overflow-y-auto rounded-pg-control border ${
                    errors.monitors
                      ? "border-pg-status-down"
                      : "border-pg-linha"
                  }`}
                >
                  {visible.map((m) => {
                    // ponytail: espelha a regra do backend (pausado ou com janela viva não entra)
                    const blocked = !m.enabled || !!m.maintenance;
                    return (
                      <li key={m.id}>
                        <label
                          className={`flex items-center gap-2.5 px-2 py-2 rounded-pg-control ${
                            blocked
                              ? "cursor-not-allowed opacity-60"
                              : "cursor-pointer hover:bg-pg-hover"
                          }`}
                        >
                          <input
                            type="checkbox"
                            checked={ids.includes(m.id)}
                            disabled={blocked}
                            onChange={() => toggle(m.id)}
                            className="w-4 h-4 flex-shrink-0 accent-pg-link"
                          />
                          <span className="truncate text-pg-ui text-pg-texto">
                            {m.name}
                          </span>
                          {blocked && (
                            <span className="ml-auto flex-shrink-0 text-pg-caption text-pg-texto-2">
                              {!m.enabled ? "pausado" : "já agendada"}
                            </span>
                          )}
                        </label>
                      </li>
                    );
                  })}
                  {visible.length === 0 && (
                    <li className="px-2 py-3 text-pg-caption text-pg-texto-2">
                      Nenhum monitor encontrado.
                    </li>
                  )}
                </ul>
              </>
            )}
            {fieldError("monitors")}
          </fieldset>
        )}

        <div
          className={`grid grid-cols-1 sm:grid-cols-2 gap-4 ${locked ? "mb-2" : "mb-5"}`}
        >
          {dateField("starts", "Início", starts, setStarts)}
          {dateField("ends", "Fim", ends, setEnds)}
        </div>
        {locked && (
          <p id={lockedNoteId} className="mb-5 text-pg-caption text-pg-texto-2">
            {LOCKED_NOTE}
          </p>
        )}

        <ReasonField
          label="Motivo"
          value={reason}
          onChange={setReason}
          placeholder="Ex.: atualização do banco de dados"
        />

        {submitError && (
          <p
            role="alert"
            className="-mt-2 mb-4 text-pg-caption text-pg-status-down"
          >
            {submitError}
          </p>
        )}

        <div className="flex gap-2.5 justify-end">
          <Button variant="ghost" onClick={onClose}>
            Cancelar
          </Button>
          <Button type="submit" disabled={saving || list == null}>
            {edit
              ? saving
                ? "Salvando..."
                : "Salvar"
              : saving
                ? "Agendando..."
                : "Agendar"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
