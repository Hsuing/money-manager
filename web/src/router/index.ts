import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue')
  },
  {
    path: '/ledger',
    name: 'Ledger',
    component: () => import('../views/Ledger.vue')
  },
  {
    path: '/analysis',
    name: 'Analysis',
    component: () => import('../views/Analysis.vue')
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/Settings.vue')
  },
  {
    path: '/settings/income-expense',
    name: 'IncomeExpenseSettings',
    component: () => import('../views/settings/IncomeExpenseSettings.vue')
  },
  {
    path: '/settings/category',
    name: 'CategorySettings',
    component: () => import('../views/settings/CategorySettings.vue')
  },
  {
    path: '/settings/account',
    name: 'AccountSettings',
    component: () => import('../views/settings/AccountSettings.vue')
  },
  {
    path: '/settings/template',
    name: 'TemplateSettings',
    component: () => import('../views/settings/TemplateSettings.vue')
  },
  {
    path: '/settings/schedule',
    name: 'ScheduleSettings',
    component: () => import('../views/settings/ScheduleSettings.vue')
  },
  {
    path: '/settings/system',
    name: 'SystemSettings',
    component: () => import('../views/settings/SystemSettings.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.name !== 'Login' && !token) {
    next({ name: 'Login' })
  } else if (to.name === 'Login' && token) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
