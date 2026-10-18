<script lang="ts">
  import LogoutButton from "@/components/LogoutButton.svelte";
  import api, { OpenAPI } from "@/api";
  import { Input } from "$lib/components/ui/input";
  import {
    FieldGroup,
    Field,
    FieldLabel,
    FieldError,
  } from "$lib/components/ui/field";
  const id = $props.id();
  import { Button } from "$lib/components/ui/button";
  import * as Card from "$lib/components/ui/card";

  let { onlogin }: { onlogin?: () => void | Promise<void> } = $props();
  let password = $state("");
  let error = $state("");
  let pending = $state(false);
  const valid = $derived(password.length >= 8);
  async function submit(event: SubmitEvent) {
    event.preventDefault();
    if (!valid || pending) return;
    pending = true;
    error = "";
    try {
      OpenAPI.TOKEN = await api.authenticate({ requestBody: { password } });
      password = "";
      await onlogin?.();
    } catch (err) {
      error =
        err && typeof err === "object" && "status" in err && err.status === 401
          ? "Incorrect master password."
          : err &&
              typeof err === "object" &&
              "status" in err &&
              err.status === 403
            ? "This vault belongs to a different account. Logout and sign in with its owner."
            : "Unable to unlock. Check your master password and server connection.";
      password = "";
    } finally {
      pending = false;
    }
  }
</script>

<Card.Root class="mx-auto w-full max-w-sm">
  <Card.Header>
    <Card.Title class="text-2xl">Unlock vault</Card.Title>
    <Card.Description>Unlock Myst with your master password.</Card.Description>
  </Card.Header>
  <Card.Content>
    <form onsubmit={submit}>
      <FieldGroup>
        <Field>
          <FieldLabel for="password-{id}">Master password</FieldLabel>
          <Input
            id="password-{id}"
            bind:value={password}
            type="password"
            autocomplete="current-password"
            placeholder="Enter your master password"
            required
            minlength={8}
            disabled={pending}
            aria-invalid={!!error}
            aria-describedby={error ? `error-${id}` : undefined}
          />
          {#if error}<FieldError id="error-{id}">{error}</FieldError>{/if}
        </Field>
        <Field>
          <Button type="submit" class="w-full" disabled={!valid || pending}>
            {pending ? "Unlocking…" : "Unlock vault"}
          </Button>
        </Field>
      </FieldGroup>
    </form>
    <div class="mt-4 grid"><LogoutButton /></div>
  </Card.Content>
</Card.Root>
