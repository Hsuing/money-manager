import { createApp } from 'vue'
import './assets/styles/global.css'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import App from './App.vue'
import router from './router'

// 强制开启 Element Plus 暗色模式
document.documentElement.classList.add('dark')

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.mount('#app')
