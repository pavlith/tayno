<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import type { Invitation } from "@/api";
  import Modal from "@/components/Modal.svelte";
  const submit = () => {
    onsubmit?.();
  };

  interface Props {
    onsubmit?: (value: void) => void;
    show?: boolean;
    invitation: Invitation;
  }

  let { onsubmit, show = $bindable(false), invitation }: Props = $props();
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
        Are you sure you want to accept this invitation?
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
        <Button variant="default" class="button green" type="submit"
          >Accept</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
