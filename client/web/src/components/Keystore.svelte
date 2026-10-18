<script lang="ts">
  import { PaneGroup, Pane, PaneResizer } from "paneforge";
  import { MediaQuery } from "svelte/reactivity";
  const narrow = new MediaQuery("(max-width: 767px)");
  import { Button } from "$lib/components/ui/button";
  import api from "@/api";
  import CreateEntryModal from "@/components/CreateEntryModal.svelte";
  import CreateInvitationModal from "@/components/CreateInvitationModal.svelte";
  import DeleteKeystoreModal from "@/components/DeleteKeystoreModal.svelte";
  import Entry from "@/components/Entry.svelte";
  import EntryPlaceholder from "@/components/EntryPlaceholder.svelte";
  import Link from "@/components/Link.svelte";
  import { getInvitations } from "@/stores/invitations";
  import { getKeystores } from "@/stores/keystores";
  import { showMessage } from "@/stores/messages";
  import { navigate, useFocus, useParams } from "@/lib/router";

  let { keystore }: { keystore: import("@/api").Keystore } = $props();

  let showCreateInvitationModal: boolean = $state(false);
  let showCreateEntryModal: boolean = $state(false);
  let showDeleteKeystoreModal: boolean = $state(false);

  const params = useParams();
  const registerFocus = useFocus();

  let entry = $derived(keystore?.entries.find((e) => e.id === $params.entryId));

  function onInvitationCreated(e: { id: string }) {
    showCreateInvitationModal = false;
    getInvitations();
  }

  function onEntryCreated(e: { id: string }) {
    showCreateEntryModal = false;
    getKeystores();
  }

  function deleteKeystore() {
    api.deleteKeystore({ keystoreId: keystore.id }).then(async () => {
      await navigate("/", { replace: true });
      showMessage("Keystore deleted");
      showDeleteKeystoreModal = false;
      await getKeystores();
    });
  }
</script>

<PaneGroup
  direction={narrow.current ? "vertical" : "horizontal"}
  autoSaveId="myst-entries-layout"
  class="min-h-0 min-w-0 flex-1"
>
  <Pane defaultSize={40} minSize={20}>
    <div class="entries-list">
      <div class="entries" use:registerFocus>
        {#if !keystore.readOnly}
          <div class="entries-list-header">
            <Button onclick={() => (showCreateEntryModal = true)} class="button"
              >Create Entry</Button
            >
            <Button
              onclick={() => (showCreateInvitationModal = true)}
              variant="secondary">Share Keystore</Button
            >
            <Button
              onclick={() => (showDeleteKeystoreModal = true)}
              variant="destructive">Delete Keystore</Button
            >
          </div>
        {/if}
        <h1 class="mb-1 break-all text-lg font-semibold">{keystore.name}</h1>
        <p class="mb-5 text-xs text-muted-foreground">
          {keystore.entries.length} entries{keystore.readOnly
            ? " · Read only"
            : ""}
        </p>
        {#each keystore.entries as entry (entry.id)}
          <Link path={`/keystore/${keystore.id}/entry/${entry.id}`}>
            <div class="entry-row">
              <span class="icon">
                <span class="placeholder"></span>
              </span>
              <div class="info">
                <span class="name">
                  {entry.website}
                </span>
                <span class="user">
                  {entry.username}
                </span>
              </div>
            </div>
          </Link>
        {:else}
          <p class="py-8 text-center text-sm text-muted-foreground">
            No entries yet.
          </p>
        {/each}
      </div>
    </div>
  </Pane>
  <PaneResizer
    class="z-10 bg-border hover:bg-primary focus-visible:bg-primary data-[direction=horizontal]:w-1 data-[direction=vertical]:h-1"
  />
  <Pane minSize={30}>
    {#if entry}
      <Entry {entry} {keystore} />
    {:else}
      <EntryPlaceholder />
    {/if}
  </Pane>
</PaneGroup>

<DeleteKeystoreModal
  bind:show={showDeleteKeystoreModal}
  onsubmit={(e) => {
    deleteKeystore();
  }}
/>

<CreateInvitationModal
  bind:show={showCreateInvitationModal}
  {keystore}
  oncreated={(e) => {
    onInvitationCreated(e);
  }}
/>

<CreateEntryModal
  bind:show={showCreateEntryModal}
  {keystore}
  oncreated={(e) => {
    onEntryCreated(e);
  }}
/>

<style>
  .entries-list {
    flex: 0 0 50%;
    min-width: 0;
    height: 100%;
  }
  .entries {
    height: 100%;
    overflow-y: auto;
    padding: 24px;
  }
  .entries :global(.link a) {
    display: block;
    border-radius: 10px;
    color: inherit;
    text-decoration: none;
  }
  .entries :global(.link a:hover) {
    background: var(--muted);
  }
  .entries :global(.link.active a) {
    background: var(--accent);
    color: var(--accent-foreground);
  }
  .entry-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
  }
  .info {
    display: grid;
    gap: 3px;
    min-width: 0;
  }
  .name,
  .user {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .name {
    font-weight: 500;
  }
  .user {
    font-size: 12px;
    color: var(--muted-foreground);
  }
  .placeholder {
    display: block;
    width: 32px;
    height: 32px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--secondary);
  }
</style>
