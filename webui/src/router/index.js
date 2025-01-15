import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/session' },
    { path: '/home', component: HomeView },
    { path: '/search', component: HomeView },
    { path: '/session', component: LoginView },
    { path: '/some/:id/link', component: HomeView },
  ]
})

export default router
