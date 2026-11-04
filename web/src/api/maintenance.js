import { client } from "./client";

export const maintenanceApi = {
  // { monitor_ids, starts_at, ends_at, reason } → { windows }
  create: (data) =>
    client.post("/maintenance-windows", data).then((r) => r.data),
  cancel: (id) =>
    client.delete(`/maintenance-windows/${id}`).then((r) => r.data),
};
