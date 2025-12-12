import { api } from "./client";
import type { User } from "../auth/auth";

export interface LoginResp {
  token: string;
  user: User;
}

export interface RegisterResp {
  user: User;
}

export async function login(username: string, password: string) {
  return api<LoginResp>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ username, password }),
  });
}

export async function register(username: string, password: string) {
  return api<RegisterResp>("/auth/register", {
    method: "POST",
    body: JSON.stringify({ username, password }),
  });
}
