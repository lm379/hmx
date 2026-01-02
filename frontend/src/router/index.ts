import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import LoginView from '../views/LoginView.vue';
import RegisterView from '../views/RegisterView.vue';
import ForgetPasswordView from '../views/ForgetPasswordView.vue';
import ProfileView from '../views/ProfileView.vue';
import ArtistListView from '../views/ArtistListView.vue';
import PlayView from '../views/PlayView.vue';
import ArtistProfileView from '../views/ArtistProfileView.vue';
import NotFound from '../views/NotFound.vue';
import AdminLayout from '../views/admin/AdminLayout.vue';
import DashboardView from '../views/admin/DashboardView.vue';
import VideoManager from '../views/admin/VideoManager.vue';
import ArtistManager from '../views/admin/ArtistManager.vue';
import UserManager from '../views/admin/UserManager.vue';
import { useAuthStore } from '../stores/auth';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: HomeView,
    meta: { showCategoryBar: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterView
  },
  {
    path: '/forget-password',
    name: 'ForgetPassword',
    component: ForgetPasswordView
  },
  {
    path: '/profile',
    name: 'Profile',
    component: ProfileView,
    meta: { requiresAuth: true }
  },
  {
    path: '/artists',
    name: 'Artists',
    component: ArtistListView,
    meta: { showCategoryBar: true }
  },
  {
    path: '/artist/:id',
    name: 'ArtistProfile',
    component: ArtistProfileView
  },
  {
    path: '/video/:id',
    name: 'Play',
    component: PlayView,
    meta: { showCategoryBar: true }
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: '仪表盘',
        component: DashboardView
      },
      {
        path: 'videos',
        name: '视频管理',
        component: VideoManager
      },
      {
        path: 'artists',
        name: '艺术家管理',
        component: ArtistManager
      },
      {
        path: 'users',
        name: '用户管理',
        component: UserManager
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: NotFound
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
