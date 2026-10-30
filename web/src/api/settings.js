// src/api/settings.js
import { client } from "./client";

export const settingsApi = {
  get: () => client.get("/settings").then((r) => r.data),
  update: (data) => client.patch("/settings", data).then((r) => r.data),
  info: () => client.get("/info").then((r) => r.data),
  // url vazia → a API testa a URL salva
  testWebhook: (url) =>
    client
      .post("/settings/webhook/test", url ? { url } : {})
      .then((r) => r.data),
  export: () =>
    client.get("/export", { responseType: "blob" }).then((r) => r.data),
};
