<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import type { Invitation, Keystore } from "@/api";
  import api from "@/api";
  import Field from "@/components/Field.svelte";
  import InputField from "@/components/InputField.svelte";
  import Modal from "@/components/Modal.svelte";
  import { showError, showMessage } from "@/stores/messages";
  import { onMount } from "svelte";

  interface Props {
    oncreated?: (value: { id: string }) => void;
    show?: boolean;
    keystore: Keystore;
  }

  let { oncreated, show = $bindable(false), keystore }: Props = $props();

  let user = $state("");

  let userValid = $derived(user.trim() !== "");

  let allowSubmit = $derived(userValid);

  $effect(() => {
    if (!show) {
      user = "";
    }
  });

  onMount(async () => {
    user = "";
  });

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    api
      .createInvitation({
        keystoreId: keystore.id,
        requestBody: {
          invitee: user,
        },
      })
      .then((inv: Invitation) => {
        showMessage("Invitation sent");
        oncreated?.({ id: inv.id });
      })
      .catch((err) => {
        showError("Create Invitation Failed");
        console.error(err);
      });
  };
</script>

<form
  class="create-invitation-modal"
  onsubmit={(event) => {
    event.preventDefault();
    submit();
  }}
>
  <Modal bind:show>
    {#snippet header()}
      <div class="modal-header">
        <div class="title">
          <span>Share with...</span>
        </div>
      </div>
    {/snippet}

    <div class="modal-content">
      <Field label="Keystore" value={keystore.name} />
      <InputField bind:value={user} label="User" />
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
          >Invite User</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
