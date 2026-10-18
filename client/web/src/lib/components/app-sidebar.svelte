<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar";
  import Shield from "@lucide/svelte/icons/shield-check";
  import Key from "@lucide/svelte/icons/key-round";
  import Mail from "@lucide/svelte/icons/mail";
  import Plus from "@lucide/svelte/icons/plus";
  import LogOut from "@lucide/svelte/icons/log-out";
  import { useNavigate, useLocation } from "@/lib/router";
  import { currentUser } from "@/stores/user";
  import { invitations } from "@/stores/invitations";
  import { logout } from "@/stores/authState";
  import { showError } from "@/stores/messages";
  import type { Keystore } from "@/api";
  let { keystores, oncreate }: { keystores: Keystore[]; oncreate: () => void } =
    $props();
  const navigate = useNavigate();
  const location = useLocation();
  const sidebar = Sidebar.useSidebar();
  function go(path: string) {
    navigate(path);
    sidebar.setOpenMobile(false);
  }
  let pending = $state(false);
  async function signOut() {
    if (pending) return;
    pending = true;
    try {
      await logout();
    } catch {
      showError("Unable to log out. Please try again.");
    } finally {
      pending = false;
    }
  }
</script>

<Sidebar.Root collapsible="icon" variant="inset">
  <Sidebar.Header>
    <Sidebar.Menu
      ><Sidebar.MenuItem>
        <Sidebar.MenuButton
          size="lg"
          onclick={() => go("/")}
          tooltipContent="Myst"
        >
          <div
            class="flex size-8 shrink-0 items-center justify-center rounded-lg bg-primary text-primary-foreground"
          >
            <Shield />
          </div>
          <span class="font-semibold">Myst</span>
        </Sidebar.MenuButton>
      </Sidebar.MenuItem></Sidebar.Menu
    >
  </Sidebar.Header>
  <Sidebar.Content>
    <Sidebar.Group>
      <Sidebar.GroupLabel>Workspace</Sidebar.GroupLabel>
      <Sidebar.Menu>
        <Sidebar.MenuItem
          ><Sidebar.MenuButton
            onclick={oncreate}
            tooltipContent="Create keystore"
            ><Plus /><span>Create keystore</span></Sidebar.MenuButton
          ></Sidebar.MenuItem
        >
        <Sidebar.MenuItem
          ><Sidebar.MenuButton
            onclick={() => go("/invitations")}
            isActive={$location.pathname === "/invitations"}
            tooltipContent="Invitations"
            ><Mail /><span
              >Invitations ({$invitations.filter((i) => i.status === "pending")
                .length})</span
            ></Sidebar.MenuButton
          ></Sidebar.MenuItem
        >
      </Sidebar.Menu>
    </Sidebar.Group>
    <Sidebar.Group>
      <Sidebar.GroupLabel>Keystores</Sidebar.GroupLabel>
      <Sidebar.Menu>
        {#each keystores as store (store.id)}
          <Sidebar.MenuItem
            ><Sidebar.MenuButton
              onclick={() => go(`/keystore/${store.id}`)}
              isActive={$location.pathname.startsWith(`/keystore/${store.id}`)}
              tooltipContent={store.name}
              ><Key /><span>{store.name}</span></Sidebar.MenuButton
            ></Sidebar.MenuItem
          >
        {/each}
      </Sidebar.Menu>
    </Sidebar.Group>
  </Sidebar.Content>
  <Sidebar.Footer>
    <p
      class="truncate px-2 text-xs text-muted-foreground group-data-[collapsible=icon]:hidden"
    >
      {$currentUser?.username ?? "Myst"}
    </p>
    <Sidebar.Menu
      ><Sidebar.MenuItem
        ><Sidebar.MenuButton
          onclick={signOut}
          aria-disabled={pending}
          tooltipContent="Logout"
          ><LogOut /><span>{pending ? "Logging out…" : "Logout"}</span
          ></Sidebar.MenuButton
        ></Sidebar.MenuItem
      ></Sidebar.Menu
    >
  </Sidebar.Footer>
  <Sidebar.Rail />
</Sidebar.Root>
