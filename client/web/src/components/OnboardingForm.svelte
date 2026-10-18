<script lang="ts">
  import { onMount } from "svelte";
  import api, { OpenAPI } from "@/api";
  import PasswordInputField from "./PasswordInputField.svelte";
  import { Button } from "$lib/components/ui/button";
  import * as Card from "$lib/components/ui/card";

  let { oninitialized }: { oninitialized?: () => void } = $props();
  let password = $state("");
  let passwordInput = $state<HTMLInputElement | null>(null);
  let pending = $state(false);
  let error = $state("");
  const valid = $derived(password.length >= 8);
  onMount(() => passwordInput?.focus());

  async function submit(event: SubmitEvent) {
    event.preventDefault();
    if (!valid || pending) return;
    pending = true;
    error = "";
    try {
      OpenAPI.TOKEN = await api.createEnclave({ requestBody: { password } });
      password = "";
      oninitialized?.();
    } catch {
      error = "Setup failed. Please try again.";
    } finally {
      pending = false;
    }
  }
</script>

<div class="flex min-h-dvh w-full items-center justify-center p-4 sm:p-8">
  <Card.Root class="w-full max-w-md py-6">
    <Card.Header class="space-y-2 px-6">
      <div
        class="text-xs font-semibold uppercase tracking-[0.25em] text-primary"
      >
        Myst
      </div>
      <Card.Title class="text-2xl">Protect your secrets</Card.Title>
      <Card.Description
        >Set a master password to encrypt your data. Choose a strong password
        and keep it safe.</Card.Description
      >
    </Card.Header>
    <Card.Content class="px-6">
      <form onsubmit={submit}>
        <PasswordInputField
          bind:value={password}
          bind:ref={passwordInput}
          label="Master password"
          placeholder="At least 8 characters"
          autocomplete="new-password"
          disabled={pending}
        />
        {#if error}<p role="alert" class="mb-4 text-sm text-destructive">
            {error}
          </p>{/if}
        <Button type="submit" class="w-full" disabled={!valid || pending}
          >{pending ? "Creating…" : "Create encrypted vault"}</Button
        >
      </form>
    </Card.Content>
  </Card.Root>
</div>
