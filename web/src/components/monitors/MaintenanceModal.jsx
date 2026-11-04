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
import { defaultRange, toApiDate, validateWindow } from "../../lib/maintenance";
import { getTimeZone } from "../../lib/timezone";

// Agendar uma janela para um ou mais monitores. `monitors` é a lista para
// escolher (sem ela, o modal busca); `initialIds` vem marcado (detalhe).
export function MaintenanceModal({
  monitors,
  initialIds = [],
  onClose,
  onCreated,
}) {
  const [list, setList] = useState(monitors ?? null);
  const [ids, setIds] = useState(initialIds);
  const [query, setQuery] = useState("");
  const [range] = useState(() => defaultRange());
  const [starts, setStarts] = useState(range.starts);
  const [ends, setEnds] = useState(range.ends);
  const [reason, setReason] = useState("");
  const [errors, setErrors] = useState({});
  const [submitError, setSubmitError] = useState(null);
  const [saving, setSaving] = useState(false);
  const formRef = useRef(null);
  const uid = useId();
  const fieldId = (name) => `${uid}-${name}`;
  const errorId = (name) => `${uid}-${name}-error`;
  const zone = getTimeZone();

  useEffect(() => {
    if (monitors) return;
    let alive = true;
    monitorsApi
      .list({ limit: 100 })
      .then((r) => alive && setList(r.data))
      .catch((e) => alive && setSubmitError(e.message));
    return () => {
      alive = false;
    };
  }, [monitors]);

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

  const submit = async (e) => {
    e.preventDefault();
    if (saving) return;
    setSubmitError(null);
    const errs = validateWindow({ ids, starts, ends });
    setErrors(errs);
    if (Object.keys(errs).length > 0) {
      focusFirstError(errs);
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
    `w-full px-3 py-2.5 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border ${
      invalid ? "border-pg-down" : "border-pg-borda"
    }`;
  const fieldError = (name) =>
    errors[name] && (
      <p id={errorId(name)} className="mt-1 mb-0 text-pg-caption text-pg-down">
        {errors[name]}
      </p>
    );
  const dateField = (name, label, value, set) => (
    <div className="min-w-0">
      <label
        htmlFor={fieldId(name)}
        className="block mb-2 text-pg-ui font-medium text-pg-ardosia"
      >
        {label} <span className="font-normal text-pg-neblina">({zone})</span>
      </label>
      <input
        id={fieldId(name)}
        type="datetime-local"
        value={value}
        onChange={(e) => set(e.target.value)}
        aria-invalid={errors[name] ? true : undefined}
        aria-describedby={errors[name] ? errorId(name) : undefined}
        className={inputCls(errors[name])}
      />
      {fieldError(name)}
    </div>
  );

  return (
    <Modal
      title="Agendar manutenção"
      subtitle="Durante a janela, nenhum check roda, nenhum incidente abre e ninguém é avisado."
      onClose={onClose}
    >
      <form ref={formRef} onSubmit={submit} noValidate>
        <fieldset
          id={fieldId("monitors")}
          className="mb-5 border-0 p-0 m-0 min-w-0"
          aria-invalid={errors.monitors ? true : undefined}
          aria-describedby={errors.monitors ? errorId("monitors") : undefined}
        >
          <legend className="mb-2 text-pg-ui font-medium text-pg-ardosia">
            Monitores{" "}
            <span className="font-normal text-pg-neblina">
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
                    className="absolute left-3 top-1/2 -translate-y-1/2 text-pg-neblina"
                  />
                  <input
                    type="search"
                    aria-label="Filtrar monitores"
                    placeholder="Filtrar por nome ou endereço"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    className="w-full pl-9 pr-3 py-2 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda placeholder:text-pg-neblina"
                  />
                </div>
              )}
              <ul
                className={`list-none m-0 p-1 max-h-52 overflow-y-auto rounded-pg-control border ${
                  errors.monitors ? "border-pg-down" : "border-pg-borda"
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
                            : "cursor-pointer hover:bg-pg-nuvem"
                        }`}
                      >
                        <input
                          type="checkbox"
                          checked={ids.includes(m.id)}
                          disabled={blocked}
                          onChange={() => toggle(m.id)}
                          className="w-4 h-4 flex-shrink-0 accent-pg-ciano-fundo"
                        />
                        <span className="truncate text-pg-ui text-pg-ardosia">
                          {m.name}
                        </span>
                        {blocked && (
                          <span className="ml-auto flex-shrink-0 text-pg-caption text-pg-neblina">
                            {!m.enabled ? "pausado" : "já agendada"}
                          </span>
                        )}
                      </label>
                    </li>
                  );
                })}
                {visible.length === 0 && (
                  <li className="px-2 py-3 text-pg-caption text-pg-neblina">
                    Nenhum monitor encontrado.
                  </li>
                )}
              </ul>
            </>
          )}
          {fieldError("monitors")}
        </fieldset>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-5">
          {dateField("starts", "Início", starts, setStarts)}
          {dateField("ends", "Fim", ends, setEnds)}
        </div>

        <ReasonField
          label="Motivo"
          value={reason}
          onChange={setReason}
          placeholder="Ex.: atualização do banco de dados"
        />

        {submitError && (
          <p role="alert" className="-mt-2 mb-4 text-pg-caption text-pg-down">
            {submitError}
          </p>
        )}

        <div className="flex gap-2.5 justify-end">
          <Button variant="ghost" onClick={onClose}>
            Cancelar
          </Button>
          <Button type="submit" disabled={saving || list == null}>
            {saving ? "Agendando..." : "Agendar"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
