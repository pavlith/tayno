<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import type { HTMLInputAttributes } from "svelte/elements";
  const uid = $props.id();
  let {
    ref = $bindable(null),
    value = $bindable(""),
    label = "",
    error = "",
    class: className = "",
    id = uid,
    type = "text",
    ...rest
  }: Omit<HTMLInputAttributes, "value" | "type" | "files"> & {
    type?: Exclude<HTMLInputAttributes["type"], null | "file">;
    ref?: HTMLInputElement | null;
    value?: string;
    label?: string;
    error?: string | false;
  } = $props();
</script>

<div class={`field mb-5 grid gap-2 ${className}`}>
  {#if label}<Label for={id} class="leading-relaxed">{label}</Label>{/if}
  <Input
    bind:ref
    bind:value
    {id}
    {type}
    {...rest}
    aria-invalid={!!error}
    aria-describedby={error ? `${id}-error` : undefined}
    class="h-10"
  />
  {#if error}<p
      id={`${id}-error`}
      class="text-sm text-destructive"
      role="alert"
    >
      {error}
    </p>{/if}
</div>
