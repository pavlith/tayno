<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import api from "@/api";
  import type { Invitation } from "@/api/generated";
  import AcceptInvitationModal from "@/components/AcceptInvitationModal.svelte";
  import DeclineInvitationModal from "@/components/DeclineInvitationModal.svelte";
  import FinalizeInvitationModal from "@/components/FinalizeInvitationModal.svelte";
  import { hash } from "@/lib/color-hash";
  import { format } from "@/lib/date";
  import { getInvitations, invitations } from "@/stores/invitations";
  import { currentUser } from "@/stores/user";
  import { useFocus } from "@/lib/router";
  import { onMount } from "svelte";
  import VerifyInvitationModal from "@/components/VerifyInvitationModal.svelte";

  const registerFocus = useFocus();

  let incomingInvitations = $derived(
    $invitations.filter(
      (inv) =>
        inv.invitee.id === $currentUser?.id &&
        (inv.status === "pending" || inv.status === "accepted"),
    ),
  );
  let outgoingInvitations = $derived(
    $invitations.filter(
      (inv) =>
        inv.inviter.id === $currentUser?.id &&
        (inv.status === "pending" || inv.status === "accepted"),
    ),
  );
  let pastInvitations = $derived(
    $invitations.filter(
      (inv) => inv.status !== "pending" && inv.status !== "accepted",
    ),
  );

  let invitation = $state<Invitation>();
  let showAcceptInvitationModal: boolean = $state(false);
  let showVerifyInvitationModal: boolean = $state(false);
  let showDeclineInvitationModal: boolean = $state(false);
  let showFinalizeInvitationModal: boolean = $state(false);

  function acceptInvitation() {
    api
      .acceptInvitation({
        invitationId: invitation!.id,
      })
      .then(() => {
        showAcceptInvitationModal = false;
        getInvitations();
      });
  }

  function declineInvitation() {
    api
      .declineOrCancelInvitation({
        invitationId: invitation!.id,
      })
      .then(() => {
        showDeclineInvitationModal = false;
        getInvitations();
      });
  }

  function finalizeInvitation() {
    api
      .finalizeInvitation({
        invitationId: invitation!.id,
        requestBody: {
          remoteKeystoreId: invitation!.keystore.remoteId,
          inviteePublicKey: invitation!.invitee.publicKey,
        },
      })
      .then(() => {
        showFinalizeInvitationModal = false;
        getInvitations();
      });
  }

  function showAcceptInvitationModalFunc(inv: Invitation) {
    invitation = inv;
    showAcceptInvitationModal = true;
  }

  function showVerifyInvitationModalFunc(inv: Invitation) {
    invitation = inv;
    showVerifyInvitationModal = true;
  }

  function showDeclineInvitationModalFunc(inv: Invitation) {
    invitation = inv;
    showDeclineInvitationModal = true;
  }

  function showFinalizeInvitationModalFunc(inv: Invitation) {
    invitation = inv;
    showFinalizeInvitationModal = true;
  }

  onMount(() => {
    getInvitations();
  });
</script>

<div class="invitations-list" use:registerFocus>
  <h1 class="mb-6 text-xl font-semibold">Invitations</h1>
  {#if $invitations.length === 0}<p class="text-muted-foreground">
      No invitations yet. Share a keystore to invite someone.
    </p>{/if}
  {#if incomingInvitations.length > 0}
    <div class="section">
      <h5>Incoming Invitations</h5>
      {#each incomingInvitations as inv}
        <div class="invitation">
          <span class="icon">
            <span style="background-color: {hash(inv.inviter.username)}">
              {inv.inviter.username.slice(0, 2).toUpperCase()}
            </span>
          </span>
          <div class="info">
            <span class="name">
              {inv.keystore.name}
            </span>
            <span class="user">
              {#if inv.status === "accepted"}
                Invitation from <strong
                  style="color: {hash(inv.inviter.username)}"
                  >{inv.inviter.username}</strong
                >
                accepted {format(inv.acceptedAt)}.
              {:else}
                Invited by <strong style="color: {hash(inv.inviter.username)}"
                  >{inv.inviter.username}</strong
                >
                {format(inv.createdAt)}.
              {/if}
            </span>
          </div>

          <div class="actions">
            {#if inv.status === "pending"}
              <Button
                variant="default"
                class="button green"
                onclick={() => {
                  showAcceptInvitationModalFunc(inv);
                }}>Accept</Button
              >
              <Button
                variant="destructive"
                class="button red"
                onclick={() => {
                  showDeclineInvitationModalFunc(inv);
                }}>Decline</Button
              >
            {:else if inv.status === "accepted"}
              <Button
                variant="default"
                class="button green"
                onclick={() => {
                  showVerifyInvitationModalFunc(inv);
                }}>Verify</Button
              >
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}

  {#if outgoingInvitations.length > 0}
    <div class="section">
      <h5>Outgoing Invitations</h5>

      {#each outgoingInvitations as inv}
        <div class="invitation">
          <span class="icon">
            <span style="background-color: {hash(inv.invitee.username)}">
              {inv.invitee.username.slice(0, 2).toUpperCase()}
            </span>
          </span>
          <div class="info">
            <span class="name">
              {inv.keystore.name}
            </span>
            <span class="user">
              Invited <strong style="color: {hash(inv.invitee.username)}"
                >{inv.invitee.username}</strong
              >
              {format(inv.createdAt)}.
            </span>
          </div>

          <div class="actions">
            {#if inv.status === "accepted"}
              <Button
                variant="default"
                class="button green"
                onclick={() => {
                  showFinalizeInvitationModalFunc(inv);
                }}>Verify</Button
              >
            {/if}
            <Button
              variant="destructive"
              class="button red"
              onclick={() => {
                showDeclineInvitationModalFunc(inv);
              }}>Delete</Button
            >
          </div>
        </div>
      {/each}
    </div>
  {/if}

  {#if pastInvitations.length > 0}
    <div class="section">
      <h5>Past Invitations</h5>
      {#each pastInvitations as inv}
        <div class="invitation">
          <span class="icon">
            <span
              style="background-color: {hash(
                inv.inviter.id === $currentUser?.id
                  ? inv.invitee.username
                  : inv.inviter.username,
              )}"
            >
              {#if inv.inviter.id === $currentUser?.id}
                {inv.invitee.username.slice(0, 2).toUpperCase()}
              {:else}
                {inv.inviter.username.slice(0, 2).toUpperCase()}
              {/if}
            </span>
          </span>
          <div class="info">
            <span class="name">
              {inv.keystore.name}
            </span>
            <span class="user">
              {#if inv.status === "finalized"}
                {#if inv.inviter.id === $currentUser?.id}
                  Shared with <strong
                    style="color: {hash(inv.invitee.username)}"
                    >{inv.invitee.username}</strong
                  >
                  since {format(inv.acceptedAt)}.
                {:else}
                  Being shared by <strong
                    style="color: {hash(inv.inviter.username)}"
                    >{inv.inviter.username}</strong
                  >
                  since {format(inv.acceptedAt)}.
                {/if}
              {:else if inv.status === "deleted"}
                {#if inv.inviter.id === $currentUser?.id}
                  Invitation to <strong
                    style="color: {hash(inv.invitee.username)}"
                    >{inv.invitee.username}</strong
                  >
                  deleted {format(inv.deletedAt)}.
                {/if}
              {:else if inv.status === "declined"}
                {#if inv.inviter.id === $currentUser?.id}
                  Invitation declined by <strong
                    style="color: {hash(inv.invitee.username)}"
                    >{inv.invitee.username}</strong
                  >
                  {format(inv.declinedAt)}.
                {:else}
                  Declined invitation from <strong
                    style="color: {hash(inv.inviter.username)}"
                    >{inv.inviter.username}</strong
                  >
                  {format(inv.declinedAt)}.
                {/if}
              {/if}
            </span>
          </div>
          <div class="actions"></div>
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if invitation}
  <AcceptInvitationModal
    bind:show={showAcceptInvitationModal}
    {invitation}
    onsubmit={() => {
      acceptInvitation();
    }}
  />
  <DeclineInvitationModal
    bind:show={showDeclineInvitationModal}
    {invitation}
    onsubmit={() => {
      declineInvitation();
    }}
  />
  <FinalizeInvitationModal
    bind:show={showFinalizeInvitationModal}
    {invitation}
    onsubmit={() => {
      finalizeInvitation();
    }}
  />
  <VerifyInvitationModal bind:show={showVerifyInvitationModal} {invitation} />
{/if}

<style>
  .invitations-list {
    width: 100%;
    padding: 32px;
    overflow-y: auto;
  }
  .section {
    margin-bottom: 32px;
  }
  h5 {
    margin-bottom: 12px;
    color: var(--muted-foreground);
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }
  .invitation {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 16px;
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 16px;
    margin-bottom: 12px;
  }
  .icon > span {
    display: grid;
    place-items: center;
    height: 40px;
    width: 40px;
    border-radius: 12px;
    color: var(--background);
    font-weight: 600;
  }
  .info {
    flex: 1;
    min-width: 160px;
    display: grid;
    gap: 4px;
  }
  .name {
    font-weight: 500;
  }
  .user {
    color: var(--muted-foreground);
    font-size: 13px;
  }
  .actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }
</style>
