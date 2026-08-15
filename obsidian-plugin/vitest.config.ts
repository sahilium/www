import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    environment: "node",
    globals: true,
    include: ["**/*.test.ts"],
    coverage: {
      provider: "v8",
      include: ["main.ts"],
      exclude: ["main.test.ts"],
      reporter: ["text", "lcov", "json-summary"],
    },
  },
});