<script lang="ts">
  import { messages, MessageType } from "@/stores/messages";
  import { slide } from "svelte/transition";
</script>

<div
  class="pointer-events-none fixed right-4 bottom-4 z-[100] grid w-[min(22rem,calc(100vw-2rem))] gap-2"
  aria-live="polite"
  aria-atomic="false"
>
  {#each $messages as message}
    <div
      class="rounded-xl border border-border bg-popover px-4 py-3 text-sm shadow-lg"
      transition:slide
    >
      <div
        class={message.type === MessageType.Error
          ? "font-medium text-destructive"
          : "font-medium text-primary"}
      >
        {message.title}
      </div>
      {#if message.message}<p class="mt-1 text-muted-foreground">
          {message.message}
        </p>{/if}
    </div>
  {/each}
</div>
