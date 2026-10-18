<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar";
  import { Separator } from "$lib/components/ui/separator";
  import { Button } from "$lib/components/ui/button";
  import { toggleMode } from "mode-watcher";
  import Sun from "@lucide/svelte/icons/sun";
  import Moon from "@lucide/svelte/icons/moon";
  let { title }: { title: string } = $props();
</script>

<header class="flex h-14 shrink-0 items-center gap-3 border-b px-4 lg:px-6">
  <Sidebar.Trigger />
  <Separator orientation="vertical" class="h-4" />
  <h1 class="min-w-0 flex-1 truncate text-sm font-medium">{title}</h1>
  <Button
    variant="ghost"
    size="icon"
    onclick={toggleMode}
    aria-label="Toggle color theme"
  >
    <Sun class="hidden size-4 dark:block" /><Moon class="size-4 dark:hidden" />
  </Button>
</header>
