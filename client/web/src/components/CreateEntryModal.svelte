<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import api from "@/api";
  import type { Entry, Invitation, Keystore } from "@/api";
  import InputField from "@/components/InputField.svelte";
  import Modal from "@/components/Modal.svelte";
  import { showError, showMessage } from "@/stores/messages";

  interface Props {
    oncreated?: (value: { id: string }) => void;
    keystore: Keystore;
    show?: boolean;
  }

  let { oncreated, keystore, show = $bindable(false) }: Props = $props();

  let website = $state("");

  let username = $state("");

  let password = $state("");

  let notes = $state("");

  let showErrors = $state(false);

  $effect(() => {
    website;
    if (website !== "") {
      showErrors = true;
    }

    username;
    if (username !== "") {
      showErrors = true;
    }

    password;
    if (password !== "") {
      showErrors = true;
    }
  });

  let websiteValid = $derived(website.trim() !== "");
  let usernameValid = $derived(username.trim() !== "");
  let passwordValid = $derived(password.trim() !== "");

  let allowSubmit = $derived(passwordValid && websiteValid && usernameValid);

  $effect(() => {
    if (!show) {
      website = "";
      username = "";
      password = "";
      notes = "";
    }
  });

  const reset = () => {
    website = "";
    username = "";
    password = "";
    notes = "";
    showErrors = false;
  };

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    api
      .createEntry({
        keystoreId: keystore.id,
        requestBody: {
          website,
          username,
          password,
          notes,
        },
      })
      .then((e: Entry) => {
        showMessage("Entry created");
        oncreated?.({ id: e.id });
        reset();
      })
      .catch((err) => {
        showError("Create Entry Failed");
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
      <div class="create-title">Create entry</div>
    {/snippet}

    <div class="modal-content">
      <InputField
        bind:value={website}
        label="Website"
        error={!websiteValid && showErrors && "Website cannot be empty"}
      />
      <InputField
        bind:value={username}
        label="Username"
        error={!usernameValid && showErrors && "Username cannot be empty"}
      />
      <InputField
        bind:value={password}
        error={!passwordValid && showErrors && "Password cannot be empty"}
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
        <Button variant="default" class="button green" type="submit"
          >Create Entry</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
