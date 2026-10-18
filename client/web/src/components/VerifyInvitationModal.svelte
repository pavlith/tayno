<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import type { Invitation } from "@/api";
  import Modal from "@/components/Modal.svelte";
  import { hash } from "@/lib/color-hash.js";

  interface Props {
    show?: boolean;
    invitation: Invitation;
  }

  let { show = $bindable(false), invitation }: Props = $props();
</script>

<form onsubmit={(event) => event.preventDefault()}>
  <Modal bind:show>
    <div class="modal-content">
      Make sure the identity icons on both yours and the inviter's side are the
      same.

      <div class="identity">
        <strong style="color: {hash(invitation.inviter.username)}"
          >{invitation.inviter.username}</strong
        >'s identity:
        <div>
          {#if invitation.inviter.icon}
            <img
              style="width: 64px; height: 64px;"
              src={"data:image/svg+xml," +
                encodeURIComponent(invitation.inviter.icon)}
              alt=""
            />
          {/if}
        </div>
      </div>
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
          variant="default"
          class="button green"
          onclick={() => (show = false)}
          type="button">Verify</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
