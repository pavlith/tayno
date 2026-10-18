<script lang="ts">
  import { ModeWatcher } from "mode-watcher";
  import AppRoute from "@/AppRoute.svelte";
</script>

<ModeWatcher defaultMode="system" />
<AppRoute />
