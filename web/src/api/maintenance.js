import { client } from "./client";

export const maintenanceApi = {
  // → { windows: [{ id, monitor_id, monitor_name, starts_at, ends_at, reason, state }] }
  list: () => client.get("/maintenance-windows").then((r) => r.data),
  // patch: { starts_at?, ends_at?, reason? } (só os campos que mudaram)
  update: (id, patch) =>
    client.patch(`/maintenance-windows/${id}`, patch).then((r) => r.data),
  // { monitor_ids, starts_at, ends_at, reason } → { windows }
  create: (data) =>
    client.post("/maintenance-windows", data).then((r) => r.data),
  cancel: (id) =>
    client.delete(`/maintenance-windows/${id}`).then((r) => r.data),
};
