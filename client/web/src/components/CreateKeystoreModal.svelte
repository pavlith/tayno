<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  import type { Keystore } from "@/api";
  import api from "@/api";
  import InputField from "@/components/InputField.svelte";
  import Modal from "@/components/Modal.svelte";
  import { showError, showMessage } from "@/stores/messages";
  import { keystores } from "@/stores/keystores";

  interface Props {
    oncreated?: (value: { id: string }) => void;
    show?: boolean;
  }

  let { oncreated, show = $bindable(false) }: Props = $props();

  let name = $state("");

  let showErrors = $state(false);

  $effect(() => {
    name;

    if (name !== "") {
      showErrors = true;
    }
  });

  const reset = () => {
    name = "";
    showErrors = false;
  };

  $effect(() => {
    show;

    if (!show) {
      reset();
    }
  });

  let nameEmpty = $derived(name.trim() === "");
  let nameTooLong = $derived(name.trim().length > 24);
  let nameAlreadyExists = $derived(
    $keystores.find((k) => k.name === name.trim()) !== undefined,
  );

  let nameValid = $derived(!nameEmpty && !nameTooLong && !nameAlreadyExists);
  let nameError = $derived(
    nameEmpty
      ? "Name cannot be empty"
      : nameTooLong
        ? "Name cannot be longer than 24 characters"
        : nameAlreadyExists
          ? "Name already exists"
          : "",
  );

  let allowSubmit = $derived(nameValid);

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    api
      .createKeystore({
        requestBody: {
          name,
        },
      })
      .then((k: Keystore) => {
        showMessage("Keystore created");
        oncreated?.({ id: k.id });
        show = false;
      })
      .catch((err) => {
        showError("Create Keystore Failed");
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
      <div class="create-title">Create Keystore</div>
    {/snippet}

    <div class="modal-content">
      <InputField
        bind:value={name}
        label="Name"
        error={!nameValid && showErrors && nameError}
      />
    </div>

    {#snippet footer()}
      <div class="modal-footer">
        <Button
          variant="ghost"
          class="button transparent"
          onclick={() => {
            show = false;
          }}
          type="button"
          >Cancel
        </Button>
        <Button
          disabled={!allowSubmit}
          variant="default"
          class="button green"
          type="submit">Create</Button
        >
      </div>
    {/snippet}
  </Modal>
</form>
