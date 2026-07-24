import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [react()],
	// Keep .gitkeep so the Go embed target exists on a fresh checkout.
	build: {emptyOutDir: false}
})
