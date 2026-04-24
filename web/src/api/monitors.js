import { client } from "./client";

export const monitorsApi = {
  list: (params) => client.get("/monitors", { params }).then((r) => r.data),
  get: (id) => client.get(`/monitors/${id}`).then((r) => r.data),
  create: (data) => client.post("/monitors", data).then((r) => r.data),
  update: (id, data) =>
    client.patch(`/monitors/${id}`, data).then((r) => r.data),
  delete: (id) => client.delete(`/monitors/${id}`).then((r) => r.data),
  checks: (id, params) =>
    client.get(`/monitors/${id}/checks`, { params }).then((r) => r.data),
  incidents: (id, params) =>
    client.get(`/monitors/${id}/incidents`, { params }).then((r) => r.data),
};
