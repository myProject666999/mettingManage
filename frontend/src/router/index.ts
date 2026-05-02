import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: { title: '登录', requiresAuth: false }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/register/index.vue'),
      meta: { title: '注册', requiresAuth: false }
    },
    {
      path: '/user',
      name: 'UserLayout',
      component: () => import('@/layouts/UserLayout.vue'),
      meta: { requiresAuth: true, role: 'user' },
      children: [
        {
          path: '',
          redirect: '/user/home'
        },
        {
          path: 'home',
          name: 'UserHome',
          component: () => import('@/views/user/home/index.vue'),
          meta: { title: '首页' }
        },
        {
          path: 'meeting-rooms',
          name: 'UserMeetingRooms',
          component: () => import('@/views/user/meeting-rooms/index.vue'),
          meta: { title: '会议室' }
        },
        {
          path: 'meetings',
          name: 'UserMeetings',
          component: () => import('@/views/user/meetings/index.vue'),
          meta: { title: '我的会议' }
        },
        {
          path: 'participations',
          name: 'UserParticipations',
          component: () => import('@/views/user/participations/index.vue'),
          meta: { title: '参会信息' }
        },
        {
          path: 'announcements',
          name: 'UserAnnouncements',
          component: () => import('@/views/user/announcements/index.vue'),
          meta: { title: '公告' }
        },
        {
          path: 'reminders',
          name: 'UserReminders',
          component: () => import('@/views/user/reminders/index.vue'),
          meta: { title: '会议提醒' }
        },
        {
          path: 'profile',
          name: 'UserProfile',
          component: () => import('@/views/user/profile/index.vue'),
          meta: { title: '个人信息' }
        },
        {
          path: 'change-password',
          name: 'UserChangePassword',
          component: () => import('@/views/user/change-password/index.vue'),
          meta: { title: '修改密码' }
        }
      ]
    },
    {
      path: '/admin',
      name: 'AdminLayout',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true, role: 'admin' },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard'
        },
        {
          path: 'dashboard',
          name: 'AdminDashboard',
          component: () => import('@/views/admin/dashboard/index.vue'),
          meta: { title: '仪表盘' }
        },
        {
          path: 'meeting-rooms',
          name: 'AdminMeetingRooms',
          component: () => import('@/views/admin/meeting-rooms/index.vue'),
          meta: { title: '会议室管理' }
        },
        {
          path: 'meetings',
          name: 'AdminMeetings',
          component: () => import('@/views/admin/meetings/index.vue'),
          meta: { title: '会议管理' }
        },
        {
          path: 'participants',
          name: 'AdminParticipants',
          component: () => import('@/views/admin/participants/index.vue'),
          meta: { title: '参会信息管理' }
        },
        {
          path: 'documents',
          name: 'AdminDocuments',
          component: () => import('@/views/admin/documents/index.vue'),
          meta: { title: '会议资料管理' }
        },
        {
          path: 'reminders',
          name: 'AdminReminders',
          component: () => import('@/views/admin/reminders/index.vue'),
          meta: { title: '会议提醒管理' }
        },
        {
          path: 'announcements',
          name: 'AdminAnnouncements',
          component: () => import('@/views/admin/announcements/index.vue'),
          meta: { title: '公告管理' }
        },
        {
          path: 'carousels',
          name: 'AdminCarousels',
          component: () => import('@/views/admin/carousels/index.vue'),
          meta: { title: '轮播图管理' }
        },
        {
          path: 'users',
          name: 'AdminUsers',
          component: () => import('@/views/admin/users/index.vue'),
          meta: { title: '用户管理' }
        },
        {
          path: 'admins',
          name: 'AdminAdmins',
          component: () => import('@/views/admin/admins/index.vue'),
          meta: { title: '管理员管理' }
        },
        {
          path: 'profile',
          name: 'AdminProfile',
          component: () => import('@/views/admin/profile/index.vue'),
          meta: { title: '个人信息' }
        },
        {
          path: 'change-password',
          name: 'AdminChangePassword',
          component: () => import('@/views/admin/change-password/index.vue'),
          meta: { title: '修改密码' }
        }
      ]
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  const token = userStore.token

  document.title = to.meta.title ? `${to.meta.title} - 会议室管理系统` : '会议室管理系统'

  if (to.meta.requiresAuth) {
    if (!token) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }

    // 检查角色
    const requiredRole = to.meta.role
    if (requiredRole && userStore.role !== requiredRole) {
      // 如果角色不匹配，重定向到对应角色的首页
      if (userStore.role === 'user') {
        next({ name: 'UserHome' })
      } else if (userStore.role === 'admin' || userStore.role === 'super_admin') {
        next({ name: 'AdminDashboard' })
      } else {
        next({ name: 'Login' })
      }
      return
    }
  }

  next()
})

export default router
