import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import ProfileView from '../views/ProfileView.vue'

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/LoginView.vue')
    },
    {
        path: '/',
        name: 'Dashboard',
        component: () => import('../views/DashboardView.vue'),
        meta: { requiresAuth: true } // Custom field to protect this route
    },
    {
        path: '/board/:id',
        name: 'Board',
        component: () => import('../views/BoardView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/signup',
        name: 'Signup',
        component: () => import('../views/SignUpView.vue')
    },
    {
        path: '/profile',
        name: 'profile',
        component: ProfileView,
        meta: { requiresAuth: true } // Assuming you protect authenticated routes
    },
    {
        path: '/forgot-password',
        name: 'ForgotPassword',
        component: () => import('../views/ForgotPasswordView.vue')
    },
    {
        path: '/reset-password',
        name: 'ResetPassword',
        component: () => import('../views/ResetPasswordView.vue')
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

// Navigation Guard: Checks JWT before every page load
router.beforeEach((to, from, next) => {
    const auth = useAuthStore();

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        next('/login');
    } else {
        next();
    }
});

export default router;