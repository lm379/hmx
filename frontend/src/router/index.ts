import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/HomeView.vue'),
    meta: { showCategoryBar: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/LoginView.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/RegisterView.vue')
  },
  {
    path: '/forget-password',
    name: 'ForgetPassword',
    component: () => import('../views/ForgetPasswordView.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/artists',
    name: 'Artists',
    component: () => import('../views/ArtistListView.vue'),
    meta: { showCategoryBar: true }
  },
  {
    path: '/ranking',
    name: 'Ranking',
    component: () => import('../views/RankingView.vue'),
    meta: { showCategoryBar: true }
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('../views/SearchView.vue'),
    meta: { showCategoryBar: true }
  },
  {
    path: '/artist/:id',
    name: 'ArtistProfile',
    component: () => import('../views/ArtistProfileView.vue')
  },
  {
    path: '/video/:id',
    name: 'Play',
    component: () => import('../views/PlayView.vue'),
    meta: { showCategoryBar: true }
  },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: '仪表盘',
        component: () => import('../views/admin/DashboardView.vue')
      },
      {
        path: 'videos',
        name: '视频管理',
        component: () => import('../views/admin/VideoManager.vue')
      },
      {
        path: 'artists',
        name: '艺术家管理',
        component: () => import('../views/admin/ArtistManager.vue')
      },
      {
        path: 'users',
        name: '用户管理',
        component: () => import('../views/admin/UserManager.vue')
      },
      {
        path: 'tasks',
        name: '任务队列',
        component: () => import('../views/admin/TaskQueueView.vue')
      },
      {
        path: 'search',
        name: '搜索管理',
        component: () => import('../views/admin/SearchManager.vue')
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore();

  if (to.meta.requiresAuth) {
    if (!authStore.isLoggedIn) {
      // Try to restore session
      try {
        await authStore.checkLoginStatus();
        if (!authStore.isLoggedIn) {
          next({ name: 'Login', query: { redirect: to.fullPath } });
          return;
        }
      } catch (e) {
        next({ name: 'Login', query: { redirect: to.fullPath } });
        return;
      }
    }
  }

  if (to.meta.requiresAdmin) {
    // Ensure user info is loaded (might be logged in but user object is stale if page refreshed)
    // Actually checkLoginStatus does this.
    // Check role
    if (authStore.user?.role !== 'Administrator') {
      next({ name: 'Home' });
      return;
    }
  }

  next();
});

export default router;
