<script lang="ts">
  import type { Snippet } from "svelte";
  import { useLocation, navigate } from "@/lib/router";
  let { path = "", children }: { path?: string; children?: Snippet } = $props();
  const location = useLocation();
  const active = $derived(
    $location.pathname === path || $location.pathname.startsWith(path + "/"),
  );
</script>

<div class="link" class:active style="max-width: 100%">
  <a
    href={path}
    aria-current={active ? "page" : undefined}
    onclick={(event) => {
      if (
        event.button === 0 &&
        !event.metaKey &&
        !event.ctrlKey &&
        !event.shiftKey &&
        !event.altKey
      ) {
        event.preventDefault();
        navigate(path);
      }
    }}>{@render children?.()}</a
  >
</div>
