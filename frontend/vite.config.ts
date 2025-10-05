import { defineConfig } from 'vite';

export default defineConfig({
    build: {
        outDir: 'dist',
    },
    test: {
        globals: true,
    },
    server: {
        proxy: {
            // string shorthand: http://localhost:5173/api -> http://localhost:8080/api
            '/api': {
                target: 'http://localhost:8080', // Your Go backend server
                changeOrigin: true, // Needed for virtual hosted sites
                secure: false,      // Can be false if you're not using HTTPS
                ws: true,           // If you want to proxy websockets
            }
        }
    },
});