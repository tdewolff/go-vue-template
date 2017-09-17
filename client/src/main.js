import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

window.$ = window.jQuery = require('jquery')
require('semantic-ui-css/semantic.css')
require('semantic-ui-css/semantic.js')

import router from './router'
import store from './store'
import Auth from '@tdewolff/auth'
import App from './App'

Vue.config.productionTip = false

Vue.use(VueAxios, axios)
Vue.use(Auth, store, router)

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  template: '<App/>',
  components: { App }
})
