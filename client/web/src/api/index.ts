import { DefaultService, OpenAPI } from "@/api/generated";
import type { Error as ApiError } from "@/api/generated";

OpenAPI.BASE = (import.meta.env.VITE_API_URL?.trim() || "/api").replace(/\/+$/, "");

// The generated client includes error bodies in its success unions. Reject these
// centrally so callers receive only successful response data.
type Service = {
  [K in keyof typeof DefaultService]: (typeof DefaultService)[K] extends (
    ...args: infer A
  ) => PromiseLike<infer R>
    ? (...args: A) => Promise<Exclude<R, ApiError>>
    : (typeof DefaultService)[K];
};
const api = new Proxy(DefaultService, {
  get(target, key) {
    const method = Reflect.get(target, key);
    if (typeof method !== "function") return method;
    return async (...args: unknown[]) => {
      const result: unknown = await Reflect.apply(method, target, args);
      if (
        result &&
        typeof result === "object" &&
        "code" in result &&
        "message" in result
      ) {
        throw new globalThis.Error(String(result.message));
      }
      return result;
    };
  },
}) as Service;
export default api;
export * from "@/api/generated";
