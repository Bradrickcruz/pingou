import axios from "axios";

export function createClient() {
  const apiKey = sessionStorage.getItem("pingou_api_key") || "";

  return axios.create({
    baseURL: "/api",
    headers: {
      "X-API-Key": apiKey,
      "Content-Type": "application/json",
    },
  });
}

export const client = {
  get: (...args) => createClient().get(...args),
  post: (...args) => createClient().post(...args),
  patch: (...args) => createClient().patch(...args),
  delete: (...args) => createClient().delete(...args),
};

// interceptor global de erro
const base = axios.create();
base.interceptors.response.use(
  (res) => res,
  (err) => {
    const msg = err.response?.data?.error || err.message;
    return Promise.reject(new Error(msg));
  },
);
