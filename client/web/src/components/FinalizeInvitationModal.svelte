<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import type { Invitation } from "@/api";
  import Modal from "@/components/Modal.svelte";
  import { hash } from "@/lib/color-hash.js";
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
        Are you sure you want to verify this invitation?
      </div>
    {/snippet}

    <div class="modal-content">
      Make sure the invitee's identity icon below matches the one on their side,
      and then proceed to verify.

      <div class="identity">
        <strong style="color: {hash(invitation.invitee.username)}"
          >{invitation.invitee.username}</strong
        >'s identity:
        <div>
          {#if invitation.inviter.icon}
            <img
              style="width: 64px; height: 64px;"
              src={"data:image/svg+xml," +
                encodeURIComponent(invitation.inviter.icon)}
              alt=""
            />
          {:else if invitation.invitee.icon}
            <img
              style="width: 64px; height: 64px;"
              src={"data:image/svg+xml," +
                encodeURIComponent(invitation.invitee.icon)}
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
        <Button variant="default" class="button green" type="submit"
          >Verify</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
