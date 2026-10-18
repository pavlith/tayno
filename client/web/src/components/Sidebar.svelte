<script lang="ts">
  import LogoutButton from "./LogoutButton.svelte";
  import Link from "./Link.svelte";
  import { invitations } from "@/stores/invitations";
  import { currentUser } from "@/stores/user";
  import { Button } from "$lib/components/ui/button";
  import { Badge } from "$lib/components/ui/badge";
  import { Separator } from "$lib/components/ui/separator";
  import Plus from "@lucide/svelte/icons/plus";
  import ShieldCheck from "@lucide/svelte/icons/shield-check";
  import type { Keystore } from "@/api";
  let {
    keystores,
    showCreateKeystoreModal = $bindable(false),
    showRegisterModal = $bindable(false),
  }: {
    keystores: Keystore[];
    showCreateKeystoreModal?: boolean;
    showRegisterModal?: boolean;
  } = $props();
  const newInvitationsCount = $derived(
    $invitations.filter(
      (inv) => inv.status === "pending" || inv.status === "accepted",
    ).length,
  );
</script>

<aside class="sidebar flex flex-col gap-6 p-5" aria-label="Main navigation">
  <div class="flex items-center gap-3 px-2 py-2 text-xl font-semibold">
    <ShieldCheck class="size-6 text-primary" />Myst
  </div>
  <Separator />
  <div>
    <div class="mb-3 flex items-center justify-between px-2">
      <h2
        class="text-xs font-medium uppercase tracking-wider text-muted-foreground"
      >
        Keystores
      </h2>
      <Button
        variant="ghost"
        size="icon-sm"
        aria-label="Create keystore"
        onclick={() => (showCreateKeystoreModal = true)}><Plus /></Button
      >
    </div>
    <nav class="grid gap-1">
      {#each keystores as keystore (keystore.id)}
        <Link path={`/keystore/${keystore.id}`}
          ><span class="truncate">{keystore.name}</span></Link
        >
      {/each}
    </nav>
  </div>
  <div class="mt-auto grid gap-4">
    {#if $currentUser}
      <Link path="/invitations"
        >Invitations {#if newInvitationsCount}<Badge class="ml-auto"
            >{newInvitationsCount}</Badge
          >{/if}</Link
      >
      <p class="break-all px-2 text-xs text-muted-foreground">
        Signed in as <span class="text-foreground">{$currentUser.username}</span
        >
      </p>
    {:else if $currentUser === null}
      <Button variant="outline" onclick={() => (showRegisterModal = true)}
        >Register for sync</Button
      >
    {/if}
    <LogoutButton />
    <p class="flex items-center gap-2 px-2 text-xs text-muted-foreground">
      <ShieldCheck class="size-3.5" />Encrypted keystores
    </p>
  </div>
</aside>

<style>
  nav :global(.link a),
  aside :global(.link a) {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    border-radius: 8px;
    color: var(--muted-foreground);
    text-decoration: none;
  }
  aside :global(.link a:hover) {
    background: var(--muted);
    color: var(--foreground);
  }
  aside :global(.link.active a) {
    background: var(--accent);
    color: var(--accent-foreground);
    font-weight: 500;
  }
</style>
