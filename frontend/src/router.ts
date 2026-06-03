import { createRouter, createWebHistory } from 'vue-router';

import AboutView from '@/views/AboutView.vue';
import CategoriesView from '@/views/CategoriesView.vue';
import DiagnosticsView from '@/views/DiagnosticsView.vue';
import FeedbackView from '@/views/FeedbackView.vue';
import PrivacyView from '@/views/PrivacyView.vue';
import SearchView from '@/views/SearchView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/search' },
    { path: '/search', name: 'search', component: SearchView },
    { path: '/categories', name: 'categories', component: CategoriesView },
    { path: '/categories/:kind/:subcategory', name: 'category-detail', component: CategoriesView },
    { path: '/about', name: 'about', component: AboutView },
    { path: '/feedback', name: 'feedback', component: FeedbackView },
    { path: '/privacy/aip-food-lookup', name: 'privacy', component: PrivacyView },
    { path: '/diagnostics', name: 'diagnostics', component: DiagnosticsView },
  ],
});

export default router;
