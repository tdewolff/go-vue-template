import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import Toasted from 'vue-toasted'

window.$ = window.jQuery = require('jquery')
require('semantic-ui-css/semantic.css')
require('semantic-ui-css/semantic.js')

import config from '../../config'
import router from './router'
import store from './store'
import Auth from '@tdewolff/auth'
import App from './App'

Vue.config.productionTip = false

Vue.use(VueAxios, axios)
Vue.use(Toasted, {
  'position': 'top-center',
  'duration': 5000
})
Vue.use(Auth, store, config.URL)

Vue.axios.interceptors.request.use(request => {
  if (process.env.NODE_ENV === 'development') {
    console.log('Request: ' + JSON.stringify(request))
  }
  return request
})

Vue.axios.interceptors.response.use(response => {
  if (process.env.NODE_ENV === 'development') {
    console.log('Response: ' + JSON.stringify(response))
  }
  return response
}, error => {
  if (error.response && error.response.status === 401) {
    router.push('/auth?referrer=' + encodeURIComponent(router.history.current.path))
  }
  return Promise.reject(error)
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  template: '<App/>',
  components: { App }
})
