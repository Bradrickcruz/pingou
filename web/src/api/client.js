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
  // mensagem de erro da API (ex.: validação) no lugar do texto genérico do axios;
  // status fica no erro para a tela tratar casos como 429
  instance.interceptors.response.use(
    (res) => res,
    (err) =>
      Promise.reject(
        Object.assign(new Error(err.response?.data?.error || err.message), {
          status: err.response?.status,
        }),
      ),
  );
  return instance;
}

export const client = {
  get: (...args) => createClient().get(...args),
  post: (...args) => createClient().post(...args),
  patch: (...args) => createClient().patch(...args),
  delete: (...args) => createClient().delete(...args),
};
