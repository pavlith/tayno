<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import type { Entry } from "@/api";
  import Modal from "@/components/Modal.svelte";
  const submit = () => {
    onsubmit?.();
  };

  interface Props {
    onsubmit?: (value: void) => void;
    show?: boolean;
    entry: Entry;
  }

  let { onsubmit, show = $bindable(false), entry }: Props = $props();
</script>

<form
  onsubmit={(event) => {
    event.preventDefault();
    submit();
  }}
>
  <Modal bind:show>
    {#snippet header()}
      <div class="delete-title">
        Are you sure you want to delete this entry?
      </div>
    {/snippet}
    {#snippet footer()}
      <div class="modal-footer">
        <Button
          variant="ghost"
          class="button transparent"
          onclick={() => (show = false)}
          type="button">Cancel</Button
        >
        <Button variant="destructive" class="button red" type="submit"
          >Delete</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
