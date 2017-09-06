import Vue from 'vue'
import Router from 'vue-router'
import API from '@/components/API'
import Auth from '@/components/Auth'
import AuthCallback from '@/components/AuthCallback'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
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
      path: '/',
      name: 'API',
      component: API
    }
  ]
})
