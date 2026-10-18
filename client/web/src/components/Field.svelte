<script lang="ts">
  let {
    label,
    value,
    class: className = "",
  }: { label: string; value?: string; class?: string } = $props();
</script>

<div class={`mb-6 min-w-0 ${className}`}>
  <div
    class="mb-2 text-xs font-medium uppercase tracking-wider text-muted-foreground"
  >
    {label}
  </div>
  <div class="whitespace-pre-wrap break-all text-sm leading-relaxed">
    {value || "—"}
  </div>
</div>
