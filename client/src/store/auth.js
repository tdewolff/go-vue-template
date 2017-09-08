import jwtDecode from 'jwt-decode'

export default {
  state: {
    jwt: localStorage.getItem('jwt')
  },
  mutations: {
    LOGIN: (state, jwt) => {
      state.jwt = jwt
    },
    LOGOUT: (state) => {
      state.jwt = null
    }
  },
  actions: {
    login ({commit}, jwt) {
      localStorage.setItem('jwt', jwt)
      commit('LOGIN', jwt)
    },
    logout ({commit}) {
      localStorage.removeItem('jwt')
      commit('LOGOUT')
    }
  },
  getters: {
    isLoggedIn: state => {
      if (state.jwt) {
        var claims = jwtDecode(state.jwt)
        var now = Math.floor(Date.now() / 1000)
        if (claims['exp'] > now) {
          return true
        }
      }
      return false
    },
    jwt: (state, getters) => {
      if (getters.isLoggedIn) {
        return state.jwt
      }
      return null
    },
    user: (state, getters) => {
      if (getters.isLoggedIn) {
        return jwtDecode(state.jwt)
      }
      return null
    }
  }
}
