// Períodos da página de detalhe do monitor (mesmas chaves do /stats).
export const RANGES = [
  { key: "24h", label: "24h", short: "24 h", days: 1, since: "nas últimas 24 h" },
  { key: "7d", label: "7 dias", short: "7 d", days: 7, since: "nos últimos 7 dias" },
  { key: "14d", label: "14 dias", short: "14 d", days: 14, since: "nos últimos 14 dias" },
  { key: "30d", label: "30 dias", short: "30 d", days: 30, since: "nos últimos 30 dias" },
];

// "24h" e "7 dias" valem sempre: a retenção mínima é 7 dias.
export const MIN_RETENTION = 7;

// clampRange: maior período permitido pela retenção, sem passar do pedido.
// Chave desconhecida vira "24h".
export function clampRange(key, retentionDays = MIN_RETENTION) {
  const asked = RANGES.findIndex((r) => r.key === key);
  if (asked < 0) return "24h";
  for (let i = asked; i >= 0; i--) {
    if (RANGES[i].days <= retentionDays) return RANGES[i].key;
  }
  return "24h";
}
