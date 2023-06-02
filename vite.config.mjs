import { defineConfig } from "vite";

export default defineConfig({
  build: {
    lib: {
      entry: ".",
      name: "XheWC",
      formats: ["umd"],
    },
  },
});
