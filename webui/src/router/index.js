import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ChatView from '../views/ChatView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/session' },
    { path: '/home', component: HomeView },
    { path: '/session', component: LoginView },
    { path: '/conversation', component: ChatView },
    { path: '/conversation/:convId', component: ChatView },
  ]
})

export default router
