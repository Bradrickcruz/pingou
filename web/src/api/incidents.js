// src/api/incidents.js
import { client } from "./client";

export const incidentsApi = {
  list: (params) => client.get("/incidents", { params }).then((r) => r.data),
};
