<script lang="ts">
  import LoginForm from "$lib/components/login-form.svelte";
  let { onlogin }: { onlogin?: () => void } = $props();
</script>

<div class="flex min-h-svh w-full items-center justify-center px-4 py-8">
  <LoginForm {onlogin} />
</div>
