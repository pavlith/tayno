<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import type { Keystore } from "@/api";
  let { keystores, pending }: { keystores: Keystore[]; pending: number } =
    $props();
  const entries = $derived(
    keystores.reduce((total, store) => total + store.entries.length, 0),
  );
</script>

<div class="grid grid-cols-3 gap-2 sm:gap-3">
  {#each [{ label: "Keystores", value: keystores.length }, { label: "Entries", value: entries }, { label: "Pending invitations", value: pending }] as stat}
    <Card.Root class="gap-2 py-3 sm:py-4"
      ><Card.Header class="px-3 sm:px-6"
        ><Card.Description class="text-xs sm:text-sm"
          >{stat.label}</Card.Description
        ><Card.Title class="text-2xl tabular-nums">{stat.value}</Card.Title
        ></Card.Header
      ></Card.Root
    >
  {/each}
</div>
