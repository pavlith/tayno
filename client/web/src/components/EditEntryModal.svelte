<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import type { Entry } from "@/api";
  import Field from "@/components/Field.svelte";
  import InputField from "@/components/InputField.svelte";
  import Modal from "@/components/Modal.svelte";
  import { onMount } from "svelte";
  import { showMessage } from "@/stores/messages";

  interface Props {
    onsubmit?: (value: { password: string; notes: string }) => void;
    show?: boolean;
    entry: Entry;
  }

  let { onsubmit, show = $bindable(false), entry }: Props = $props();

  let password = $state("");

  let notes = $state("");

  let passwordValid = $derived(password.trim() !== "");
  let passwordChanged = $derived(entry && password.trim() !== entry.password);
  let notesChanged = $derived(entry && notes.trim() !== entry.notes);

  const allowSubmit = $derived(
    passwordValid && (passwordChanged || notesChanged),
  );

  $effect(() => {
    if (!show && entry) {
      password = entry.password;
      notes = entry.notes;
    }
  });

  onMount(async () => {
    password = entry.password;
    notes = entry.notes;
  });

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    onsubmit?.({
      password,
      notes,
    });
  };
</script>

<form
  class="edit-entry-modal"
  onsubmit={(event) => {
    event.preventDefault();
    submit();
  }}
>
  <Modal bind:show>
    {#snippet header()}
      <div class="modal-header">
        <div class="title">
          <span class="website">{entry.website}</span>
          <span class="username">{entry.username}</span>
        </div>
      </div>
    {/snippet}

    <div class="modal-content">
      <Field label="Website" value={entry.website} />
      <Field label="Username" value={entry.username} />
      <InputField
        bind:value={password}
        error={!passwordValid && "Password cannot be empty"}
        label="Password"
      />
      <InputField bind:value={notes} label="Notes" />
    </div>

    {#snippet footer()}
      <div class="modal-footer">
        <Button
          variant="ghost"
          class="button transparent"
          onclick={() => (show = false)}
          type="button">Cancel</Button
        >
        <Button
          disabled={!allowSubmit}
          variant="default"
          class="button green"
          type="submit">Save</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
