import axios from "axios";

export function createClient() {
  const apiKey = localStorage.getItem("pingou_api_key") || "";

  const instance = axios.create({
    baseURL: "/api",
    headers: {
      "X-API-Key": apiKey,
      "Content-Type": "application/json",
    },
  });
  // mensagem de erro da API (ex.: validação) no lugar do texto genérico do axios
  instance.interceptors.response.use(
    (res) => res,
    (err) => Promise.reject(new Error(err.response?.data?.error || err.message)),
  );
  return instance;
}

export const client = {
  get: (...args) => createClient().get(...args),
  post: (...args) => createClient().post(...args),
  patch: (...args) => createClient().patch(...args),
  delete: (...args) => createClient().delete(...args),
};
