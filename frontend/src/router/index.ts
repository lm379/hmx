import { createRouter, createWebHistory } from 'vue-router';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/HomeView.vue'),
    meta: { showCategoryBar: true }
  },
  {
    path: '/artists',
    name: 'Artists',
    component: () => import('../views/ArtistListView.vue'),
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
    component: () => import('../views/ProfileView.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
