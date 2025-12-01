import path from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import AutoImport from 'unplugin-auto-import/vite';
import Icons from 'unplugin-icons/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Components from 'unplugin-vue-components/vite'
// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss(),
    Icons({ /* options */ }), 
    Components({
      resolvers: [
        IconsResolver({
          prefix: false,
          
        }),
      ],
    }),
    AutoImport({
                include: [
                    /\.[tj]sx?$/, // .ts, .tsx, .js, .jsx
                    /\.vue$/,
                    /\.vue\?vue/, // .vue
                    /\.md$/ // .md
                ],
                // global imports to register
                imports: [ // presets
                    'vue',
                    { // custom
                        '@vueuse/core': [
                            // named imports
                            'useMouse', // import { useMouse } from '@vueuse/core',
                            // alias
                            ['useFetch', 'useMyFetch']
                        ],
                        'axios': [
                            // default imports
                            ['default', 'axios']
                        ],
                        'vue': ['PropType', 'defineProps', 'InjectionKey', 'Ref']
                    }
                ],
                dirs: ['src/components/**', 'src/components/ui/**'],
                dts: 'src/components.d.ts',
                vueTemplate: false,
                eslintrc: {
                    enabled: false, 
                    filepath: './.eslintrc-auto-import.json',
                    globalsPropValue: true
                }
            })],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        secure: false,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      }
    }
  }  
})