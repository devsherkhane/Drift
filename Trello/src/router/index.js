import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import ProfileView from '../views/ProfileView.vue'

const routes = [
    {
        path: '/welcome',
        name: 'Landing',
        component: () => import('../views/LandingView.vue'),
        meta: { guestOnly: true }
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/LoginView.vue'),
        meta: { guestOnly: true }
    },
    {
        path: '/',
        name: 'Dashboard',
        component: () => import('../views/DashboardView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/board/:id',
        name: 'Board',
        component: () => import('../views/BoardView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/board/:id/activity',
        name: 'BoardActivity',
        component: () => import('../views/BoardActivityView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/signup',
        name: 'Signup',
        component: () => import('../views/SignUpView.vue'),
        meta: { guestOnly: true }
    },
    {
        path: '/profile',
        name: 'profile',
        component: ProfileView,
        meta: { requiresAuth: true }
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
    {
        path: '/search',
        name: 'Search',
        component: () => import('../views/SearchView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/activity',
        name: 'Activity',
        component: () => import('../views/ActivityView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/archived',
        name: 'ArchivedBoards',
        component: () => import('../views/ArchivedBoardsView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/templates',
        name: 'Templates',
        component: () => import('../views/TemplatesView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('../views/NotFoundView.vue')
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

// Navigation Guard: Checks JWT before every page load
router.beforeEach((to, from, next) => {
    const auth = useAuthStore();

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        next('/welcome');
    } else if (to.meta.guestOnly && auth.isAuthenticated) {
        next('/');
    } else {
        next();
    }
});

export default router;