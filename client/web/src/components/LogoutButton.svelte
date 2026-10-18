<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { logout } from "@/stores/authState";
  let pending = $state(false);
  let error = $state("");
  async function submit() {
    if (pending) return;
    pending = true;
    error = "";
    try { await logout(); }
    catch { error = "Unable to log out. Please try again."; }
    finally { pending = false; }
  }
</script>

<Button type="button" variant="outline" disabled={pending} onclick={submit}>
  {pending ? "Logging out…" : "Logout"}
</Button>
{#if error}<p role="alert" class="text-sm text-destructive">{error}</p>{/if}
