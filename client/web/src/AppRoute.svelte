<script lang="ts">
  import AccountLoginForm from "@/components/AccountLoginForm.svelte";
  import LoginForm from "@/components/LoginForm.svelte";
  import Messages from "@/components/Messages.svelte";
  import Main from "@/pages/Main.svelte";
  import { onDestroy, onMount } from "svelte";
  import { useLocation, useNavigate } from "@/lib/router";
  import api from "@/api";
  import { AuthState, authState, getAuthState } from "@/stores/authState";
  import { showError } from "@/stores/messages";

  const location = useLocation();
  const navigate = useNavigate();

  $effect(() => {
    const state = $authState;

    if (
      (state === AuthState.Locked || state === AuthState.SignedOut) &&
      $location.pathname !== "/"
    ) {
      navigate("/", { replace: true });
    }
  });

  let checking = false;
  const checkState = async () => {
    if (checking) return;
    checking = true;
    try {
      const oldState = $authState;
      const state = await getAuthState();
      if (state === AuthState.SignedIn) {
        await api.healthCheck();
      } else if (
        state === AuthState.Locked &&
        oldState === AuthState.SignedIn
      ) {
        // was logged out, notify
        showError("Vault locked", "Unlock your vault to continue.");
      }
    } catch {
      // Keep the current screen during a temporary connection failure.
    } finally {
      checking = false;
    }
  };

  const interval = setInterval(checkState, 1000);

  onDestroy(() => clearInterval(interval));

  onMount(async () => {
    await checkState();
  });
</script>

{#if $authState !== undefined}
  {#if $authState === AuthState.SignedOut}
    <AccountLoginForm onlogin={checkState} />
  {:else if $authState === AuthState.Locked}
    <LoginForm onlogin={checkState} />
  {:else if $authState === AuthState.SignedIn}
    <Main />
  {/if}
{/if}

<Messages />
