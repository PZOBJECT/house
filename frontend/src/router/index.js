import { createRouter, createWebHistory } from 'vue-router'
import RoomManage from '../views/RoomManage.vue'
import BillManage from '../views/BillManage.vue'
import BillHistory from '../views/BillHistory.vue'
import Dashboard from '../views/Dashboard.vue'

const routes = [
  { path: '/', name: 'RoomManage', component: RoomManage },
  { path: '/bills', name: 'BillManage', component: BillManage },
  { path: '/history', name: 'BillHistory', component: BillHistory },
  { path: '/dashboard', name: 'Dashboard', component: Dashboard }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
