import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/pages/Home'
import Auth from '@/pages/Auth'
import AuthCallback from '@/pages/AuthCallback'
import API from '@/pages/API'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    },
    {
      path: '/auth',
      name: 'Auth',
      component: Auth
    },
    {
      path: '/auth/callback',
      name: 'AuthCallback',
      component: AuthCallback
    },
    {
      path: '/page',
      name: 'API',
      component: API
    }
  ]
})
