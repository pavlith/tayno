<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import type { Entry, Keystore } from "@/api";
  import api from "@/api";
  import DeleteEntryModal from "@/components/DeleteEntryModal.svelte";
  import EditEntryModal from "@/components/EditEntryModal.svelte";
  import Field from "@/components/Field.svelte";
  import { getKeystores } from "@/stores/keystores";

  let showDeleteModal = $state(false);
  let showEditModal = $state(false);

  const updateEntry = async (password: string, notes: string) => {
    api
      .updateEntry({
        keystoreId: keystore.id,
        entryId: entry.id,
        requestBody: {
          password: password != entry.password ? password : undefined,
          notes: notes != entry.notes ? notes : undefined,
        },
      })
      .then(async () => {
        showEditModal = false;
        await getKeystores();
      });
  };

  interface Props {
    keystore: Keystore;
    entry: Entry;
  }

  let { keystore, entry }: Props = $props();

  const deleteEntry = async () => {
    await api
      .deleteEntry({
        keystoreId: keystore.id,
        entryId: entry.id,
      })
      .then(() => {
        showDeleteModal = false;
      });

    await getKeystores();
  };
</script>

<div class="entry" class:empty={!entry}>
  {#if entry}
    <div class="header">
      <div class="title">
        <h5 class="break-all text-xl font-semibold">{entry.website}</h5>
        <span class="text-sm text-muted-foreground">{entry.username}</span>
      </div>
      {#if !keystore.readOnly}
        <div class="actions">
          <Button onclick={() => (showEditModal = true)} variant="secondary">
            Edit</Button
          >
          <Button
            onclick={() => (showDeleteModal = true)}
            variant="destructive"
          >
            Delete</Button
          >
        </div>
      {/if}
    </div>
    <div class="fields">
      <Field label="Website" value={entry.website} />
      <Field label="Username" value={entry.username} />
      <Field label="Password" value={entry.password} />
      <Field label="Notes" value={entry.notes} />
    </div>
  {/if}
</div>

<EditEntryModal
  bind:show={showEditModal}
  {entry}
  onsubmit={(e) => {
    updateEntry(e.password, e.notes);
  }}
/>

<DeleteEntryModal
  bind:show={showDeleteModal}
  {entry}
  onsubmit={() => {
    deleteEntry();
  }}
/>
