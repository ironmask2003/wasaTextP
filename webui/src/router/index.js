import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ChatView from '../views/ChatView.vue'
import GroupInfo from '../views/GroupInfo.vue'
import ProgileView from '../views/ProgileView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/session', redirect: '/' },
    { path: '/login', redirect: '/' },
    { path: '/', component: LoginView },
    { path: '/home', component: HomeView },
    { path: '/conversation', component: ChatView },
    { path: '/conversation/:convId', component: ChatView },
    { path: '/groups/:groupId', component: GroupInfo },
    { path: '/profile/:userId', component: ProgileView },
  ]
})

export default router
