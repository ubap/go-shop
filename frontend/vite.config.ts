import {defineConfig} from 'vite';
import {resolve} from 'path';

// export default is needed for vite config
export default defineConfig({
    build: {
        outDir: 'dist',
        rollupOptions: {
            input: {
                main: resolve(__dirname, 'index.html'),
                another: resolve(__dirname, 'basket.html'),
            },
        },
    },
    test: {
        globals: true,
    },
    server: {
        proxy: {
            // string shorthand: http://localhost:5173/api -> http://localhost:8080/api
            '/ws': {
                target: 'http://localhost:8080', // Your Go backend server
                changeOrigin: true, // Needed for virtual hosted sites
                secure: false,      // Can be false if you're not using HTTPS
                ws: true,           // If you want to proxy websockets
            }
        }
    },
});