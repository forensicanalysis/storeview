import Vue from 'vue';
import vuetify from './plugins/vuetify';
import axios from 'axios';
import VueAxios from 'vue-axios';
import VueLodash from 'vue-lodash';


import App from './App.vue';

import './registerServiceWorker';
import router from './router';
import store from './store';

import 'roboto-fontface/css/roboto/roboto-fontface.css';
import 'roboto-fontface/css/roboto-slab/roboto-slab-fontface.css';
import '@mdi/font/css/materialdesignicons.css';


Vue.use(VueAxios, axios);
Vue.use(VueLodash);


Vue.config.productionTip = false;

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App),
}).$mount('#app');
