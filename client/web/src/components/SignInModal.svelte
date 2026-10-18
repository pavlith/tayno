<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import api from "@/api";
  import InputField from "@/components/InputField.svelte";
  import Modal from "@/components/Modal.svelte";
  import { showError, showMessage } from "@/stores/messages";

  interface Props {
    // import {getCurrentUser} from "@/stores/user";
    show?: boolean;
  }

  let { show = $bindable(false) }: Props = $props();

  let username = $state(""),
    password = $state("");

  let usernameValid = $derived(username.trim() !== "");
  let passwordValid = $derived(password.trim() !== "");

  let allowSubmit = $derived(usernameValid && passwordValid);

  $effect(() => {
    if (!show) {
      username = "";
      password = "";
    }
  });

  const reset = () => {
    username = "";
    password = "";
  };

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    api
      .loginDeprecated({
        requestBody: {
          username,
          password,
        },
      })
      .then(() => {
        showMessage("Signed in.");
        reset();
        // getCurrentUser()
        show = false;
      })
      .catch((err) => {
        showError("Signing in failed.");
        console.error(err);
      });
  };
</script>

<form
  class="create-entry-modal"
  onsubmit={(event) => {
    event.preventDefault();
    submit();
  }}
>
  <Modal bind:show>
    {#snippet header()}
      <div class="create-title">Sign in</div>
    {/snippet}

    <div class="modal-content">
      <InputField bind:value={username} label="Username" />
      <InputField bind:value={password} label="Password" />
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
          >Sign In</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
