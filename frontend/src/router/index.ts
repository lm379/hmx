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
  next();
});

export default router;
