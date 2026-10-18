# Myst UI

Svelte 5 + TypeScript + Vite, Tailwind CSS 4 and shadcn-svelte (Nova).
Node.js 22.12+ or 24+ is required by Vite 8.

```sh
npm ci
npm run dev
npm run check
npm run build
```

The development UI runs on port 8082 and calls the API on port 8081.
Production builds are written to `../../static` for the existing Go server.

Components use `$state`, `$derived`, `$effect`, `$props` and `$bindable`.
Component events are callback props, and content is passed as snippets.
Shared API stores remain Svelte stores so polling and subscriptions retain
one source of truth. The small History API router in `src/lib/router.ts`
supports deep links and browser back/forward navigation.

The shadcn components live in `src/lib/components/ui`. Theme tokens are in
`src/app.css`; aliases and registry configuration are in `components.json`.
To add another component:

```sh
npx shadcn-svelte@latest add <component>
```

`Modal.svelte` wraps the shadcn Dialog. Its portal is disabled deliberately:
submit controls must remain inside their enclosing form. Dialog focus trapping,
Escape dismissal and focus restoration are provided by Bits UI.

References: [Svelte 5 migration](https://svelte.dev/docs/svelte/v5-migration-guide),
[shadcn-svelte with Vite](https://www.shadcn-svelte.com/docs/installation/vite).
