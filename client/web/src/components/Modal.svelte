<script lang="ts">
  import type { Snippet } from "svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  let {
    show = $bindable(false),
    header,
    children,
    footer,
  }: {
    show?: boolean;
    header?: Snippet;
    children?: Snippet;
    footer?: Snippet;
  } = $props();
</script>

<Dialog.Root bind:open={show}>
  <Dialog.Content
    portalProps={{ disabled: true }}
    class="max-h-[90dvh] overflow-y-auto p-6 sm:max-w-lg"
  >
    <Dialog.Header>
      <Dialog.Title
        >{#if header}{@render header()}{:else}Myst{/if}</Dialog.Title
      >
      <Dialog.Description class="sr-only"
        >Manage your encrypted keystore.</Dialog.Description
      >
    </Dialog.Header>
    <div class="min-w-0">{@render children?.()}</div>
    {#if footer}<Dialog.Footer>{@render footer()}</Dialog.Footer>{/if}
  </Dialog.Content>
</Dialog.Root>
