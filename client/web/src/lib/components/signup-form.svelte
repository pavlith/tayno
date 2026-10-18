<script lang="ts">
  import { signupAccount } from "@/stores/authState";
  import { Input } from "$lib/components/ui/input";
  import * as Field from "$lib/components/ui/field";
  const id = $props.id();

  import { Button } from "$lib/components/ui/button";
  import * as Card from "$lib/components/ui/card";

  let { onregistered }: { onregistered?: () => void | Promise<void> } =
    $props();
  let username = $state("");
  let password = $state("");
  let confirmation = $state("");
  let masterPassword = $state("");
  let masterConfirmation = $state("");
  let pending = $state(false);
  let error = $state("");
  const valid = $derived(
    username.trim().length > 0 &&
      password.length >= 8 &&
      password === confirmation &&
      masterPassword.length >= 8 &&
      masterPassword === masterConfirmation,
  );

  async function submit(event: SubmitEvent) {
    event.preventDefault();
    if (!valid || pending) return;
    pending = true;
    error = "";
    try {
      await signupAccount(username.trim(), password, masterPassword);
      password = "";
      confirmation = "";
      masterPassword = "";
      masterConfirmation = "";
      await onregistered?.();
    } catch {
      error =
        "Registration failed. Check your username and server connection, then try again.";
    } finally {
      pending = false;
    }
  }
</script>

<Card.Root>
  <Card.Header>
    <Card.Title>Create an account</Card.Title>
    <Card.Description
      >Create your account and encrypted local vault. For an existing
      unregistered vault, enter its current master password.</Card.Description
    >
  </Card.Header>
  <Card.Content>
    <form onsubmit={submit}>
      <Field.Group>
        <Field.Field>
          <Field.Label for="username-{id}">Username</Field.Label>
          <Input
            id="username-{id}"
            bind:value={username}
            autocomplete="username"
            required
            disabled={pending}
          />
        </Field.Field>
        <Field.Field>
          <Field.Label for="password-{id}">Account password</Field.Label>
          <Input
            id="password-{id}"
            bind:value={password}
            type="password"
            autocomplete="new-password"
            minlength={8}
            required
            disabled={pending}
          />
          <Field.Description
            >Must be at least 8 characters long.</Field.Description
          >
        </Field.Field>
        <Field.Field>
          <Field.Label for="confirm-{id}">Confirm password</Field.Label>
          <Input
            id="confirm-{id}"
            bind:value={confirmation}
            type="password"
            autocomplete="new-password"
            required
            disabled={pending}
          />
          {#if confirmation && password !== confirmation}<Field.Error
              >Passwords do not match.</Field.Error
            >{/if}
        </Field.Field>
        <Field.Field>
          <Field.Label for="master-{id}">Master password</Field.Label>
          <Input
            id="master-{id}"
            bind:value={masterPassword}
            type="password"
            autocomplete="new-password"
            minlength={8}
            required
            disabled={pending}
          />
          <Field.Description
            >Encrypts your local vault. Keep it safe; it is separate from your
            account password.</Field.Description
          >
        </Field.Field>
        <Field.Field>
          <Field.Label for="master-confirm-{id}"
            >Confirm master password</Field.Label
          >
          <Input
            id="master-confirm-{id}"
            bind:value={masterConfirmation}
            type="password"
            autocomplete="new-password"
            required
            disabled={pending}
          />
        </Field.Field>
        {#if error}<Field.Error>{error}</Field.Error>{/if}
        <Field.Field>
          <Button type="submit" disabled={!valid || pending}
            >{pending ? "Registering…" : "Create Account"}</Button
          >
        </Field.Field>
      </Field.Group>
    </form>
  </Card.Content>
</Card.Root>
