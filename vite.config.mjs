import { defineConfig } from "vite";

export default defineConfig({
  build: {
    lib: {
      entry: "browser.mjs",
      name: "XheWC",
      formats: ["umd"],
    },
  },
});
