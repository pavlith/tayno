import { writable, derived } from "svelte/store";

const location = writable({ pathname: window.location.pathname });
window.addEventListener("popstate", () =>
  location.set({ pathname: window.location.pathname }),
);
export function navigate(path: string, options: { replace?: boolean } = {}) {
  if (window.location.pathname === path) return;
  window.history[options.replace ? "replaceState" : "pushState"]({}, "", path);
  location.set({ pathname: window.location.pathname });
}
export const useNavigate = () => navigate;
export const useLocation = () => location;
export const useParams = () =>
  derived(location, ({ pathname }) => {
    const match = pathname.match(
      /^\/keystore\/([^/]+)(?:\/entry\/([^/]+))?\/?$/,
    );
    return { keystoreId: match?.[1], entryId: match?.[2] };
  });
export const useFocus = () => (node: HTMLElement) => {
  node.tabIndex = -1;
  return {};
};
