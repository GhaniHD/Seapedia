import type { ApiErrorBody } from "./types";

// Backend (Gin) always wraps success payloads as { data: ... } or { message: ... }
// and errors as { error: ... }. This envelope is unwrapped here once so the rest
// of the app can work directly with typed DTOs.
interface ApiEnvelope<T> {
  data?: T;
  message?: string;
}

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
const TOKEN_KEY = "seapedia_token";

export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string | null) {
  if (typeof window === "undefined") return;
  if (token) window.localStorage.setItem(TOKEN_KEY, token);
  else window.localStorage.removeItem(TOKEN_KEY);
}

interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  body?: unknown;
  auth?: boolean; // send Authorization header (default true)
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, auth = true } = opts;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (auth) {
    const token = getToken();
    if (token) headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
    cache: "no-store",
  });

  if (res.status === 204) {
    return undefined as T;
  }

  const text = await res.text();
  let data: unknown = undefined;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      // non-JSON response
    }
  }

  if (!res.ok) {
    const body = data as ApiErrorBody | undefined;
    const message = body?.error || body?.message || res.statusText || "Terjadi kesalahan";
    throw new ApiError(message, res.status);
  }

  const envelope = data as ApiEnvelope<T> | undefined;
  if (envelope && typeof envelope === "object" && "data" in envelope) {
    return envelope.data as T;
  }
  return data as T;
}

export const api = {
  get: <T>(path: string, auth = true) => request<T>(path, { method: "GET", auth }),
  post: <T>(path: string, body?: unknown, auth = true) =>
    request<T>(path, { method: "POST", body, auth }),
  put: <T>(path: string, body?: unknown, auth = true) =>
    request<T>(path, { method: "PUT", body, auth }),
  del: <T>(path: string, auth = true) => request<T>(path, { method: "DELETE", auth }),
};
