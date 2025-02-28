import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    disableConsoleIntercept: true,
    environment: 'node',
    hookTimeout: 30000,
    testTimeout: 30000,
  },
});
