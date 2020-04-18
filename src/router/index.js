import Vue from 'vue';
import VueRouter from 'vue-router';
import items from '../views/Items.vue';
import file from '../views/File.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    redirect: '/items',
  },
  {
    path: '/items',
    name: 'items',
    component: items,
  },
  {
    path: '/file',
    name: 'file',
    component: file,
  },
  {
    path: '/workflows',
    name: 'workflows',
    component: () => import('../views/Workflows.vue'),
  },
];

const router = new VueRouter({
  routes,
});

export default router;
