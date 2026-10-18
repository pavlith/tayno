<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import Keystore from "@/components/Keystore.svelte";
  import { useNavigate, useParams } from "@/lib/router";
  import type { Keystore as KeystoreModel } from "@/api";

  let {
    keystores,
    showCreateKeystoreModal = $bindable(false),
  }: {
    keystores: KeystoreModel[];
    showCreateKeystoreModal?: boolean;
  } = $props();
  const navigate = useNavigate();
  const params = useParams();
  const keystore = $derived(
    keystores.find((store) => store.id === $params.keystoreId),
  );

  $effect(() => {
    if (!keystore && keystores.length) {
      navigate(`/keystore/${keystores[0].id}`, { replace: true });
    } else if (!keystores.length && $params.keystoreId) {
      navigate("/", { replace: true });
    }
  });
</script>

{#if keystore}
  <Keystore {keystore} />
{:else}
  <div
    class="empty flex min-h-64 w-full flex-col items-center justify-center gap-4 p-6 text-center"
  >
    <h1 class="text-lg font-medium">Create a new keystore to get started</h1>
    <Button onclick={() => (showCreateKeystoreModal = true)}>Create</Button>
  </div>
{/if}
