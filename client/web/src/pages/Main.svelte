<script lang="ts">
  import CreateKeystoreModal from "@/components/CreateKeystoreModal.svelte";
  import * as Sidebar from "$lib/components/ui/sidebar";
  import AppSidebar from "$lib/components/app-sidebar.svelte";
  import SiteHeader from "$lib/components/site-header.svelte";
  import SectionCards from "$lib/components/section-cards.svelte";
  import { invitations } from "@/stores/invitations";
  import Invitations from "@/pages/Invitations.svelte";
  import Keystores from "@/pages/Keystores.svelte";
  import { getKeystores, keystores } from "@/stores/keystores";
  import { useLocation } from "@/lib/router";
  const location = useLocation();
  import { onDestroy, onMount } from "svelte";
  import { getCurrentUser } from "@/stores/user";
  import { getInvitations } from "@/stores/invitations";
  import { useNavigate } from "@/lib/router";

  const navigate = useNavigate();

  let showCreateKeystoreModal: boolean = $state(false);

  let ready: boolean = $state(false);

  onMount(async () => {
    await refresh();
    ready = true;
  });

  const refresh = async () => {
    await getKeystores();
    const u = await getCurrentUser();
    if (u) {
      await getInvitations();
    }
  };

  const interval = setInterval(refresh, 5000);

  onDestroy(() => clearInterval(interval));

  const onKeystoreCreated = async (keystore: { id: string }) => {
    await refresh();
    navigate(`/keystore/${keystore.id}`);
  };
</script>

{#if ready}
  <Sidebar.Provider class="h-svh min-h-0 overflow-hidden">
    <AppSidebar
      keystores={$keystores}
      oncreate={() => (showCreateKeystoreModal = true)}
    />
    <Sidebar.Inset class="min-h-0 min-w-0 overflow-hidden">
      <SiteHeader
        title={$location.pathname === "/invitations"
          ? "Invitations"
          : "Keystores"}
      />
      <div class="shrink-0 p-4 lg:px-6">
        <SectionCards
          keystores={$keystores}
          pending={$invitations.filter((i) => i.status === "pending").length}
        />
      </div>
      <div class="workspace flex min-h-0 flex-1 overflow-auto border-t">
        {#if $location.pathname === "/invitations"}
          <Invitations />
        {:else}
          <Keystores keystores={$keystores} bind:showCreateKeystoreModal />
        {/if}
      </div>
    </Sidebar.Inset>
  </Sidebar.Provider>

  <CreateKeystoreModal
    bind:show={showCreateKeystoreModal}
    oncreated={(e) => {
      onKeystoreCreated(e);
    }}
  />
{/if}
