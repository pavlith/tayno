<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import Modal from "@/components/Modal.svelte";
  const submit = () => {
    onsubmit?.();
  };

  interface Props {
    onsubmit?: (value: void) => void;
    show?: boolean;
  }

  let { onsubmit, show = $bindable(false) }: Props = $props();
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
        Are you sure you want to delete this keystore?
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
