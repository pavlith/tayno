<script lang="ts">
  let { invitation }: { invitation: import("@/api").Invitation } = $props();
</script>

<div class="rounded-lg border border-border p-4">
  <p class="font-medium">{invitation.keystore.name}</p>
  <p class="text-sm text-muted-foreground">{invitation.status}</p>
</div>
