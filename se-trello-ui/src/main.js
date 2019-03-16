import Vue from 'vue'

import App from './App.vue'
import Board from './components/Boards.vue'
import List from './components/List.vue'

import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import VueRouter from 'vue-router'

Vue.config.productionTip = false

Vue.use(BootstrapVue) 
Vue.use(VueRouter)

const router = new VueRouter({
  routes: [
    { path: '/', component: Board },
    { path: '/list/:id', component: List }
  ]
})

new Vue({
  render: h => h(App),
  router,
}).$mount('#app')
