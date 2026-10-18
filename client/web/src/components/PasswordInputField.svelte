<script lang="ts">
  import InputField from "./InputField.svelte";
  import type { ComponentProps } from "svelte";
  let {
    value = $bindable(""),
    ref = $bindable(null),
    type: _type,
    autocomplete = "current-password",
    ...rest
  }: ComponentProps<typeof InputField> = $props();
</script>

<InputField bind:value bind:ref {...rest} type="password" {autocomplete} />
