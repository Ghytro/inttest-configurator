import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";
import { pluginSvgr } from "@rsbuild/plugin-svgr";

export default defineConfig({
  plugins: [pluginReact(), pluginSvgr({ mixedImport: true })],
  html: {
    template: "./public/index.html",
  },
  server: {
    proxy: {
      "/api": "http://localhost:8080"
    }
  }
});
