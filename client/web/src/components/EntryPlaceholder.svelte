<script lang="ts">
  import KeyRound from "@lucide/svelte/icons/key-round";
</script>

<div
  class="entry-placeholder flex min-w-0 basis-1/2 flex-col items-center justify-center gap-3 border-l border-border p-8 text-center"
>
  <div class="rounded-2xl bg-muted p-4">
    <KeyRound class="size-7 text-primary" />
  </div>
  <h2 class="text-base font-medium">Your secrets, organized</h2>
  <p class="max-w-xs text-sm text-muted-foreground">
    Select an entry to view its details, or create one to get started.
  </p>
</div>
