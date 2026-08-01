import { defineConfig } from 'vite'
import path from "path"
import react from '@vitejs/plugin-react-swc'
import tailwindcss from '@tailwindcss/vite'


// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  server: {
    host: true,
    // Stands in for Traefik in local development: serves the API under the same origin at
    // /api, strips the prefix, and sets the Remote-User header that Authelia's ForwardAuth
    // adds in production. Keeps dev and prod on the same request shape.
    proxy: {
      "/api": {
        target: process.env.VITE_DEV_API_TARGET ?? "http://localhost:8000",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
        headers: {
          "Remote-User": process.env.VITE_DEV_OWNER_USERNAME ?? "dev-owner",
        },
      },
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  }
})
