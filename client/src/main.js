// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

window.$ = window.jQuery = require('jquery')
require('semantic-ui-css/semantic.css')
require('semantic-ui-css/semantic.js')

import router from './router'
import store from './store'
import App from './App'

Vue.config.productionTip = false

Vue.use(VueAxios, axios)

Vue.axios.interceptors.request.use((request) => {
  if (process.env.NODE_ENV === 'development') {
    console.log('Request Interceptor: ' + JSON.stringify(request))
  }
  if (store.getters.isLoggedIn) {
    request.headers['authorization'] = store.getters.jwt
  }
  return request
})

Vue.axios.interceptors.response.use((response) => {
  if (process.env.NODE_ENV === 'development') {
    console.log('Response Interceptor: ' + JSON.stringify(response))
  }
  if (response.headers['set-authorization']) {
    store.dispatch('login', response.headers['set-authorization'])
  }
  return response
}, function (error) {
  if (error.response && error.response.status === 401) {
    router.push('/auth?referer=' + encodeURIComponent(router.history.current.path))
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
