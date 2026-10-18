import api, { OpenAPI } from "@/api";
import { request } from "@/api/generated/core/request";
import { writable } from "svelte/store";
import type { Readable } from "svelte/store";

export enum AuthState {
  SignedOut = "SIGNED_OUT",
  Locked = "LOCKED",
  SignedIn = "SIGNED_IN",
}

const state = writable<AuthState | undefined>(undefined);
export const authState = { subscribe: state.subscribe } as Readable<
  AuthState | undefined
>;
const accountKey = "myst-account-session";
let accountToken = sessionStorage.getItem(accountKey) ?? "";
let revision = 0;
let exiting = false;
OpenAPI.HEADERS = async (): Promise<Record<string, string>> =>
  accountToken ? { "X-Myst-Account": accountToken } : {};

export async function loginAccount(username: string, password: string) {
  const result = await request<{ token: string }>(OpenAPI, {
    method: "POST",
    url: "/auth/login",
    body: { username, password },
    mediaType: "application/json",
  });
  setAccount(result.token);
}

export async function signupAccount(
  username: string,
  password: string,
  masterPassword: string,
) {
  const result = await request<{ token: string }>(OpenAPI, {
    method: "POST",
    url: "/auth/signup",
    body: { username, password, masterPassword },
    mediaType: "application/json",
  });
  setAccount(result.token);
}

function setAccount(token: string) {
  revision++;
  accountToken = token;
  sessionStorage.setItem(accountKey, token);
  OpenAPI.TOKEN = undefined;
  state.set(AuthState.Locked);
}

export async function logout(): Promise<void> {
  if (exiting) return;
  exiting = true;
  revision++;
  try {
    await request<void>(OpenAPI, { method: "POST", url: "/logout" });
    accountToken = "";
    sessionStorage.removeItem(accountKey);
    OpenAPI.TOKEN = undefined;
    state.set(AuthState.SignedOut);
    window.location.replace("/");
  } finally {
    exiting = false;
  }
}

export async function getAuthState(): Promise<AuthState | undefined> {
  if (exiting) return undefined;
  const started = revision;
  let next: AuthState;
  if (!accountToken) {
    next = AuthState.SignedOut;
  } else {
    try {
      await request(OpenAPI, { method: "GET", url: "/account" });
    } catch (error: any) {
      if (error.status !== 401) throw error;
      if (started !== revision || exiting) return undefined;
      accountToken = "";
      sessionStorage.removeItem(accountKey);
      OpenAPI.TOKEN = undefined;
      state.set(AuthState.SignedOut);
      return AuthState.SignedOut;
    }
    try {
      await api.enclave();
      await api.currentUser();
      next = AuthState.SignedIn;
    } catch (error: any) {
      if (![401, 403, 404].includes(error.status)) throw error;
      next = AuthState.Locked;
    }
  }
  if (started !== revision || exiting) return undefined;
  state.set(next);
  return next;
}
