import { getToken } from "../auth/auth";
const API_BASE = import.meta.env.VITE_API_BASE as string;

export interface ApiError {
  code: string;
  message: string;
}

export async function api<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();

  // Use Headers to make setting additional keys type-safe
  const headers = new Headers(options.headers || {});
  headers.set("Content-Type", "application/json");
  if (token) headers.set("Authorization", `Bearer ${token}`);

  const resp = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  const body = await resp.json();

  if (!resp.ok) {
    throw (
      (body.error as ApiError) ?? {
        code: "UNKNOWN",
        message: "unknown error",
      }
    );
  }

  return body.data as T;
}
