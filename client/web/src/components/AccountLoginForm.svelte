<script lang="ts">
  import { loginAccount } from "@/stores/authState";
  import SignupForm from "$lib/components/signup-form.svelte";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import * as Card from "$lib/components/ui/card";
  import * as Field from "$lib/components/ui/field";
  let { onlogin }: { onlogin: () => Promise<void> } = $props();
  const id = $props.id();
  let signup = $state(false);
  let username = $state("");
  let password = $state("");
  let pending = $state(false);
  let error = $state("");
  async function submit(event: SubmitEvent) {
    event.preventDefault();
    if (pending || !username.trim() || !password) return;
    pending = true;
    error = "";
    try {
      await loginAccount(username.trim(), password);
      password = "";
      await onlogin();
    } catch (err: any) {
      error =
        err.status === 401
          ? "Incorrect username or account password."
          : "Unable to connect to the account server. Please try again.";
    } finally {
      pending = false;
    }
  }
</script>

<div class="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
  <div class="w-full max-w-sm">
    {#if signup}
      <SignupForm onregistered={onlogin} />
      <Button
        variant="link"
        class="mt-3 w-full"
        onclick={() => (signup = false)}>Already have an account? Login</Button
      >
    {:else}
      <Card.Root>
        <Card.Header>
          <Card.Title class="text-2xl">Login</Card.Title>
          <Card.Description
            >Sign in to your Myst account. You will unlock your vault next.</Card.Description
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
                  autocomplete="current-password"
                  required
                  disabled={pending}
                />
              </Field.Field>
              {#if error}<Field.Error>{error}</Field.Error>{/if}
              <Field.Field
                ><Button
                  type="submit"
                  disabled={pending || !username.trim() || !password}
                  >{pending ? "Signing in…" : "Login"}</Button
                ></Field.Field
              >
            </Field.Group>
          </form>
          <Button
            variant="link"
            class="mt-3 w-full"
            onclick={() => (signup = true)}>Create an account</Button
          >
        </Card.Content>
      </Card.Root>
    {/if}
  </div>
</div>
