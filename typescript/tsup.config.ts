import { defineConfig } from 'tsup';

export default defineConfig({
  splitting: true,
  clean: true,
  format: ['esm', 'cjs'],
  bundle: true,
  skipNodeModulesBundle: true,
  watch: false,
  shims: true,
  entry: {
    index: 'radius/index.ts',
  },
  treeshake: true,
  outDir: 'dist',
  dts: true,
  minify: process.env.NODE_ENV === 'production',
  sourcemap: process.env.NODE_ENV !== 'production',
});
