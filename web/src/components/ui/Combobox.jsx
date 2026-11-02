import { useEffect, useId, useState } from "react";

// Campo digitável com lista que filtra enquanto se digita (padrão ARIA 1.2
// "combobox com listbox"). Só aceita valores da lista: ao sair do campo sem
// escolher, volta ao valor atual. `filter(options, texto)` decide o que mostrar.
export function Combobox({
  id,
  value,
  options,
  onChange,
  filter,
  emptyText = "Nada encontrado",
  ...rest
}) {
  const listId = useId();
  const [text, setText] = useState(null); // null = não está digitando
  const [open, setOpen] = useState(false);
  const [active, setActive] = useState(0);

  const shown = text == null ? options : filter(options, text);
  const optionId = (i) => `${listId}-${i}`;

  useEffect(() => {
    if (open) {
      document
        .getElementById(optionId(active))
        ?.scrollIntoView({ block: "nearest" });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, active]);

  const openList = () => {
    if (open) return;
    setActive(Math.max(0, options.indexOf(value)));
    setOpen(true);
  };

  const close = () => {
    setOpen(false);
    setText(null);
  };

  const choose = (opt) => {
    onChange(opt);
    close();
  };

  const onKeyDown = (e) => {
    if (e.key === "ArrowDown" || e.key === "ArrowUp") {
      e.preventDefault();
      if (!open) return openList();
      const step = e.key === "ArrowDown" ? 1 : -1;
      setActive((a) => Math.min(Math.max(a + step, 0), shown.length - 1));
    } else if (e.key === "Enter" && open) {
      e.preventDefault(); // não envia o formulário de Configurações
      if (shown[active]) choose(shown[active]);
    } else if (e.key === "Escape" && open) {
      e.preventDefault();
      close();
    }
  };

  return (
    <div className="relative w-full md:max-w-[360px]">
      <input
        {...rest}
        id={id}
        type="text"
        role="combobox"
        aria-expanded={open}
        aria-controls={listId}
        aria-autocomplete="list"
        aria-activedescendant={
          open && shown[active] ? optionId(active) : undefined
        }
        autoComplete="off"
        spellCheck={false}
        value={text ?? value}
        onFocus={(e) => {
          e.target.select();
          openList();
        }}
        onClick={openList}
        onChange={(e) => {
          setText(e.target.value);
          setActive(0);
          setOpen(true);
        }}
        onKeyDown={onKeyDown}
        onBlur={close}
        className="w-full px-3 py-2.5 rounded-pg-control font-pg-mono text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-neblina focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
      />
      {open && (
        <ul
          id={listId}
          role="listbox"
          className="absolute z-10 left-0 right-0 mt-1 max-h-60 overflow-auto m-0 p-1 list-none rounded-pg-control bg-pg-branco border border-pg-borda"
        >
          {shown.length === 0 && (
            <li role="status" className="px-3 py-2 text-pg-ui text-pg-neblina">
              {emptyText}
            </li>
          )}
          {shown.map((opt, i) => (
            <li
              key={opt}
              id={optionId(i)}
              role="option"
              aria-selected={opt === value}
              // mouseDown (e não click) com preventDefault: o input não perde o foco antes da escolha
              onMouseDown={(e) => {
                e.preventDefault();
                choose(opt);
              }}
              onMouseEnter={() => setActive(i)}
              className={`px-3 py-2 rounded font-pg-mono text-pg-ui text-pg-ardosia cursor-pointer ${
                i === active ? "bg-pg-ciano-nevoa" : ""
              } ${opt === value ? "font-semibold" : ""}`}
            >
              {opt}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
