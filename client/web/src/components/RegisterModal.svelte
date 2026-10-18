<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import api from "@/api";
  import InputField from "@/components/InputField.svelte";
  import Modal from "@/components/Modal.svelte";
  import { showError, showMessage } from "@/stores/messages";
  import { getCurrentUser } from "@/stores/user";
  import PasswordInputField from "@/components/PasswordInputField.svelte";
  import { getInvitations } from "@/stores/invitations";

  interface Props {
    // import {getCurrentUser} from "@/stores/user";
    show?: boolean;
  }

  let { show = $bindable(false) }: Props = $props();

  let username: string = $state("");
  let password: string = $state("");

  let usernameChanged = $state(false);
  let passwordChanged = $state(false);

  $effect(() => {
    username;

    if (username !== "") {
      usernameChanged = true;
    }
  });

  $effect(() => {
    password;

    if (password !== "") {
      passwordChanged = true;
    }
  });

  let showErrors = $derived(usernameChanged && passwordChanged);

  let usernameEmpty = $derived(username.trim() === "");
  let passwordEmpty = $derived(password.trim() === "");

  let usernameValid = $derived(!usernameEmpty);
  let passwordValid = $derived(!passwordEmpty);

  let usernameError = $derived(
    usernameEmpty ? "Username cannot be empty." : "",
  );
  let passwordError = $derived(
    passwordEmpty ? "Password cannot be empty." : "",
  );

  let allowSubmit = $derived(!usernameEmpty && !passwordEmpty);

  const reset = () => {
    username = "";
    password = "";
  };

  $effect(() => {
    show;

    if (!show) {
      reset();
    }
  });

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    api
      .register({
        requestBody: {
          username,
          password,
        },
      })
      .then(async () => {
        showMessage("Signed in.");
        await getCurrentUser();
        await getInvitations();
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
      <div class="create-title">Register</div>
    {/snippet}

    <div class="modal-content">
      <InputField
        bind:value={username}
        label="Username"
        name="text"
        error={!usernameValid && showErrors && usernameError}
      />
      <PasswordInputField
        bind:value={password}
        label="Password"
        name="telephone"
        error={!passwordValid && showErrors && passwordError}
      />
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
          type="submit">Register</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
