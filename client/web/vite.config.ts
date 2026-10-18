import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { fileURLToPath, URL } from "node:url";

export default defineConfig({
  envDir: fileURLToPath(new URL("../../", import.meta.url)),
  build: { outDir: "../../static", emptyOutDir: true },
  server: {
    port: 8082,
    proxy: {
      "/api": {
        target: process.env.VITE_API_PROXY_TARGET || "http://localhost:8081",
        changeOrigin: true,
      },
    },
  },
  plugins: [tailwindcss(), svelte()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
      $lib: fileURLToPath(new URL("./src/lib", import.meta.url)),
    },
  },
});
